package models

import (
  "time"
)

type UsageReport struct {
  ID uint `gorm:"primaryKey"`
  Timestamp time.Time `json:"timestamp"`
  Sponsor string `json:"sponsor"`
  ProjectName string `json:"project_name"`
  VCPUQuota int `json:"vcpu_quota"`
  VCPUUsage int `json:"vcpu_usage"`
  RAMQuota int `json:"ram_quota"`
  RAMUsage int `json:"ram_usage"`
}

type ReportSummary struct {
  Sponsor string `json:"sponsor"`
  Timestamp time.Time `json:"timestamp"`
  VCPUQuota int `json:"vcpu_quota"`
  VCPUUsage int `json:"vcpu_usage"`
  RAMQuota int `json:"ram_quota"`
  RAMUsage int `json:"ram_usage"`
}

type ProjectData struct {
  ID string
  Sponsor string
  Name string
}

type ProjectSumUsage struct {
  VcpuUsage int
  RAMUsage int
}
