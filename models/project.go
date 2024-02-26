package models

type UsageReport struct {
  Timestamp string
  Sponsor string
  VCPUQuota int
  VCPUUsage int
  RAMQuota int
  RAMUsage int
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
