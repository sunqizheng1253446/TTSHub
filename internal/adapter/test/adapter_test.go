package test

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"ttshub/internal/adapter"
	"github.com/stretchr/testify/assert"
)

// MockHTTPResponse 创建模拟HTTP响应
func MockHTTPResponse(status int, body []byte) *httptest.Server {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(status)
		w.Write(body)
	})
	return httptest.NewServer(handler)
}

// MockAdapter 实现一个测试用的模拟适配器
func MockAdapter() adapter.Adapter {
	return &mockAdapter{}
}

// mockAdapter 是一个测试用的适配器实现

type mockAdapter struct{}

type mockConfig struct {
	APIKey string `json:"api_key"`
}

func (m *mockAdapter) Init(config map[string]interface{}) error {
	// 简单的初始化逻辑
	return nil
}

func (m *mockAdapter) ValidateConfig(config map[string]interface{}) error {
	// 验证配置中是否包含APIKey
	if _, ok := config["api_key"]; !ok {
		return errors.New("missing required field: api_key")
	}
	return nil
}

func (m *mockAdapter) ConvertRequest(req *adapter.TTSRequest) (adapter.ProviderRequest, error) {
	// 转换请求为模拟的提供者请求
	return adapter.ProviderRequest{
		Method:  "POST",
		URL:     "https://mock-api.example.com/tts",
		Headers: map[string]string{"Content-Type": "application/json"},
		Body:    []byte(`{"text":"` + req.Text + `"}`),
	}, nil
}

func (m *mockAdapter) ConvertResponse(resp *http.Response) (*adapter.TTSResponse, error) {
	// 从响应中读取数据
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// 创建TTS响应
	return &adapter.TTSResponse{
		AudioData: body,
		Format:    "mp3",
		Duration:  1000,
	}, nil
}

func (m *mockAdapter) GetVoiceMapping() map[string]string {
	return map[string]string{
		"alloy":  "mock_alloy",
		"echo":   "mock_echo",
		"fable":  "mock_fable",
		"onyx":   "mock_onyx",
		"nova":   "mock_nova",
		"shimmer": "mock_shimmer",
	}
}

func (m *mockAdapter) GetAvailableVoices() []adapter.Voice {
	return []adapter.Voice{
		{ID: "mock_alloy", Name: "Mock Alloy", Gender: "female", Locale: "en-US"},
		{ID: "mock_echo", Name: "Mock Echo", Gender: "male", Locale: "en-US"},
	}
}

// 测试适配器管理器的功能
func TestAdapterManager(t *testing.T) {
	// 初始化适配器管理器
	manager := adapter.NewAdapterManager()

	// 测试1: 注册适配器
	err := manager.RegisterAdapter("mock", MockAdapter)
	assert.NoError(t, err, "注册适配器失败")

	// 测试2: 获取已注册的适配器
	adapterFactory, exists := manager.GetAdapter("mock")
	assert.True(t, exists, "应该找到已注册的适配器")

	// 测试3: 创建适配器实例
	adapterInstance, err := adapterFactory()
	assert.NoError(t, err, "创建适配器实例失败")

	// 测试4: 初始化适配器
	config := map[string]interface{}{"api_key": "test-key"}
	err = adapterInstance.Init(config)
	assert.NoError(t, err, "初始化适配器失败")

	// 测试5: 验证配置
	err = adapterInstance.ValidateConfig(config)
	assert.NoError(t, err, "验证有效配置失败")

	// 测试6: 验证无效配置
	invalidConfig := map[string]interface{}{}
	err = adapterInstance.ValidateConfig(invalidConfig)
	assert.Error(t, err, "验证无效配置应该返回错误")

	// 测试7: 转换请求
	ttsReq := &adapter.TTSRequest{
		Text:  "这是一段测试文本",
		Voice: "alloy",
		Speed: 1.0,
		Pitch: 1.0,
	}
	providerReq, err := adapterInstance.ConvertRequest(ttsReq)
	assert.NoError(t, err, "转换请求失败")
	assert.Equal(t, "POST", providerReq.Method)
	assert.Equal(t, "https://mock-api.example.com/tts", providerReq.URL)

	// 测试8: 获取语音映射
	voiceMapping := adapterInstance.GetVoiceMapping()
	assert.Len(t, voiceMapping, 6, "语音映射数量不匹配")
	assert.Equal(t, "mock_alloy", voiceMapping["alloy"])

	// 测试9: 获取可用语音
	voices := adapterInstance.GetAvailableVoices()
	assert.Len(t, voices, 2, "可用语音数量不匹配")
	assert.Equal(t, "Mock Alloy", voices[0].Name)

	// 测试10: 获取未注册的适配器
	_, exists = manager.GetAdapter("non-existent")
	assert.False(t, exists, "不应该找到未注册的适配器")
}

// 测试HTTP客户端的功能
func TestHTTPClient(t *testing.T) {
	// 创建模拟HTTP服务器，返回成功响应
	mockServer := MockHTTPResponse(http.StatusOK, []byte("mock audio data"))
	defer mockServer.Close()

	// 初始化HTTP客户端
	client := adapter.NewHTTPClient(10)

	// 创建请求
	req := adapter.ProviderRequest{
		Method:  "GET",
		URL:     mockServer.URL,
		Headers: map[string]string{},
		Body:    nil,
	}

	// 发送请求
	resp, err := client.SendRequest(req)
	assert.NoError(t, err, "发送HTTP请求失败")
	assert.NotNil(t, resp, "响应应该不为空")
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// 读取响应体
	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err, "读取响应体失败")
	assert.Equal(t, "mock audio data", string(body))

	// 测试失败的HTTP请求
	failedServer := MockHTTPResponse(http.StatusInternalServerError, []byte("error"))
	defer failedServer.Close()

	failedReq := adapter.ProviderRequest{
		Method:  "GET",
		URL:     failedServer.URL,
		Headers: map[string]string{},
		Body:    nil,
	}

	failedResp, err := client.SendRequest(failedReq)
	assert.NoError(t, err, "即使状态码不是200，也不应该返回错误")
	assert.NotNil(t, failedResp, "即使失败，响应也应该不为空")
	assert.Equal(t, http.StatusInternalServerError, failedResp.StatusCode)
}