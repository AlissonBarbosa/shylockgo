package main

import (
  "database/sql"
  "fmt"
  "os"
  "time"

  _ "github.com/mattn/go-sqlite3"
  "github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/servers"
)

func main()  {
  opts, err := openstack.AuthOptionsFromEnv()
  if err != nil {
    fmt.Println("Error reading environment variables:", err)
    os.Exit(1)
  }

  provider, err := openstack.AuthenticatedClient(opts)
  if err != nil {
    fmt.Println("Error creating authenticated client", err)
    os.Exit(1)
  }

  client, err := openstack.NewComputeV2(provider, gophercloud.EndpointOpts{
    Region: os.Getenv("OS_REGION_NAME"),
  })
  if err != nil {
    fmt.Println("Error creating ComputeV2 client:", err)
    os.Exit(1)
  }

  db, err := connectToDatabase()
  if err != nil {
    fmt.Println("Error connecting to database", err)
    os.Exit(1)
  }
  defer db.Close()

  _, err = db.Exec(`
    CREATE TABLE IF NOT EXISTS instances (
      id INTEGER NOT NULL PRIMARY KEY,
      timestamp TEXT,
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
  saveInstancesToDatabase(client, db)
}

func connectToDatabase() (*sql.DB, error) {
  db, err := sql.Open("sqlite3", os.Getenv("DATABASEFILE"))
  if err != nil {
    return nil, err
  }

  return db, nil
}

func saveInstancesToDatabase(client *gophercloud.ServiceClient, db *sql.DB) {
  listOpts := servers.ListOpts{
		AllTenants: true,
	}

  rows, err := servers.List(client, listOpts).AllPages()
  if err != nil {
    fmt.Println("Error listing instances:", err)
    return
  }
  
  instanceList, err := servers.ExtractServers(rows)

  if err != nil {
    fmt.Println("Error extracting instance list:", err)
    return
  }

  timestamp := time.Now()
  for _, instance := range instanceList {
    
    _, err := db.Exec("INSERT INTO instances (timestamp, instanceid, name, status, createdat, flavorid, tenantid, hostid, userid) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)",
      timestamp, instance.ID, instance.Name, instance.Status, instance.Created, instance.Flavor["id"], instance.TenantID, instance.HostID, instance.UserID)
    if err != nil {
      fmt.Println("Error inserting instance data:", err)
    }
  }
}
