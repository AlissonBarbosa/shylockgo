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

type ProjectDesc struct {
  ID uint `gorm:"primaryKey"`
  Timestamp int64 `json:"timestamp"`
  ProjectID string `json:"project_id"`
  ProjectName string `json:"project_name"`
  ProjectSponsor string `json:"project_sponsor"`
}

type ProjectQuota struct {
  ID uint `gorm:"primaryKey"`
  Timestamp int64 `json:"timestamp"`
  ProjectID string `json:"project_id"`
  QuotaRam int64 `json:"quota_ram"`
  QuotaVcpu int64 `json:"quota_vcpu"` 
  QuotaBlockstorage int64 `json:"quota_blockstorage"` 
}

type ProjectQuotaUsage struct {
  ID uint `gorm:"primaryKey"`
  Timestamp int64 `json:"timestamp"`
  ProjectID string `json:"project_id"`
  RamUsage int64 `json:"ram_usage"`
  VcpuUsage int64 `json:"vcpu_usage"` 
  BlockstorageUsage int64 `json:"blockstorage_usage"` 
}

type ProjectUsers struct {
  ID uint `gorm:"primaryKey"`
  Timestamp int64 `json:"timestamp"`
  ProjectID string `json:"project_id"`
  UserID string `json:"user_id"` 
}
