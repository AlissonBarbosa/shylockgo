package models

type ServerDesc struct {
  ID uint `gorm:"primaryKey"`
  Timestamp int64 `json:"timestamp"`
  ServerID string `json:"server_id"`
  ServerName string `json:"server_name"`
  ServerAddress string `json:"server_address"`
}

type ServerSpec struct {
  ID uint `gorm:"primaryKey"`
  Timestamp int64 `json:"timestamp"`
  ServerID string `json:"server_id"`
  FlavorID string `json:"flavor_id"`
}

type ServerUsage struct {
  ID uint `gorm:"primaryKey"`
  Timestamp int64 `json:"timestamp"`
  ServerID string `json:"server_id"`
  VcpuUsage string `json:"vcpu_usage"`
  RamUsage string `json:"ram_usage"`
  Domain string `json:"domain"`
  HostID string `json:"compute_id"`
}

type ServerOwnership struct {
  ID uint `gorm:"primaryKey"`
  Timestamp int64 `json:"timestamp"`
  ServerID string `json:"server_id"`
  UserID string `json:"user_id"`
  ProjectID string `json:"project_id"`
}

type ServerData struct {
  ID string `json:"id"`
  Name string `json:"name"`
  ProjectID string `json:"project_id"`
  HostID string `json:"host_id"`
  Domain string `json:"domain"`
  MemoryUsage int64 `json:"memory_usage"`
}

type ServerMeta struct {
  ID uint `json:"id"`
  ServerID string `json:"server_id"`
  Name string `json:"server_name"`
  ProjectID string `json:"project_id"`
  HostID string `json:"host_id"`
  Domain string `json:"serverdomain"`
  MemoryUsage int64 `json:"memory_usage"`
}

type DomainQueryResponse struct {
  Status string `json:"status"`
    Data   struct {
        ResultType string `json:"resultType"`
        Result     []struct {
            Metric struct {
                Domain string `json:"domain"`
            } `json:"metric"`
        } `json:"result"`
    } `json:"data"`
}
