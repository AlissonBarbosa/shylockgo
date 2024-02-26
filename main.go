package main

import (
	"net/http"

	"github.com/AlissonBarbosa/shylockgo/common"
	"github.com/AlissonBarbosa/shylockgo/project"
	"github.com/gin-gonic/gin"
)

func main()  {
  router := gin.Default()

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
  
  router.Run("localhost:8080")

}
