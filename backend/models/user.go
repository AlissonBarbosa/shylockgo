package models

type UserDesc struct {
  ID uint `gorm:"primaryKey"`
  Timestamp int64 `json:"timestamp"`
  UserID string `json:"user_id"`
  UserName string `json:"user_name"`
}

type UserProject struct {
  ID uint `gorm:"primaryKey"`
  Timestamp int64 `json:"timestamp"`
  UserID string `json:"user_id"`
  ProjectID string `json:"project_id"`
}
