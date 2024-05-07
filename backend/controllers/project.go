package controllers

import (
	"log/slog"
	"os"
	"regexp"
	"time"

	"github.com/AlissonBarbosa/shylockgo/models"
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/quotasets"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/usage"
	"github.com/gophercloud/gophercloud/openstack/identity/v3/projects"
	"github.com/gophercloud/gophercloud/pagination"
)

func GetAllProjects(provider *gophercloud.ProviderClient) ([]models.ProjectData, error) {
  var projectsListOutput []models.ProjectData
  client, err := openstack.NewIdentityV3(provider, gophercloud.EndpointOpts{})
  if err != nil {
    slog.Error("Error creating openstack client", err)
    return nil, err
  }

  listOpts := projects.ListOpts{
    Enabled: gophercloud.Enabled,
    DomainID: os.Getenv("DOMAIN_ID"),
  }
  rows, err := projects.List(client, listOpts).AllPages()
  if err != nil {
    slog.Error("Error getting project list", err)
    return nil, err
  }

  projectList, err := projects.ExtractProjects(rows)
  if err != nil {
    slog.Error("Error extracting project list", err)
    return nil, err
  }

  for _, project := range projectList {
    sponsor := project.Description
    re := regexp.MustCompile(`Responsavel(?:\(is\))?:\s+(\S+)@`)
    match := re.FindStringSubmatch(project.Description)
    if len(match) > 1 {
      sponsor = match[1]
    }

    projectsListOutput = append(projectsListOutput, models.ProjectData{ID: project.ID, Sponsor: sponsor, Name: project.Name})
  }

  slog.Info("Returning all projects")
  return projectsListOutput, nil
}

func SaveProjectsDesc(provider *gophercloud.ProviderClient) error {
  projectList, err := GetAllProjects(provider)
  if err != nil {
    slog.Error("Error getting all projects", err)
    return err
  }
  epoch := time.Now().Unix()

  for _, project := range projectList {
    projectToSave := models.ProjectDesc{Timestamp: epoch, ProjectID: project.ID, ProjectName: project.Name, ProjectSponsor: project.Sponsor}
    models.DB.Create(&projectToSave)
  }
  slog.Info("Projects Descriptions saved on database")
  return nil
}

func GetProjectQuota(provider *gophercloud.ProviderClient) error {
  projectList, err := GetAllProjects(provider)
  if err != nil {
    slog.Error("Error getting all projects", err)
    return err
  }

  client, err := openstack.NewComputeV2(provider, gophercloud.EndpointOpts{
    Region: os.Getenv("OS_REGION_NAME"),
  })
  if err != nil {
    slog.Error("Error creating openstack client", err)
    return err
  }

  epoch := time.Now().Unix()

  for _, project := range projectList {
    quotas, err := quotasets.Get(client, project.ID).Extract()
    if err != nil {
      slog.Error("Error getting project quota", err)
      return err
    }
    quotaToSave := models.ProjectQuota{Timestamp: epoch, ProjectID: project.ID, QuotaRam: int64(quotas.RAM), QuotaVcpu: int64(quotas.Cores)}
    models.DB.Create(&quotaToSave)
  }
  
  slog.Info("Projects quotas saved on database")
  return nil
}

func GetProjectUsage(provider *gophercloud.ProviderClient) error {
  projectList, err := GetAllProjects(provider)
  if err != nil {
    slog.Error("Error getting all projects", err)
    return err
  }

  client, err := openstack.NewComputeV2(provider, gophercloud.EndpointOpts{
    Region: os.Getenv("OS_REGION_NAME"),
  })
  if err != nil {
    slog.Error("Error creating openstack client", err)
    return err
  }

  start := time.Now().Add(time.Duration(-5) * time.Minute)
  end := time.Now()
  singleTenantOpts := usage.SingleTenantOpts{
    Start: &start,
    End: &end,
  }
  epoch := time.Now().Unix()
  
  for _, project := range projectList {
    VCPUSum := 0
    MemorySum := 0
    err = usage.SingleTenant(client, project.ID, singleTenantOpts).EachPage(func(page pagination.Page) (bool, error) {
      tenantUsage, err := usage.ExtractSingleTenant(page)
      if err != nil {
        return false, err
      }
      for _, server := range tenantUsage.ServerUsages {
        VCPUSum += server.VCPUs
        MemorySum += server.MemoryMB
      }
      return true, nil
    })

    if err != nil {
      slog.Error("Error getting project quota usage", err)
      return err
    }
    projectSumUsage := models.ProjectQuotaUsage{Timestamp:epoch, ProjectID: project.ID ,VcpuUsage: int64(VCPUSum), RamUsage: int64(MemorySum)}

    models.DB.Create(&projectSumUsage)
  }

  slog.Info("Project quota usage saved on database")
  return nil
}
