package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/AlissonBarbosa/shylockgo/common"
	"github.com/AlissonBarbosa/shylockgo/controllers"
	"github.com/AlissonBarbosa/shylockgo/models"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
)



func main()  {
  l := slog.New(slog.NewJSONHandler(os.Stdout, nil))
  slog.SetDefault(l)

  models.ConnectDatabase()

  provider, err := common.GetProvider()
  if err != nil {
    slog.Error("Error getting provider:", err)
    return
  }
  
  // Populate database on start
  controllers.SaveProjectsDesc(provider)
  controllers.SaveAllServers(provider)
  controllers.GetProjectQuota(provider)
  controllers.GetProjectUsage(provider)


  c := cron.New()

  _, err = c.AddFunc("*/30 * * * *", func(){
    controllers.SaveProjectsDesc(provider)
  })
  if err != nil {
    slog.Error("Error adding task", err)
    return
  }

  _, err = c.AddFunc("*/31 * * * *", func(){
    controllers.SaveAllServers(provider)
  })
  if err != nil {
    slog.Error("Error adding task", err)
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

    c.JSON(http.StatusOK, aggregatedReports)
    
  })

  router.GET("/servers", func (c *gin.Context) {
    var servers []models.ServerMeta
    servers, err = controllers.GetAllServers(provider)
    if err != nil {
      c.JSON(http.StatusInternalServerError, gin.H{
        "status": "error",
        "message": "Something wrong on server side",
      })
      return
    }
    c.JSON(http.StatusOK, gin.H{
      "status": "success",
      "data": servers,
      "message": "Success",
    })
  })
  
  router.Run("0.0.0.0:5000")

}
