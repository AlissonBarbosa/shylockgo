package controllers

import (
	"errors"
	"fmt"
	"log/slog"
	"strconv"

	//"strconv"
	"time"

	"net/url"

	"github.com/AlissonBarbosa/shylockgo/models"
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/servers"
)

//func GetAllServers(provider *gophercloud.ProviderClient) ([]models.ServerMeta, error) {
//  var allServersMeta []models.ServerMeta
//  var maxEpoch int64
//
//  models.DB.Model(&models.ServerMeta{}).Select("MAX(epoch)").Scan(&maxEpoch)
//  models.DB.Where("epoch = ?", maxEpoch).Find(&allServersMeta)
//
//  return allServersMeta, nil
//}

func getServersList(provider *gophercloud.ProviderClient) ([]servers.Server, error) {
  client, err := openstack.NewComputeV2(provider, gophercloud.EndpointOpts{})
  if err != nil {
    slog.Error("Error creating openstack client", err)
    return nil, err
  }

  listOpts := servers.ListOpts{
    AllTenants: true,
  }
  rows, err := servers.List(client, listOpts).AllPages()
  if err != nil {
    slog.Error("Error getting server list", err)
    return nil, err
  }

  serverList, err := servers.ExtractServers(rows)
  if err != nil {
    slog.Error("Error extracting server list", err)
    return nil, err
  }
  return serverList, nil
}

func getFixedAddress(addresses map[string]interface{}) (string, error) {
  for _, addrPool := range addresses{
    if addrsList, ok := addrPool.([]interface{}); ok {
      for _, addr := range addrsList {
        addrMap, ok := addr.(map[string]interface{})
        if !ok {
          return "", errors.New("Failed to parse address")
        }
        addrType, ok := addrMap["OS-EXT-IPS:type"].(string)
        if !ok {
          return "", errors.New("Failed to parse address type")
        }
        if addrType == "fixed" {
          addrValue, ok := addrMap["addr"].(string)
          if !ok {
            return "", errors.New("Failed to parse address value")
          }
          return addrValue, nil
        }
      }
    }
  }
  return "None", nil
}

func getServerDomain(id string) (string, error) {
  query := fmt.Sprintf("libvirt_domain_info_meta{uuid='%s'}", url.QueryEscape(id))
  result := QueryGetPrometheus(query)
  if result.Error != nil {
    return "None", result.Error
  }
  if len(result.Data.(models.QueryResponse).Data.Result) > 0 {
    if result.Data.(models.QueryResponse).Data.Result[0].Metric.Domain != "" {
      domain := fmt.Sprintf("%v", result.Data.(models.QueryResponse).Data.Result[0].Metric.Domain)
      return domain, nil
    }
  }
  return "None", nil
}

func getServerMemoryUsage(domain string) (string, error) {
  query := fmt.Sprintf("libvirt_domain_memory_stats_used_percent{domain='%s'}", url.QueryEscape(domain))
  result := QueryGetPrometheus(query)
  if result.Error != nil {
    return "None", result.Error
  }
  if len(result.Data.(models.QueryResponse).Data.Result) > 0 {
    if len(result.Data.(models.QueryResponse).Data.Result[0].Value) > 1 {
      memoryUsage := fmt.Sprintf("%v",result.Data.(models.QueryResponse).Data.Result[0].Value[1])
      memoryConverted, err := strconv.ParseFloat(memoryUsage, 64)
      if err != nil {
        slog.Error("Error converting memory usage to float")
        return "None", err
      }
      return fmt.Sprintf("%.2f", memoryConverted), nil
    }
  }
  return "None", nil
}

func getServerVcpuUsage(domain string) (string, error) {
  query := fmt.Sprintf("irate(libvirt_domain_info_cpu_time_seconds_total{domain='%s'}[5m])*100", url.QueryEscape(domain))
  result := QueryGetPrometheus(query)
  if result.Error != nil {
    return "None", result.Error
  }
  if len(result.Data.(models.QueryResponse).Data.Result) > 0 {
    if len(result.Data.(models.QueryResponse).Data.Result[0].Value) > 1 {
      vcpuUsage := fmt.Sprintf("%v", result.Data.(models.QueryResponse).Data.Result[0].Value[1])
      vcpuConverted, err := strconv.ParseFloat(vcpuUsage, 64)
      if err != nil {
        slog.Error("Error converting vcpu usage to float")
        return "None", err
      }
      return fmt.Sprintf("%.2f", vcpuConverted), nil
    }
  }
  return "None", nil
}

func SaveServersDesc(provider *gophercloud.ProviderClient) error {
  serverList, err := getServersList(provider)
  if err != nil {
    slog.Error("Error getting server list from openstack", err)
    return err
  }
  epoch := time.Now().Unix()
  for _, server := range serverList {
    address, err := getFixedAddress(server.Addresses)
    if err != nil {
      slog.Error("Error getting server address", err)
      return err
    }
    serverDesc := models.ServerDesc{Timestamp: epoch,ServerID: server.ID, ServerName: server.Name, ServerAddress: address}
    models.DB.Create(&serverDesc)
  }
  slog.Info("Saving all servers description")
  return nil
}

func SaveServersSpec(provider *gophercloud.ProviderClient) error {
  serverList, err := getServersList(provider)
  if err != nil {
    slog.Error("Error getting server list from openstack", err)
    return err
  }
  epoch := time.Now().Unix()
  for _, server := range serverList {
    serverSpec := models.ServerSpec{Timestamp: epoch, ServerID: server.ID, FlavorID: server.Flavor["id"].(string)}
    models.DB.Create(&serverSpec)
  }
  slog.Info("Saving all servers specifications")
  return nil
}

func SaveServersUsage(provider *gophercloud.ProviderClient) error {
  serverList, err := getServersList(provider)
  if err != nil {
    slog.Error("Error getting server list from openstack", err)
    return err
  }
  epoch := time.Now().Unix()
  for _, server := range serverList {
    domain, err := getServerDomain(server.ID)
    if err != nil {
      slog.Error("Error geting server domain")
      return err
    }
    memoryUsage, err := getServerMemoryUsage(domain)
    if err != nil {
      slog.Error("Error getting memory usage", err)
      return err
    }
    vcpuUsage, err := getServerVcpuUsage(domain)
    if err != nil {
      slog.Error("Error getting vcpu usage", err)
      return err
    }
    serverUsage := models.ServerUsage{Timestamp: epoch, ServerID: server.ID, RamUsage: memoryUsage, VcpuUsage: vcpuUsage, Domain: domain, HostID: server.HostID}
    models.DB.Create(&serverUsage)
  }
  slog.Info("Saving all servers usage")
  return nil
}

func SaveServersOwnership(provider *gophercloud.ProviderClient) error {
  serverList, err := getServersList(provider)
  if err != nil {
    slog.Error("Error getting server list from openstack", err)
    return err
  }
  epoch := time.Now().Unix()
  for _, server := range serverList {
    serverOwnership := models.ServerOwnership{Timestamp: epoch, ServerID: server.ID, UserID: server.UserID, ProjectID: server.TenantID}
    models.DB.Create(&serverOwnership)
  }
  slog.Info("Saving all servers ownership")
  return nil
}
