package server

import (
  "os"
  "time"
	"database/sql"
	"github.com/google/uuid"
  "github.com/gophercloud/gophercloud"
  "github.com/gophercloud/gophercloud/openstack"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/servers"
)

func SaveInstancesToDatabase(provider *gophercloud.ProviderClient, db *sql.DB) error {
  _, err := db.Exec(`
      CREATE TABLE IF NOT EXISTS instances (
        id TEXT NOT NULL PRIMARY KEY,
        timestamp TEXT,
        epoch BIGINT,
        instanceid TEXT,
        name TEXT,
        status TEXT,
        createdat TEXT,
        flavorid TEXT,
        tenantid TEXT,
        hostid TEXT,
        userid TEXT
      )
    `)

  if err != nil {
    return err
  }

  client, err := openstack.NewComputeV2(provider, gophercloud.EndpointOpts{
		Region: os.Getenv("OS_REGION_NAME"),
	})
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

	instanceList, err := servers.ExtractServers(rows)

	if err != nil {
		return err
	}

	timestamp := time.Now()
	for _, instance := range instanceList {
		_, err := db.Exec("INSERT INTO instances (id, timestamp, epoch, instanceid, name, status, createdat, flavorid, tenantid, hostid, userid) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)",
			uuid.New(), timestamp, timestamp.Unix() ,instance.ID, instance.Name, instance.Status, instance.Created, instance.Flavor["id"], instance.TenantID, instance.HostID, instance.UserID)
		if err != nil {
			return err
		}
	}

  return nil
}
