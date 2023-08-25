package main

import (
	"fmt"
	"log"
	"os"

	"github.com/AlissonBarbosa/shylockgo/common"
	"github.com/AlissonBarbosa/shylockgo/project"
	"github.com/AlissonBarbosa/shylockgo/server"
	_ "github.com/lib/pq"
)

func main() {

  if len(os.Args) != 2 {
    fmt.Println("Usage: shylock <function>")
    os.Exit(0)
  }

  function := os.Args[1]
  
	logFile, err := os.OpenFile(os.Getenv("LOGFILE"), os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("[ERROR] Error creating log file:", err)
		os.Exit(1)
	}
	defer logFile.Close()

	log.SetOutput(logFile)

  provider, err := common.GetProvider()
  if err != nil {
    log.Println("[ERROR] Error getting provider")
    os.Exit(1)
  }

  db, err := common.ConnectToDatabase()
    if err != nil {
      log.Println("[ERROR] Error connecting to database", err)
      os.Exit(1)
    }
  defer db.Close()

  switch function {
  case "instances":
    err := server.SaveInstancesToDatabase(provider, db)
    if err != nil {
      log.Println(err)
    }
  case "projects":
    err := project.SaveProjectsToDatabase(provider, db)
    if err != nil {
      log.Println(err)
    }
  default:
    fmt.Println("Invalid function.")
  }
}
