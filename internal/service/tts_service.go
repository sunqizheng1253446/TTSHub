package service

import (
	"context"
	"ttshub/internal/adapter"
	"ttshub/internal/models"
)

// TTSService TTS服务接口
type TTSService interface {
	// Synthesize 执行文本转语音
	Synthesize(ctx context.Context, request *models.TTSRequest) (*models.TTSResponse, error)
	
	// ListChannels 获取所有渠道配置
	ListChannels() ([]*models.ChannelConfig, error)
	
	// GetChannel 获取指定渠道配置
	GetChannel(channelID uint) (*models.ChannelConfig, error)
	
	// UpdateChannel 更新渠道配置
	UpdateChannel(channel *models.ChannelConfig) error
	
	// DeleteChannel 删除渠道配置
	DeleteChannel(channelID uint) error
	
	// CreateChannel 创建新渠道配置
	CreateChannel(channel *models.ChannelConfig) error
	
	// RegisterAdapter 注册渠道适配器
	RegisterAdapter(adapterType string, adapterFactory adapter.AdapterFactory) error
	
	// ValidateChannelConfig 验证渠道配置
	ValidateChannelConfig(channelType string, config string) error
	
	// GetAvailableAdapterTypes 获取可用的适配器类型
	GetAvailableAdapterTypes() []string
	
	// Ping 测试服务可用性
	Ping() error
}