package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/servers"
	_ "github.com/mattn/go-sqlite3"
)

func main()  {
  logFile, err := os.OpenFile(os.Getenv("LOGFILE"), os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
  if err != nil {
    fmt.Println("[ERROR] Error creating log file:", err)
    os.Exit(1)
  }
  defer logFile.Close()

  log.SetOutput(logFile)

  opts, err := openstack.AuthOptionsFromEnv()
  if err != nil {
    log.Println("[ERROR] Error reading environment variables:", err)
    os.Exit(1)
  }

  provider, err := openstack.AuthenticatedClient(opts)
  if err != nil {
    log.Println("[ERROR] Error creating authenticated client", err)
    os.Exit(1)
  }

  client, err := openstack.NewComputeV2(provider, gophercloud.EndpointOpts{
    Region: os.Getenv("OS_REGION_NAME"),
  })
  if err != nil {
    log.Println("[ERROR] Error creating ComputeV2 client:", err)
    os.Exit(1)
  }
  
  for {
    db, err := connectToDatabase()
    if err != nil {
      log.Println("[ERROR] Error connecting to database", err)
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
    log.Println("[INFO] Saving instances on database")
    saveInstancesToDatabase(client, db)
    timer, err := strconv.Atoi(os.Getenv("TIMER"))
    if err != nil {
      log.Println("[ERROR] Error converting timer to integer", err)
      os.Exit(1)
    }
    time.Sleep(time.Duration(timer) * time.Minute)
  }
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
    log.Println("[ERROR] Error listing instances:", err)
    return
  }
  
  instanceList, err := servers.ExtractServers(rows)

  if err != nil {
    log.Println("[ERROR] Error extracting instance list:", err)
    return
  }

  timestamp := time.Now()
  for _, instance := range instanceList {
    
    _, err := db.Exec("INSERT INTO instances (timestamp, instanceid, name, status, createdat, flavorid, tenantid, hostid, userid) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)",
      timestamp, instance.ID, instance.Name, instance.Status, instance.Created, instance.Flavor["id"], instance.TenantID, instance.HostID, instance.UserID)
    if err != nil {
      log.Println("[ERROR] Error inserting instance data:", err)
    }
  }
}
