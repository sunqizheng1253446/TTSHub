package service

import (
	"errors"
)

// MockTTSService 是TTSService的模拟实现，用于测试

type MockTTSService struct {
	channels map[string]ChannelConfig
}

// NewMockTTSService 创建一个新的模拟TTS服务
func NewMockTTSService() *MockTTSService {
	return &MockTTSService{
		channels: make(map[string]ChannelConfig),
	}
}

// Synthesize 模拟文本转语音功能
func (m *MockTTSService) Synthesize(req *SynthesizeRequest) (*SynthesizeResponse, error) {
	// 检查渠道是否存在
	if _, exists := m.channels[req.Channel]; !exists {
		return nil, errors.New("channel not found")
	}

	// 返回模拟的响应
	return &SynthesizeResponse{
		AudioData: []byte("mock audio data"),
		Format:    "mp3",
		Duration:  1000,
	}, nil
}

// ListChannels 模拟列出所有渠道
func (m *MockTTSService) ListChannels() ([]ChannelConfig, error) {
	channels := make([]ChannelConfig, 0, len(m.channels))
	for _, channel := range m.channels {
		channels = append(channels, channel)
	}
	return channels, nil
}

// GetChannel 模拟获取单个渠道
func (m *MockTTSService) GetChannel(id string) (ChannelConfig, error) {
	channel, exists := m.channels[id]
	if !exists {
		return ChannelConfig{}, errors.New("channel not found")
	}
	return channel, nil
}

// CreateChannel 模拟创建渠道
func (m *MockTTSService) CreateChannel(config ChannelConfig) error {
	if config.ID == "" {
		return errors.New("channel id is required")
	}
	m.channels[config.ID] = config
	return nil
}

// UpdateChannel 模拟更新渠道
func (m *MockTTSService) UpdateChannel(config ChannelConfig) error {
	if _, exists := m.channels[config.ID]; !exists {
		return errors.New("channel not found")
	}
	m.channels[config.ID] = config
	return nil
}

// DeleteChannel 模拟删除渠道
func (m *MockTTSService) DeleteChannel(id string) error {
	if _, exists := m.channels[id]; !exists {
		return errors.New("channel not found")
	}
	delete(m.channels, id)
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
func (m *MockTTSService) RegisterAdapter(name string, factory AdapterFactory) error {
	return nil
}

// GetAvailableAdapterTypes 模拟获取可用的适配器类型
func (m *MockTTSService) GetAvailableAdapterTypes() []string {
	return []string{"mock", "openai", "baidu", "google", "azure"}
}