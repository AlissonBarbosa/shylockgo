package models

import (
	//"fmt"
	//"os"

	"gorm.io/gorm"
  //"gorm.io/driver/postgres"
  "github.com/glebarez/sqlite"
)

var DB *gorm.DB

func ConnectDatabase() {
// dsn := fmt.Sprintf("host=db user=%s password=%s dbname=%s port=5432 sslmode=disable TimeZone=America/Recife",
//   os.Getenv("POSTGRES_USER"), os.Getenv("POSTGRES_PASSWORD"), os.Getenv("POSTGRES_DB"))
// database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
  database, err := gorm.Open(sqlite.Open("/var/lib/shylockgo/shylock.db"), &gorm.Config{})

  if err != nil {
    panic("Failed to connect to database")
  }

  database.AutoMigrate(&ProjectDesc{})
  database.AutoMigrate(&ProjectQuota{})
  database.AutoMigrate(&ProjectQuotaUsage{})
  database.AutoMigrate(&ServerDesc{})
  database.AutoMigrate(&ServerSpec{})
  database.AutoMigrate(&ServerUsage{})
  database.AutoMigrate(&ServerOwnership{})
  database.AutoMigrate(&FlavorDesc{})
  database.AutoMigrate(&FlavorSpec{})
  database.AutoMigrate(&UserDesc{})
  database.AutoMigrate(&UserProject{})

  DB = database
}
