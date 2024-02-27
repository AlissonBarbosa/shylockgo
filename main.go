package main

import (
	"net/http"

	"github.com/AlissonBarbosa/shylockgo/common"
	"github.com/AlissonBarbosa/shylockgo/controllers"
	"github.com/gin-gonic/gin"
  "github.com/gin-contrib/cors"
)

func main()  {
  router := gin.Default()
  config := cors.DefaultConfig()
  config.AllowOrigins = []string{"*"}
  router.Use(cors.New(config))

  router.GET("/quota-summary", func (c *gin.Context)  {
    provider, err := common.GetProvider()
    if err != nil {
      c.JSON(http.StatusInternalServerError, gin.H{"error": "Error getting provider"})
      return
    }

    sponsorSummary, err := project.GetSponsorSummary(provider)
    if err != nil {
      c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
      return
    }

    c.JSON(http.StatusOK, sponsorSummary)
    
  })
  
  router.Run("0.0.0.0:8080")

}
