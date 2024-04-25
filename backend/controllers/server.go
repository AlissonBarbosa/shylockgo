package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"

	"net/http"
	"net/url"
	"os"

	"github.com/AlissonBarbosa/shylockgo/models"
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/servers"
)

func GetAllServers(provider *gophercloud.ProviderClient) ([]models.ServerData, error) {
  var serverListOutput []models.ServerData
  client, err := openstack.NewComputeV2(provider, gophercloud.EndpointOpts{})
  if err != nil {
    return nil, err
  }

  listOpts := servers.ListOpts{
    AllTenants: true,
  }
  rows, err := servers.List(client, listOpts).AllPages()
  if err != nil {
    return nil, err
  }

  serverList, err := servers.ExtractServers(rows)
  if err != nil {
    return nil, err
  }

  for _, server := range serverList {
    domain, err := GetServerDomain(server.ID)
    if err != nil {
      return nil, err
    }

    memoryUsage, err := GetServerMemoryUsage(domain)
    if err != nil {
      return nil, err
    }
    memoryConverted, err := strconv.ParseFloat(memoryUsage, 64)
    if err != nil {
      return nil, err
    }
    serverListOutput = append(serverListOutput, models.ServerData{ID: server.ID,
      Name: server.Name, ProjectID : server.TenantID, HostID: server.HostID,
      Domain: domain, MemoryUsage: int64(memoryConverted)})
    //fmt.Printf("%+v\n", server)
  }
  return serverListOutput, nil
}

// TODO: Write a generic funcion in common receiving model
func GetServerDomain(id string) (string, error) {
  query := fmt.Sprintf("libvirt_domain_info_meta{uuid='%s'}", url.QueryEscape(id))
  prometheus_url := fmt.Sprintf("%s:%s/api/v1/query?query=%s", os.Getenv("PROMETHEUS_URL"), os.Getenv("PROMETHEUS_PORT"), query)
  resp, err := http.Get(prometheus_url)
  if err != nil {
    return "Error", err
  }
  defer resp.Body.Close()
  
  body, err := ioutil.ReadAll(resp.Body)
  if err != nil {
    return "Error", err
  }

  var response models.DomainQueryResponse
  err = json.Unmarshal(body, &response)
  if err != nil {
    return "Error", err
  }
  domain := "None"
  if len(response.Data.Result) > 0 {
    domain = response.Data.Result[0].Metric.Domain
  }

  return domain, nil
}

//func GetServerMemoryUsage(domain string) (string, error) {
//  query := fmt.Sprintf("libvirt_domain_memory_stats_used_percent{domain='%s'}'", url.QueryEscape(domain))
//  result := QueryGetPrometheus(query)
//  if result.Error != nil {
//    return "Error", result.Error
//  }
//  fmt.Println(result)
//  memoryUsage := "None"
//  if len(result.Data.(models.QueryResponse).Data.Result) > 0 {
//    fmt.Printf("%v", result.Data.(models.QueryResponse).Data.Result)
//    if len(result.Data.(models.QueryResponse).Data.Result[0].Value) > 1 {
//      memoryUsage = fmt.Sprintf("%v",result.Data.(models.QueryResponse).Data.Result[0].Value[1])
//    }
//  }
//
//  return memoryUsage, nil
//}

func GetServerMemoryUsage(domain string) (string, error) {
  query := fmt.Sprintf("libvirt_domain_memory_stats_used_percent{domain='%s'}", url.QueryEscape(domain))
  prometheus_url := fmt.Sprintf("%s:%s/api/v1/query?query=%s", os.Getenv("PROMETHEUS_URL"), os.Getenv("PROMETHEUS_PORT"), query)
  resp, err := http.Get(prometheus_url)
  if err != nil {
    return "Error", err
  }
  defer resp.Body.Close()
  
  body, err := ioutil.ReadAll(resp.Body)
  if err != nil {
    return "Error", err
  }

  var response models.QueryResponse
  err = json.Unmarshal(body, &response)
  if err != nil {
    return "Error", err
  }
  memoryUsage := "none"
  if len(response.Data.Result) > 0 {
    memoryUsage = fmt.Sprintf("%v", response.Data.Result[0].Value[1]) 
  }
  return memoryUsage, nil
}
