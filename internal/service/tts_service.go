package service

import (
	"ttshub/internal/models"
)

// TTSService TTS服务接口
type TTSService interface {
	// Synthesize 执行文本转语音
	Synthesize(request *models.TTSRequest) (*models.TTSResponse, error)
	
	// ListChannels 获取所有渠道配置
	ListChannels() ([]*models.ChannelConfig, error)
	
	// GetChannel 获取指定渠道配置
	GetChannel(channelID string) (*models.ChannelConfig, error)
	
	// UpdateChannel 更新渠道配置
	UpdateChannel(channel *models.ChannelConfig) error
	
	// DeleteChannel 删除渠道配置
	DeleteChannel(channelID string) error
	
	// CreateChannel 创建新渠道配置
	CreateChannel(channel *models.ChannelConfig) error
	
	// RegisterAdapter 注册渠道适配器
	RegisterAdapter(adapterType string, adapterFactory func() interface{}) error
}