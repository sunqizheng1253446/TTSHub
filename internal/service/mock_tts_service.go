package service

import (
	"context"
	"errors"
	"time"
	"ttshub/internal/adapter"
	"ttshub/internal/models"
)

// MockTTSService 是TTSService的模拟实现，用于测试

type MockTTSService struct {
	channels map[string]*models.ChannelConfig
}

// NewMockTTSService 创建一个新的模拟TTS服务
func NewMockTTSService() *MockTTSService {
	return &MockTTSService{
		channels: make(map[string]*models.ChannelConfig),
	}
}

// Synthesize 模拟TTS合成
func (m *MockTTSService) Synthesize(ctx context.Context, req *models.TTSRequest) (*models.TTSResponse, error) {
	// 模拟处理时间
	time.Sleep(100 * time.Millisecond)
	
	// 创建模拟响应
	response := &models.TTSResponse{
		AudioData: []byte("模拟的音频数据"),
		Format:    "mp3",
		Duration:  2.5,
		Cost:      0.01,
		ChannelID: req.ChannelID,
	}
	
	return response, nil
}

// ListChannels 模拟列出所有渠道
func (m *MockTTSService) ListChannels() ([]*models.ChannelConfig, error) {
	channels := make([]*models.ChannelConfig, 0, len(m.channels))
	for _, channel := range m.channels {
		channels = append(channels, channel)
	}
	return channels, nil
}

// GetChannel 模拟获取单个渠道
func (m *MockTTSService) GetChannel(channelID uint) (*models.ChannelConfig, error) {
	idStr := string(channelID)
	channel, exists := m.channels[idStr]
	if !exists {
		return nil, errors.New("channel not found")
	}
	return channel, nil
}

// CreateChannel 模拟创建渠道
func (m *MockTTSService) CreateChannel(config *models.ChannelConfig) error {
	if config.ID == 0 {
		return errors.New("channel id is required")
	}
	m.channels[string(config.ID)] = config
	return nil
}

// UpdateChannel 模拟更新渠道
func (m *MockTTSService) UpdateChannel(config *models.ChannelConfig) error {
	if _, exists := m.channels[string(config.ID)]; !exists {
		return errors.New("channel not found")
	}
	m.channels[string(config.ID)] = config
	return nil
}

// DeleteChannel 模拟删除渠道
func (m *MockTTSService) DeleteChannel(channelID uint) error {
	idStr := string(channelID)
	if _, exists := m.channels[idStr]; !exists {
		return errors.New("channel not found")
	}
	delete(m.channels, idStr)
	return nil
}

// ValidateChannelConfig 模拟验证渠道配置
func (m *MockTTSService) ValidateChannelConfig(id string) error {
	if _, exists := m.channels[id]; !exists {
		return errors.New("channel not found")
	}
	return nil
}

// RegisterAdapter 模拟注册适配器
func (m *MockTTSService) RegisterAdapter(name string, factory adapter.AdapterFactory) error {
	return nil
}

// GetAvailableAdapterTypes 模拟获取可用的适配器类型
func (m *MockTTSService) GetAvailableAdapterTypes() []string {
	return []string{"mock", "openai", "baidu", "google", "azure"}
}