package controllers

import (
	"fmt"
	"strconv"
	"time"

	"net/url"

	"github.com/AlissonBarbosa/shylockgo/models"
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/servers"
)

func GetAllServers(provider *gophercloud.ProviderClient) ([]models.ServerMeta, error) {
  var allServersMeta []models.ServerMeta
  var maxEpoch int64

  models.DB.Model(&models.ServerMeta{}).Select("MAX(epoch)").Scan(&maxEpoch)
  models.DB.Where("epoch = ?", maxEpoch).Find(&allServersMeta)

  return allServersMeta, nil
}

func SaveAllServers(provider *gophercloud.ProviderClient) error {
  client, err := openstack.NewComputeV2(provider, gophercloud.EndpointOpts{})
  if err != nil {
    return err
  }

  listOpts := servers.ListOpts{
    AllTenants: true,
  }
  rows, err := servers.List(client, listOpts).AllPages()
  if err != nil {
    return err
  }

  serverList, err := servers.ExtractServers(rows)
  if err != nil {
    return err
  }

  epoch := time.Now().Unix()

  for _, server := range serverList {
    domain, err := GetServerDomain(server.ID)
    if err != nil {
      return err
    }

    memoryUsage, err := GetServerMemoryUsage(domain)
    if err != nil {
      return err
    }
    memoryConverted, err := strconv.ParseFloat(memoryUsage, 64)
    if err != nil {
      return err
    }

    report := models.ServerMeta{ServerID: server.ID,
      Name: server.Name, ProjectID : server.TenantID, HostID: server.HostID,
      Domain: domain, MemoryUsage: int64(memoryConverted), Epoch: epoch}
    models.DB.Create(&report)
  }
  return nil
}

func GetServerDomain(id string) (string, error) {
  query := fmt.Sprintf("libvirt_domain_info_meta{uuid='%s'}", url.QueryEscape(id))
  result := QueryGetPrometheus(query)
  if result.Error != nil {
    return "Error", result.Error
  }
  domain := "None"
  if len(result.Data.(models.QueryResponse).Data.Result) > 0 {
    if result.Data.(models.QueryResponse).Data.Result[0].Metric.Domain != "" {
      domain = fmt.Sprintf("%v", result.Data.(models.QueryResponse).Data.Result[0].Metric.Domain)
    }
  }

  return domain, nil
}

func GetServerMemoryUsage(domain string) (string, error) {
  query := fmt.Sprintf("libvirt_domain_memory_stats_used_percent{domain='%s'}", url.QueryEscape(domain))
  result := QueryGetPrometheus(query)
  if result.Error != nil {
    return "Error", result.Error
  }
  memoryUsage := "None"
  if len(result.Data.(models.QueryResponse).Data.Result) > 0 {
    if len(result.Data.(models.QueryResponse).Data.Result[0].Value) > 1 {
      memoryUsage = fmt.Sprintf("%v",result.Data.(models.QueryResponse).Data.Result[0].Value[1])
    }
  }

  return memoryUsage, nil
}
