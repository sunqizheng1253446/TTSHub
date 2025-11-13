package test

import (
	"encoding/json"
	"os"
	"testing"

	"ttshub/internal/storage"
	"github.com/stretchr/testify/assert"
)

var testDBPath = "./test_ttshub.db"

func TestSQLiteStore(t *testing.T) {
	// 测试前清理可能存在的测试数据库文件
	os.Remove(testDBPath)
	defer os.Remove(testDBPath)

	// 初始化存储
	store, err := storage.NewSQLiteStore(testDBPath)
	assert.NoError(t, err, "初始化SQLite存储失败")

	// 测试1: 创建新渠道配置
	config := storage.ChannelConfig{
		ID:           "test-openai",
		Name:         "测试OpenAI渠道",
		Type:         "openai",
		Status:       "active",
		APIKey:       "test-api-key",
		APISecret:    "test-api-secret",
		Endpoint:     "https://api.openai.com/v1",
		Model:        "tts-1",
		DefaultVoice: "alloy",
		DefaultSpeed: 1.0,
		DefaultPitch: 1.0,
		ExtraConfig:  map[string]interface{}{"max_tokens": 1000},
	}

	// 将ExtraConfig序列化为JSON
	extraConfigJSON, _ := json.Marshal(config.ExtraConfig)
	config.ExtraConfigJSON = string(extraConfigJSON)

	err = store.CreateChannelConfig(config)
	assert.NoError(t, err, "创建渠道配置失败")

	// 测试2: 获取渠道配置
	retrievedConfig, err := store.GetChannelConfig("test-openai")
	assert.NoError(t, err, "获取渠道配置失败")
	assert.Equal(t, config.ID, retrievedConfig.ID)
	assert.Equal(t, config.Name, retrievedConfig.Name)
	assert.Equal(t, config.Type, retrievedConfig.Type)
	assert.Equal(t, config.Status, retrievedConfig.Status)
	assert.Equal(t, config.APIKey, retrievedConfig.APIKey)
	assert.Equal(t, config.APISecret, retrievedConfig.APISecret)
	assert.Equal(t, config.Endpoint, retrievedConfig.Endpoint)
	assert.Equal(t, config.Model, retrievedConfig.Model)
	assert.Equal(t, config.DefaultVoice, retrievedConfig.DefaultVoice)
	assert.Equal(t, config.DefaultSpeed, retrievedConfig.DefaultSpeed)
	assert.Equal(t, config.DefaultPitch, retrievedConfig.DefaultPitch)

	// 测试3: 列出所有渠道配置
	configs, err := store.ListChannelConfigs()
	assert.NoError(t, err, "列出渠道配置失败")
	assert.Len(t, configs, 1, "渠道配置数量不匹配")

	// 测试4: 更新渠道配置
	config.Name = "更新后的测试渠道"
	config.Status = "inactive"
	config.DefaultSpeed = 0.8
	err = store.UpdateChannelConfig(config)
	assert.NoError(t, err, "更新渠道配置失败")

	updatedConfig, err := store.GetChannelConfig("test-openai")
	assert.NoError(t, err, "获取更新后的渠道配置失败")
	assert.Equal(t, "更新后的测试渠道", updatedConfig.Name)
	assert.Equal(t, "inactive", updatedConfig.Status)
	assert.Equal(t, 0.8, updatedConfig.DefaultSpeed)

	// 测试5: 删除渠道配置
	err = store.DeleteChannelConfig("test-openai")
	assert.NoError(t, err, "删除渠道配置失败")

	// 验证删除成功
	_, err = store.GetChannelConfig("test-openai")
	assert.Error(t, err, "应该找不到已删除的渠道配置")

	// 测试6: 再次列出渠道配置，应该为空
	configs, err = store.ListChannelConfigs()
	assert.NoError(t, err, "列出渠道配置失败")
	assert.Len(t, configs, 0, "渠道配置应该已清空")
}

func TestSQLiteStoreErrorCases(t *testing.T) {
	// 测试前清理可能存在的测试数据库文件
	os.Remove(testDBPath)
	defer os.Remove(testDBPath)

	// 初始化存储
	store, err := storage.NewSQLiteStore(testDBPath)
	assert.NoError(t, err, "初始化SQLite存储失败")

	// 测试: 获取不存在的渠道配置
	_, err = store.GetChannelConfig("non-existent")
	assert.Error(t, err, "获取不存在的渠道配置应该返回错误")

	// 测试: 更新不存在的渠道配置
	config := storage.ChannelConfig{
		ID: "non-existent",
	}
	err = store.UpdateChannelConfig(config)
	assert.Error(t, err, "更新不存在的渠道配置应该返回错误")

	// 测试: 删除不存在的渠道配置
	err = store.DeleteChannelConfig("non-existent")
	assert.Error(t, err, "删除不存在的渠道配置应该返回错误")
}