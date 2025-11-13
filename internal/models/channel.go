package models

import "fmt"
import "time"

// ChannelConfig 渠道配置模型
type ChannelConfig struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
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

// Validate 验证渠道配置
func (c *ChannelConfig) Validate() error {
	if c.Name == "" {
		return NewError("渠道名称不能为空")
	}
	if c.Type == "" {
		return NewError("渠道类型不能为空")
	}
	if c.Config == "" {
		return NewError("渠道配置不能为空")
	}
	return nil
}

// NewError 创建新的错误
func NewError(msg string) error {
	return fmt.Errorf(msg)
}