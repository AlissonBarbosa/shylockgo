package main

import (
	"fmt"
	"os"
	"time"

	"github.com/AlissonBarbosa/shylockgo/common"
	"github.com/AlissonBarbosa/shylockgo/models"
	"github.com/AlissonBarbosa/shylockgo/project"
)

func main()  {
  provider, err := common.GetProvider()
  if err != nil {
    fmt.Println("[ERROR] Error getting provider")
    os.Exit(1)
  }
  projects, err := project.GetProjects(provider)
  if err != nil {
    fmt.Println(err)
  }

  aggregateReports := make(map[string]models.UsageReport)

  for _, projectData := range projects {
    quotas, err := project.GetProjectQuota(provider, projectData.ID)
    if err != nil {
      fmt.Println("[ERROR] Error getting project quota")
    }
    projectUsage, err := project.GetProjectUsage(provider, projectData.ID)
    if err != nil {
      fmt.Println("[ERROR] Error getting project usage")
    }
    timestamp, _ := time.Now().MarshalText()

    report, ok := aggregateReports[projectData.Sponsor]
    if !ok {
      report = models.UsageReport{Timestamp: string(timestamp), Sponsor: projectData.Sponsor}
    }

    report.VCPUQuota += quotas.Cores
    report.VCPUUsage += projectUsage.VcpuUsage
    report.RAMQuota += quotas.RAM
    report.RAMUsage += projectUsage.RAMUsage

    aggregateReports[projectData.Sponsor] = report
  }

  fmt.Printf("Timestamp;Sponsor;vCPUQuota;vCPUUsage;RAMQuota;RAMUsage\n")
  for sponsor, report := range aggregateReports {
    fmt.Printf("%s;%s;%d;%d;%d;%d\n", report.Timestamp, sponsor, report.VCPUQuota, report.VCPUUsage, report.RAMQuota, report.RAMUsage)
  }
}
