package controllers

import (
	"fmt"
	"log/slog"
	"time"
	"github.com/AlissonBarbosa/shylockgo/backend/scraper/models"
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/flavors"
)

func getFlavorList(provider *gophercloud.ProviderClient) ([]flavors.Flavor, error) {
  client, err := openstack.NewComputeV2(provider, gophercloud.EndpointOpts{})
  if err != nil {
    slog.Error("Error creating openstack client", err)
    return nil, err
  }

  listOpts := flavors.ListOpts{
  }
  rows, err := flavors.ListDetail(client, listOpts).AllPages()
  if err != nil {
    slog.Error("Error getting flavor list", err)
    return nil, err
  }

  flavorList, err := flavors.ExtractFlavors(rows)
  if err != nil {
    slog.Error("Error extracting flavor list", err)
    return nil, err
  }
  return flavorList, nil
}

func SaveFlavorDesc(provider *gophercloud.ProviderClient) error {
  flavorList, err := getFlavorList(provider)
  if err != nil {
    slog.Error("Error getting flavor list from openstack", err)
    return err
  }
  epoch := time.Now().Unix()
  for _, flavor := range flavorList {
    flavorDesc := models.FlavorDesc{Timestamp: epoch, FlavorID: flavor.ID, FlavorName: flavor.Name}
    models.DB.Create(&flavorDesc)
  }
  slog.Info("Saving all flavors description")
  return nil
}

func SaveFlavorSpec(provider *gophercloud.ProviderClient) error {
  flavorList, err := getFlavorList(provider)
  if err != nil {
    slog.Error("Error getting flavor list from openstack", err)
    return err
  }
  epoch := time.Now().Unix()
  for _, flavor := range flavorList {
    flavorSpec := models.FlavorSpec{Timestamp: epoch, FlavorID: flavor.ID, Vcpu:fmt.Sprintf("%v", flavor.VCPUs), Ram: fmt.Sprintf("%v", flavor.RAM), Disk: fmt.Sprintf("%v", flavor.Disk)}
    models.DB.Create(&flavorSpec)
  }
  slog.Info("Saving all flavors specifications")
  return nil
}
