package project

import (
	"fmt"
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

func GetProjects(provider *gophercloud.ProviderClient) ([]models.ProjectData, error) {
  var projectsListOutput []models.ProjectData
  client, err := openstack.NewIdentityV3(provider, gophercloud.EndpointOpts{})
  if err != nil {
    return nil, err
  }

  listOpts := projects.ListOpts{
    Enabled: gophercloud.Enabled,
    DomainID: os.Getenv("DOMAIN_ID"),
  }
  rows, err := projects.List(client, listOpts).AllPages()
  if err != nil {
    return nil, err
  }

  projectList, err := projects.ExtractProjects(rows)
  if err != nil {
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
  return projectsListOutput, nil
}

func GetProjectQuota(provider *gophercloud.ProviderClient, projectID string) (*quotasets.QuotaSet, error) {
  client, err := openstack.NewComputeV2(provider, gophercloud.EndpointOpts{
    Region: os.Getenv("OS_REGION_NAME"),
  })
  if err != nil {
    return nil, err
  }

  quotas, err := quotasets.Get(client, projectID).Extract()
  if err != nil {
    return nil, err
  }

  return quotas, nil
}

func GetProjectUsage(provider *gophercloud.ProviderClient, projectID string) (*models.ProjectSumUsage, error) {
  client, err := openstack.NewComputeV2(provider, gophercloud.EndpointOpts{
    Region: os.Getenv("OS_REGION_NAME"),
  })
  if err != nil {
    return nil, err
  }

  start := time.Now().AddDate(0, 0, -1)
  end := time.Now()
  singleTenantOpts := usage.SingleTenantOpts{
    Start: &start,
    End: &end,
  }
  VCPUSum := 0
  MemorySum := 0
  err = usage.SingleTenant(client, projectID, singleTenantOpts).EachPage(func(page pagination.Page) (bool, error) {
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
    return nil, err
  }

  projectSumUsage := models.ProjectSumUsage{VcpuUsage: VCPUSum, RAMUsage: MemorySum}

  return &projectSumUsage, err
}

func SaveProjectSummary(provider *gophercloud.ProviderClient) error {
  projects, err := GetProjects(provider)
  if err != nil {
    return err
  }


  timestamp := time.Now()

  for _, projectData := range projects {
    quotas, err := GetProjectQuota(provider, projectData.ID)
    if err != nil {
      fmt.Println("[ERROR] Error getting project quota", err)
      return err
    }

    projectUsage, err := GetProjectUsage(provider, projectData.ID)
    if err != nil {
      fmt.Println("[ERROR] Error getting project usage", err)
      return err
    }

    report := models.UsageReport{Timestamp: timestamp, Sponsor: projectData.Sponsor, ProjectName: projectData.Name,
      VCPUQuota: quotas.Cores, VCPUUsage: projectUsage.VcpuUsage, RAMQuota: quotas.RAM, RAMUsage: projectUsage.RAMUsage}

    models.DB.Create(&report)

  }

  return nil

}

//func GetSponsorSummary(provider *gophercloud.ProviderClient) (map[string]models.UsageReport, error) {
//  
//  projects, err := GetProjects(provider)
//  if err != nil {
//    return nil, err
//  }
//
//  aggregateReports := make(map[string]models.UsageReport)
//
//  for _, projectData := range projects {
//    quotas, err := GetProjectQuota(provider, projectData.ID)
//    if err != nil {
//      fmt.Println("[ERROR] Error getting project quota", err)
//      continue
//    }
//
//    projectUsage, err := GetProjectUsage(provider, projectData.ID)
//    if err != nil {
//      fmt.Println("[ERROR] Error getting project usage", err)
//      continue
//    }
//    timestamp, _ := time.Now().MarshalText()
//
//    report, ok := aggregateReports[projectData.Sponsor]
//    if !ok {
//      report = models.UsageReport{Timestamp: &timestamp, Sponsor: projectData.Sponsor}
//    }
//
//    report.VCPUQuota += quotas.Cores
//    report.VCPUUsage += projectUsage.VcpuUsage
//    report.RAMQuota += quotas.RAM
//    report.RAMUsage += projectUsage.RAMUsage
//
//    aggregateReports[projectData.Sponsor] = report
//  }
//
//  return aggregateReports, nil
//}
