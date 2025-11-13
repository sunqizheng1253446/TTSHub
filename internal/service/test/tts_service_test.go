package test

import (
	"errors"
	"testing"

	"ttshub/internal/adapter"
	"ttshub/internal/service"
	"github.com/stretchr/testify/assert"
)

// MockStorage 是一个模拟的存储实现，用于测试

type MockStorage struct {
	configs map[string]service.ChannelConfig
}

func NewMockStorage() *MockStorage {
	return &MockStorage{
		configs: make(map[string]service.ChannelConfig),
	}
}

func (m *MockStorage) CreateChannelConfig(config service.ChannelConfig) error {
	m.configs[config.ID] = config
	return nil
}

func (m *MockStorage) GetChannelConfig(id string) (service.ChannelConfig, error) {
	config, ok := m.configs[id]
	if !ok {
		return service.ChannelConfig{}, errors.New("channel not found")
	}
	return config, nil
}

func (m *MockStorage) UpdateChannelConfig(config service.ChannelConfig) error {
	if _, ok := m.configs[config.ID]; !ok {
		return errors.New("channel not found")
	}
	m.configs[config.ID] = config
	return nil
}

func (m *MockStorage) DeleteChannelConfig(id string) error {
	if _, ok := m.configs[id]; !ok {
		return errors.New("channel not found")
	}
	delete(m.configs, id)
	return nil
}

func (m *MockStorage) ListChannelConfigs() ([]service.ChannelConfig, error) {
	configs := make([]service.ChannelConfig, 0, len(m.configs))
	for _, config := range m.configs {
		configs = append(configs, config)
	}
	return configs, nil
}

// MockAdapterFactory 是一个模拟的适配器工厂函数
func MockAdapterFactory() adapter.Adapter {
	return &MockAdapter{}
}

// 测试TTS服务的功能
func TestTTSService(t *testing.T) {
	// 创建模拟存储
	mockStorage := NewMockStorage()

	// 创建适配器管理器并注册模拟适配器
	adapterManager := adapter.NewAdapterManager()
	adapterManager.RegisterAdapter("mock", MockAdapterFactory)

	// 初始化TTS服务
	ttsService, err := service.NewTTSService(mockStorage, adapterManager)
	assert.NoError(t, err, "初始化TTS服务失败")

	// 测试1: 创建渠道
	createConfig := service.ChannelConfig{
		ID:           "test-channel",
		Name:         "测试渠道",
		Type:         "mock",
		Status:       "active",
		APIKey:       "test-api-key",
		APISecret:    "test-api-secret",
		Endpoint:     "https://api.example.com",
		Model:        "test-model",
		DefaultVoice: "alloy",
		DefaultSpeed: 1.0,
		DefaultPitch: 1.0,
		ExtraConfigJSON: `{"timeout": 30}`,
	}

	err = ttsService.CreateChannel(createConfig)
	assert.NoError(t, err, "创建渠道失败")

	// 测试2: 获取渠道
	getConfig, err := ttsService.GetChannel("test-channel")
	assert.NoError(t, err, "获取渠道失败")
	assert.Equal(t, "test-channel", getConfig.ID)
	assert.Equal(t, "测试渠道", getConfig.Name)

	// 测试3: 列出所有渠道
	channels, err := ttsService.ListChannels()
	assert.NoError(t, err, "列出渠道失败")
	assert.Len(t, channels, 1, "渠道数量不匹配")

	// 测试4: 更新渠道
	createConfig.Name = "更新后的测试渠道"
	createConfig.Status = "inactive"
	err = ttsService.UpdateChannel(createConfig)
	assert.NoError(t, err, "更新渠道失败")

	updatedConfig, err := ttsService.GetChannel("test-channel")
	assert.NoError(t, err, "获取更新后的渠道失败")
	assert.Equal(t, "更新后的测试渠道", updatedConfig.Name)
	assert.Equal(t, "inactive", updatedConfig.Status)

	// 测试5: 验证渠道配置
	createConfig.Status = "active"
	err = ttsService.UpdateChannel(createConfig) // 先激活渠道
	assert.NoError(t, err, "激活渠道失败")

	err = ttsService.ValidateChannelConfig("test-channel")
	assert.NoError(t, err, "验证渠道配置失败")

	// 测试6: 获取可用的适配器类型
	adapterTypes := ttsService.GetAvailableAdapterTypes()
	assert.Contains(t, adapterTypes, "mock", "应该包含mock适配器类型")

	// 测试7: 删除渠道
	err = ttsService.DeleteChannel("test-channel")
	assert.NoError(t, err, "删除渠道失败")

	// 验证删除成功
	_, err = ttsService.GetChannel("test-channel")
	assert.Error(t, err, "应该找不到已删除的渠道")

	// 测试错误情况

	// 测试8: 创建不存在的适配器类型
	invalidConfig := service.ChannelConfig{
		ID:   "invalid-channel",
		Type: "non-existent",
	}
	err = ttsService.CreateChannel(invalidConfig)
	assert.Error(t, err, "创建不存在的适配器类型应该返回错误")

	// 测试9: 获取不存在的渠道
	_, err = ttsService.GetChannel("non-existent")
	assert.Error(t, err, "获取不存在的渠道应该返回错误")

	// 测试10: 更新不存在的渠道
	_, err = ttsService.UpdateChannel(service.ChannelConfig{ID: "non-existent"})
	assert.Error(t, err, "更新不存在的渠道应该返回错误")

	// 测试11: 删除不存在的渠道
	err = ttsService.DeleteChannel("non-existent")
	assert.Error(t, err, "删除不存在的渠道应该返回错误")

	// 测试12: 验证不存在的渠道配置
	err = ttsService.ValidateChannelConfig("non-existent")
	assert.Error(t, err, "验证不存在的渠道配置应该返回错误")
}

// 测试适配器注册功能
func TestTTSServiceAdapterRegistration(t *testing.T) {
	// 创建模拟存储
	mockStorage := NewMockStorage()

	// 创建适配器管理器并注册多个模拟适配器
	adapterManager := adapter.NewAdapterManager()
	adapterManager.RegisterAdapter("mock1", MockAdapterFactory)
	adapterManager.RegisterAdapter("mock2", MockAdapterFactory)

	// 初始化TTS服务
	ttsService, err := service.NewTTSService(mockStorage, adapterManager)
	assert.NoError(t, err, "初始化TTS服务失败")

	// 检查可用适配器类型
	adapterTypes := ttsService.GetAvailableAdapterTypes()
	assert.Contains(t, adapterTypes, "mock1")
	assert.Contains(t, adapterTypes, "mock2")
	assert.Len(t, adapterTypes, 2)

	// 测试注册已存在的适配器
	err = ttsService.RegisterAdapter("mock1", MockAdapterFactory)
	assert.Error(t, err, "注册已存在的适配器应该返回错误")

	// 测试注册新适配器
	err = ttsService.RegisterAdapter("mock3", MockAdapterFactory)
	assert.NoError(t, err, "注册新适配器失败")

	// 检查新注册的适配器
	updatedTypes := ttsService.GetAvailableAdapterTypes()
	assert.Contains(t, updatedTypes, "mock3")
	assert.Len(t, updatedTypes, 3)
}