package adapter

import (
	"errors"
	"ttshub/internal/models"
)

// 定义适配器相关错误
var (
	ErrInvalidConfig   = errors.New("invalid adapter configuration")
	ErrInvalidRequest  = errors.New("invalid TTS request")
	ErrConvertFailed   = errors.New("request conversion failed")
	ErrUnsupportedFormat = errors.New("unsupported audio format")
)

// TTSAdapter TTS适配器接口
type TTSAdapter interface {
	// ConvertRequest 将渠道特定请求转换为OpenAI TTS请求
	ConvertRequest(request *models.TTSRequest) (*models.OpenAIRequest, error)
	
	// GetChannelName 获取渠道名称
	GetChannelName() string
	
	// ValidateConfig 验证渠道配置
	ValidateConfig() error
	
	// Init 初始化适配器
	Init(config string) error
}

// BaseAdapter 适配器基类，提供通用实现
type BaseAdapter struct {
	ChannelName string
	Config      string
}

// GetChannelName 获取渠道名称
func (b *BaseAdapter) GetChannelName() string {
	return b.ChannelName
}

// Init 初始化适配器
func (b *BaseAdapter) Init(config string) error {
	b.Config = config
	return nil
}

// ValidateConfig 基础配置验证（可被子类重写）
func (b *BaseAdapter) ValidateConfig() error {
	if b.Config == "" {
		return ErrInvalidConfig
	}
	return nil
}