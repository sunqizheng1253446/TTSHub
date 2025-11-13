package models

import "time"

// ChannelConfig 渠道配置模型
type ChannelConfig struct {
	ID        string    `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name" gorm:"not null"`
	Type      string    `json:"type" gorm:"not null"`
	Config    string    `json:"config" gorm:"not null"` // JSON格式的配置
	Status    int       `json:"status" gorm:"default:1"` // 1:启用, 0:禁用
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TableName 指定表名
func (ChannelConfig) TableName() string {
	return "channels"
}