package controllers

import (
	"log/slog"
  "os"
	"time"
	"github.com/AlissonBarbosa/shylockgo/models"
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	"github.com/gophercloud/gophercloud/openstack/identity/v3/users"
	"github.com/gophercloud/gophercloud/openstack/identity/v3/projects"
)

func getAllUsers(provider *gophercloud.ProviderClient) ([]users.User, error) {
  client, err := openstack.NewIdentityV3(provider, gophercloud.EndpointOpts{})
  if err != nil {
    slog.Error("Error creating identity client", err)
    return nil, err
  }
  listOpts := users.ListOpts{
    DomainID: os.Getenv("DOMAIN_ID"),
  }
  rows, err := users.List(client, listOpts).AllPages()
  if err != nil {
    slog.Error("Error getting user list from openstack", err)
    return nil, err
  }

  userList, err := users.ExtractUsers(rows)
  if err != nil {
    slog.Error("Error extracting user list", err)
    return nil, err
  }
  return userList, nil
}

func SaveUserDesc(provider *gophercloud.ProviderClient) error {
  userList, err := getAllUsers(provider)
  if err != nil {
    slog.Error("Error getting user list from openstack", err)
    return err
  }
  epoch := time.Now().Unix()
  for _, user := range userList {
    userDesc := models.UserDesc{Timestamp: epoch, UserID: user.ID, UserName: user.Name}
    models.DB.Create(&userDesc)
  }
  slog.Info("Saving all users description")
  return nil
}

func SaveUserProjects(provider *gophercloud.ProviderClient) error {
  userList, err := getAllUsers(provider)
  if err != nil {
    slog.Error("Error getting user list from openstack", err)
    return err
  }
  client, err := openstack.NewIdentityV3(provider, gophercloud.EndpointOpts{})
  if err != nil {
    slog.Error("Error creating openstack client", err)
    return err
  }

  epoch := time.Now().Unix()
  for _, user := range userList{
    rows, err := users.ListProjects(client, user.ID).AllPages()
    if err != nil {
      slog.Error("Error getting project list", err)
      return err
    }

    projectList, err := projects.ExtractProjects(rows)
    if err != nil {
      slog.Error("Error extracting project list", err)
      return err
    }

    for _, project := range projectList {
      userProject := models.UserProject{Timestamp: epoch, UserID: user.ID, ProjectID: project.ID}
      models.DB.Create(&userProject)
    }
  }

  slog.Info("Saving all project users")
  return nil
}
