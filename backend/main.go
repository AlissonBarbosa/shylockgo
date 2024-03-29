package main

import (
	"fmt"
	"net/http"

	"github.com/AlissonBarbosa/shylockgo/common"
	"github.com/AlissonBarbosa/shylockgo/controllers"
	"github.com/AlissonBarbosa/shylockgo/models"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
)

func main()  {
  models.ConnectDatabase()

  provider, err := common.GetProvider()
  if err != nil {
    fmt.Println("Error getting provider")
    return
  }

  c := cron.New()
  _, err = c.AddFunc("*/50 * * * *", func(){
    project.SaveProjectSummary(provider)
  })
  if err != nil {
    fmt.Println("Error adding task", err)
    return
  }
  c.Start()

  router := gin.Default()
  config := cors.DefaultConfig()
  config.AllowOrigins = []string{"*"}
  config.AllowHeaders = []string{"*"}
  router.Use(cors.New(config))

  router.GET("/quota-summary", func (c *gin.Context)  {
    var reports []models.ReportSummary

    models.DB.Table("(SELECT sponsor, MAX(timestamp) AS latest_timestamp FROM usage_reports GROUP BY sponsor) AS subquery").
      Select("subquery.latest_timestamp, subquery.sponsor, SUM(v_cpu_quota) AS v_cpu_quota, SUM(v_cpu_usage) AS v_cpu_usage, SUM(ram_quota) AS ram_quota, SUM(ram_usage) AS ram_usage").
      Joins("JOIN usage_reports ON usage_reports.sponsor = subquery.sponsor AND usage_reports.timestamp = subquery.latest_timestamp").
      Group("subquery.sponsor, subquery.latest_timestamp").Scan(&reports)

    aggregatedReports := make(map[string][]models.ReportSummary)

    for _, report := range reports {
      aggregatedReports[report.Sponsor] = append(aggregatedReports[report.Sponsor], report)
    }

    //sponsorSummary, err := project.GetSponsorSummary(provider)
    //if err != nil {
    //  c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
    //  return
    //}

    c.JSON(http.StatusOK, aggregatedReports)
    
  })
  
  router.Run("0.0.0.0:5000")

}
