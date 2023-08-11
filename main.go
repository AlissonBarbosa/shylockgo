package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/servers"
	"github.com/gophercloud/gophercloud/openstack/identity/v3/projects"
	_ "github.com/lib/pq"
	"github.com/robfig/cron/v3"
)

func main() {
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

  c := cron.New()

  db, err := connectToDatabase()
    if err != nil {
      log.Println("[ERROR] Error connecting to database", err)
      return
    }
  defer db.Close()

  saveInstancesToDatabase(provider, db)
  saveProjectsToDatabase(provider, db)

  _, err = c.AddFunc("@every "+os.Getenv("TIMER")+"m", func() {
    saveInstancesToDatabase(provider, db)
  })

  if err != nil {
    log.Println("[ERROR] Error scheduling instance job:", err)
    os.Exit(1)
  }

  _, err = c.AddFunc("@hourly", func() {
    saveProjectsToDatabase(provider, db)
  })

  if err != nil {
    log.Println("[Error] Error scheduling project job")
  }

  c.Start()
  select {}
}

func connectToDatabase() (*sql.DB, error) {
	dbInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		os.Getenv("DBHOST"), 5432, os.Getenv("DBUSER"), os.Getenv("DBPASSWORD"), os.Getenv("DBNAME"))
	db, err := sql.Open("postgres", dbInfo)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}

func saveInstancesToDatabase(provider *gophercloud.ProviderClient, db *sql.DB) {
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
    log.Println("[ERROR] Error writing table instances:", err)
    return
  }

  client, err := openstack.NewComputeV2(provider, gophercloud.EndpointOpts{
		Region: os.Getenv("OS_REGION_NAME"),
	})
	if err != nil {
		log.Println("[ERROR] Error creating ComputeV2 client:", err)
		os.Exit(1)
	}
  
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
		_, err := db.Exec("INSERT INTO instances (id, timestamp, epoch, instanceid, name, status, createdat, flavorid, tenantid, hostid, userid) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)",
			uuid.New(), timestamp, timestamp.Unix() ,instance.ID, instance.Name, instance.Status, instance.Created, instance.Flavor["id"], instance.TenantID, instance.HostID, instance.UserID)
		if err != nil {
			log.Println("[ERROR] Error inserting instance data:", err)
		}
	}
  log.Println("[INFO] Instances saved on database")
}

func saveProjectsToDatabase(provider *gophercloud.ProviderClient, db *sql.DB) {
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
    log.Println("[ERROR] Error writing table projects", err)
    return
  }
  client, err := openstack.NewIdentityV3(provider, gophercloud.EndpointOpts{})
  if err != nil {
    log.Println("[Error] Error listing projects", err)
    return
  }

  listOpts := projects.ListOpts{
    Enabled: gophercloud.Enabled,
  }
  rows, err := projects.List(client, listOpts).AllPages()
  if err != nil {
    log.Println("[ERROR] Error listing projects:", err)
    return
  }

  projectList, err := projects.ExtractProjects(rows)
  if err != nil {
    log.Println("[ERROR] Error extracting projects:", err)
    return
  }

  timestamp := time.Now()
  for _, project := range projectList {
    _, err := db.Exec("INSERT INTO projects (id, timestamp, epoch, projectid, name, description, domainid) VALUES ($1, $2, $3, $4, $5, $6, $7)",
      uuid.New(), timestamp, timestamp.Unix(), project.ID, project.Name, project.Description, project.DomainID)
    if err != nil {
      log.Println("[ERROR] Error inserting project data:", err)
    }
  }
  log.Println("[INFO] Projects saved on database")
}
