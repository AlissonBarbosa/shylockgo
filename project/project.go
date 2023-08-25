package project

import (
  "time"
	"database/sql"
	"github.com/google/uuid"
  "github.com/gophercloud/gophercloud"
  "github.com/gophercloud/gophercloud/openstack"
	"github.com/gophercloud/gophercloud/openstack/identity/v3/projects"
)

func SaveProjectsToDatabase(provider *gophercloud.ProviderClient, db *sql.DB) error{
  _, err := db.Exec(`
      CREATE TABLE IF NOT EXISTS projects (
        id TEXT NOT NULL PRIMARY KEY,
        timestamp TEXT,
        epoch BIGINT,
        projectid TEXT,
        name TEXT,
        description TEXT,
        domainid TEXT
      )
    `)

  if err != nil {
    return err
  }
  client, err := openstack.NewIdentityV3(provider, gophercloud.EndpointOpts{})
  if err != nil {
    return err
  }

  listOpts := projects.ListOpts{
    Enabled: gophercloud.Enabled,
  }
  rows, err := projects.List(client, listOpts).AllPages()
  if err != nil {
    return err
  }

  projectList, err := projects.ExtractProjects(rows)
  if err != nil {
    return err
  }

  timestamp := time.Now()
  for _, project := range projectList {
    _, err := db.Exec("INSERT INTO projects (id, timestamp, epoch, projectid, name, description, domainid) VALUES ($1, $2, $3, $4, $5, $6, $7)",
      uuid.New(), timestamp, timestamp.Unix(), project.ID, project.Name, project.Description, project.DomainID)
    if err != nil {
      return err
    }
  }
  return nil
}
