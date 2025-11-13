package adapter

import (
	"encoding/json"
	"fmt"
	"ttshub/internal/models"
)

// CustomChannelConfig 自定义渠道配置
type CustomChannelConfig struct {
	BaseURL    string                 `json:"base_url"`
	APIKey     string                 `json:"api_key"`
	AuthType   string                 `json:"auth_type"` // bearer, api_key, basic
	Headers    map[string]string      `json:"headers"`
	Params     map[string]interface{} `json:"params"`
	Mapping    map[string]string      `json:"mapping"` // 字段映射关系
	Format     string                 `json:"format"`
	SampleRate int                    `json:"sample_rate"`
	Speed      float64                `json:"speed"`
	Voice      string                 `json:"voice"`
	Model      string                 `json:"model"`
}

// CustomChannelAdapter 自定义渠道适配器
type CustomChannelAdapter struct {
	BaseAdapter
	Config CustomChannelConfig
}

// NewCustomChannelAdapter 创建自定义渠道适配器实例
func NewCustomChannelAdapter(channelName string) *CustomChannelAdapter {
	return &CustomChannelAdapter{
		BaseAdapter: BaseAdapter{
			ChannelName: channelName,
		},
		Config: CustomChannelConfig{
			AuthType:   "bearer",
			Headers:    make(map[string]string),
			Params:     make(map[string]interface{}),
			Mapping:    make(map[string]string),
			Format:     "mp3",
			SampleRate: 24000,
			Speed:      1.0,
		},
	}
}

// Init 初始化适配器
func (a *CustomChannelAdapter) Init(config string) error {
	// 设置默认配置
	a.Config = CustomChannelConfig{
		AuthType:   "bearer",
		Headers:    make(map[string]string),
		Params:     make(map[string]interface{}),
		Mapping:    make(map[string]string),
		Format:     "mp3",
		SampleRate: 24000,
		Speed:      1.0,
	}

	// 如果提供了配置，则解析
	if config != "" {
		if err := json.Unmarshal([]byte(config), &a.Config); err != nil {
			return fmt.Errorf("解析配置失败: %w", err)
		}
	}

	return nil
}

// ValidateConfig 验证配置
func (a *CustomChannelAdapter) ValidateConfig() error {
	// 基础验证
	if err := a.BaseAdapter.ValidateConfig(); err != nil {
		return err
	}

	// 自定义渠道特定验证
	if a.Config.BaseURL == "" {
		return fmt.Errorf("基础URL不能为空")
	}

	// 验证认证类型
	validAuthTypes := map[string]bool{
		"bearer":  true,
		"api_key": true,
		"basic":   true,
		"none":    true,
	}
	if a.Config.AuthType == "" {
		a.Config.AuthType = "bearer"
	} else if !validAuthTypes[a.Config.AuthType] {
		return fmt.Errorf("无效的认证类型: %s", a.Config.AuthType)
	}

	// 如果认证类型需要API密钥，则验证API密钥
	if (a.Config.AuthType == "bearer" || a.Config.AuthType == "api_key") && a.Config.APIKey == "" {
		return fmt.Errorf("需要API密钥但未配置")
	}

	// 确保映射存在
	if a.Config.Mapping == nil {
		a.Config.Mapping = make(map[string]string)
	}

	// 确保头部和参数映射存在
	if a.Config.Headers == nil {
		a.Config.Headers = make(map[string]string)
	}
	if a.Config.Params == nil {
		a.Config.Params = make(map[string]interface{})
	}

	return nil
}

// ConvertRequest 将自定义渠道请求转换为OpenAI TTS请求
func (a *CustomChannelAdapter) ConvertRequest(request *models.TTSRequest) (*models.OpenAIRequest, error) {
	if request == nil {
		return nil, ErrInvalidRequest
	}

	if request.Text == "" {
		return nil, fmt.Errorf("文本内容不能为空")
	}

	// 创建OpenAI请求
	openaiReq := &models.OpenAIRequest{
		Input: request.Text,
	}

	// 设置默认选项
	openaiReq.Options.Format = a.Config.Format
	openaiReq.Options.SampleRate = a.Config.SampleRate
	openaiReq.Options.Speed = a.Config.Speed

	// 应用请求中的参数（优先级高于配置）
	if request.Model != "" {
		openaiReq.Model = request.Model
	} else if a.Config.Model != "" {
		openaiReq.Model = a.Config.Model
	} else {
		// 默认模型
		openaiReq.Model = "tts-1"
	}

	if request.Voice != "" {
		openaiReq.Voice = request.Voice
	} else if a.Config.Voice != "" {
		openaiReq.Voice = a.Config.Voice
	} else {
		// 默认语音
		openaiReq.Voice = "alloy"
	}

	if request.Format != "" {
		openaiReq.Options.Format = request.Format
	}

	if request.SampleRate > 0 {
		openaiReq.Options.SampleRate = request.SampleRate
	}

	if request.Speed > 0 {
		openaiReq.Options.Speed = request.Speed
	}

	// 应用映射规则（从自定义字段到OpenAI字段）
	if a.Config.Mapping != nil {
		// 这里可以根据映射规则进行更复杂的转换
		// 例如，将自定义参数映射到OpenAI特定参数
		// 对于简单实现，我们已经处理了基本参数
	}

	return openaiReq, nil
}

// GetChannelConfig 获取完整的渠道配置
func (a *CustomChannelAdapter) GetChannelConfig() CustomChannelConfig {
	return a.Config
}

// GetAuthHeaders 获取认证头部
func (a *CustomChannelAdapter) GetAuthHeaders() map[string]string {
	headers := make(map[string]string)

	// 复制配置中的头部
	for k, v := range a.Config.Headers {
		headers[k] = v
	}

	// 根据认证类型添加认证头部
	switch a.Config.AuthType {
	case "bearer":
		headers["Authorization"] = fmt.Sprintf("Bearer %s", a.Config.APIKey)
	case "api_key":
		// API密钥可能在头部或查询参数中
		apiKeyHeader := a.Config.Mapping["api_key_header"]
		if apiKeyHeader != "" {
			headers[apiKeyHeader] = a.Config.APIKey
		}
	}

	return headers
}

// GetQueryParams 获取查询参数
func (a *CustomChannelAdapter) GetQueryParams() map[string]interface{} {
	params := make(map[string]interface{})

	// 复制配置中的参数
	for k, v := range a.Config.Params {
		params[k] = v
	}

	// 如果API密钥在查询参数中
	if a.Config.AuthType == "api_key" {
		apiKeyParam := a.Config.Mapping["api_key_param"]
		if apiKeyParam != "" {
			params[apiKeyParam] = a.Config.APIKey
		}
	}

	return params
}