package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"ttshub/internal/handlers"
	"ttshub/internal/service"
	"github.com/stretchr/testify/assert"
)

// 创建模拟TTS服务
func NewMockTTSService() *service.MockTTSService {
	return service.NewMockTTSService()
}

// 测试健康检查接口
func TestHealthCheck(t *testing.T) {
	// 创建模拟服务和处理器
	mockService := NewMockTTSService()
	handler := handlers.NewTTSHandler(mockService)

	// 创建HTTP请求
	req, err := http.NewRequest("GET", "/api/health", nil)
	assert.NoError(t, err, "创建健康检查请求失败")

	// 创建响应记录器
	recorder := httptest.NewRecorder()

	// 调用处理器
	handler.Health(recorder, req)

	// 验证响应
	resp := recorder.Result()
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// 解析响应体
	var response map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	assert.NoError(t, err, "解析健康检查响应失败")
	assert.Equal(t, "ok", response["status"])
}

// 测试Ping接口
func TestPing(t *testing.T) {
	// 创建模拟服务和处理器
	mockService := NewMockTTSService()
	handler := handlers.NewTTSHandler(mockService)

	// 创建HTTP请求
	req, err := http.NewRequest("GET", "/api/ping", nil)
	assert.NoError(t, err, "创建Ping请求失败")

	// 创建响应记录器
	recorder := httptest.NewRecorder()

	// 调用处理器
	handler.Ping(recorder, req)

	// 验证响应
	resp := recorder.Result()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "pong", recorder.Body.String())
}

// 测试获取适配器类型接口
func TestGetAvailableAdapterTypes(t *testing.T) {
	// 创建模拟服务和处理器
	mockService := NewMockTTSService()
	handler := handlers.NewTTSHandler(mockService)

	// 创建HTTP请求
	req, err := http.NewRequest("GET", "/api/v1/adapters", nil)
	assert.NoError(t, err, "创建获取适配器类型请求失败")

	// 创建响应记录器
	recorder := httptest.NewRecorder()

	// 调用处理器
	handler.GetAvailableAdapterTypes(recorder, req)

	// 验证响应
	resp := recorder.Result()
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// 解析响应体
	var response []string
	err = json.NewDecoder(resp.Body).Decode(&response)
	assert.NoError(t, err, "解析适配器类型响应失败")
	assert.NotEmpty(t, response)
}

// 测试渠道管理接口
func TestChannelManagement(t *testing.T) {
	// 创建模拟服务和处理器
	mockService := NewMockTTSService()
	handler := handlers.NewTTSHandler(mockService)

	// 测试1: 创建渠道
	channelConfig := map[string]interface{}{
		"id":            "test-channel",
		"name":          "测试渠道",
		"type":          "mock",
		"status":        "active",
		"api_key":       "test-api-key",
		"api_secret":    "test-api-secret",
		"endpoint":      "https://api.example.com",
		"model":         "test-model",
		"default_voice": "alloy",
		"default_speed": 1.0,
		"default_pitch": 1.0,
	}

	jsonData, _ := json.Marshal(channelConfig)
	req, _ := http.NewRequest("POST", "/api/v1/channels", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()
	handler.CreateChannel(recorder, req)

	resp := recorder.Result()
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	// 测试2: 获取渠道
	req, _ = http.NewRequest("GET", "/api/v1/channels/test-channel", nil)
	recorder = httptest.NewRecorder()
	handler.GetChannel(recorder, req)

	resp = recorder.Result()
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// 解析响应体
	var getResponse map[string]interface{}
	err := json.NewDecoder(resp.Body).Decode(&getResponse)
	assert.NoError(t, err, "解析获取渠道响应失败")
	assert.Equal(t, "test-channel", getResponse["id"])

	// 测试3: 列出所有渠道
	req, _ = http.NewRequest("GET", "/api/v1/channels", nil)
	recorder = httptest.NewRecorder()
	handler.ListChannels(recorder, req)

	resp = recorder.Result()
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// 解析响应体
	var listResponse []map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&listResponse)
	assert.NoError(t, err, "解析列出渠道响应失败")
	assert.Len(t, listResponse, 1)

	// 测试4: 更新渠道
	channelConfig["name"] = "更新后的测试渠道"
	jsonData, _ = json.Marshal(channelConfig)
	req, _ = http.NewRequest("PUT", "/api/v1/channels/test-channel", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	recorder = httptest.NewRecorder()
	handler.UpdateChannel(recorder, req)

	resp = recorder.Result()
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// 测试5: 删除渠道
	req, _ = http.NewRequest("DELETE", "/api/v1/channels/test-channel", nil)
	recorder = httptest.NewRecorder()
	handler.DeleteChannel(recorder, req)

	resp = recorder.Result()
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// 测试6: 获取不存在的渠道
	req, _ = http.NewRequest("GET", "/api/v1/channels/non-existent", nil)
	recorder = httptest.NewRecorder()
	handler.GetChannel(recorder, req)

	resp = recorder.Result()
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}

// 测试TTS合成接口
func TestSynthesize(t *testing.T) {
	// 创建模拟服务和处理器
	mockService := NewMockTTSService()
	handler := handlers.NewTTSHandler(mockService)

	// 首先创建一个渠道
	channelConfig := map[string]interface{}{
		"id":            "test-channel",
		"name":          "测试渠道",
		"type":          "mock",
		"status":        "active",
		"api_key":       "test-api-key",
		"api_secret":    "test-api-secret",
		"endpoint":      "https://api.example.com",
		"model":         "test-model",
		"default_voice": "alloy",
		"default_speed": 1.0,
		"default_pitch": 1.0,
	}

	jsonData, _ := json.Marshal(channelConfig)
	req, _ := http.NewRequest("POST", "/api/v1/channels", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()
	handler.CreateChannel(recorder, req)

	// 测试合成接口
	synthesizeReq := map[string]interface{}{
		"channel": "test-channel",
		"text":    "这是一段测试文本",
		"voice":   "alloy",
		"speed":   1.0,
		"pitch":   1.0,
	}

	jsonData, _ = json.Marshal(synthesizeReq)
	req, _ = http.NewRequest("POST", "/api/v1/tts", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	recorder = httptest.NewRecorder()
	handler.Synthesize(recorder, req)

	resp := recorder.Result()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "audio/mpeg", resp.Header.Get("Content-Type"))
}

// 测试错误处理
func TestErrorHandling(t *testing.T) {
	// 创建模拟服务和处理器
	mockService := NewMockTTSService()
	handler := handlers.NewTTSHandler(mockService)

	// 测试无效的JSON请求
	req, _ := http.NewRequest("POST", "/api/v1/channels", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()
	handler.CreateChannel(recorder, req)

	resp := recorder.Result()
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	// 测试缺少必要字段的请求
	invalidConfig := map[string]interface{}{
		"name": "缺少ID的渠道",
	}

	jsonData, _ := json.Marshal(invalidConfig)
	req, _ = http.NewRequest("POST", "/api/v1/channels", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	recorder = httptest.NewRecorder()
	handler.CreateChannel(recorder, req)

	resp = recorder.Result()
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}