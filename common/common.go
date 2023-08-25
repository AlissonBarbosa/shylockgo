package common

import (
  "os"
  "fmt"
	"database/sql"
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
)

func GetProvider() (*gophercloud.ProviderClient, error) {
  opts, err := openstack.AuthOptionsFromEnv()
	if err != nil {
    return nil, err
	}
	provider, err := openstack.AuthenticatedClient(opts)
	if err != nil {
    return nil, err
	}
  return provider, nil
}

func ConnectToDatabase() (*sql.DB, error) {
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

