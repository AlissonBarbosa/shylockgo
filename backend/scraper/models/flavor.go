package models

type FlavorDesc struct {
  ID uint `gorm:"primaryKey"`
  Timestamp int64 `json:"timestamp"`
  FlavorID string `json:"flavor_id"`
  FlavorName string `json:"flavor_name"`
}

type FlavorSpec struct {
  ID uint `gorm:"primaryKey"`
  Timestamp int64 `json:"timestamp"`
  FlavorID string `json:"flavor_id"`
  Vcpu string `json:"vcpu"`
  Ram string `json:"ram"`
  Disk string `json:"disk"`
}
