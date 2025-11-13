package adapter

import (
	"encoding/json"
	"fmt"
	"ttshub/internal/models"
)

// OpenAIAdapterConfig OpenAI适配器配置
type OpenAIAdapterConfig struct {
	APIKey   string `json:"api_key"`
	Model    string `json:"model"`
	Voice    string `json:"voice"`
	Format   string `json:"format"`
	Speed    float64 `json:"speed"`
	SampleRate int    `json:"sample_rate"`
}

// OpenAIAdapter OpenAI适配器实现
type OpenAIAdapter struct {
	BaseAdapter
	Config OpenAIAdapterConfig
}

// NewOpenAIAdapter 创建OpenAI适配器实例
func NewOpenAIAdapter() *OpenAIAdapter {
	return &OpenAIAdapter{
		BaseAdapter: BaseAdapter{
			ChannelName: "openai",
		},
		Config: OpenAIAdapterConfig{
			Model:      "tts-1",
			Voice:      "alloy",
			Format:     "mp3",
			Speed:      1.0,
			SampleRate: 24000,
		},
	}
}

// Init 初始化适配器
func (a *OpenAIAdapter) Init(config string) error {
	a.Config = OpenAIAdapterConfig{
		Model:      "tts-1",
		Voice:      "alloy",
		Format:     "mp3",
		Speed:      1.0,
		SampleRate: 24000,
	}

	if config != "" {
		if err := json.Unmarshal([]byte(config), &a.Config); err != nil {
			return fmt.Errorf("解析配置失败: %w", err)
		}
	}

	return nil
}

// ValidateConfig 验证配置
func (a *OpenAIAdapter) ValidateConfig() error {
	// 基础验证
	if err := a.BaseAdapter.ValidateConfig(); err != nil {
		return err
	}

	// OpenAI特定验证
	if a.Config.Model == "" {
		a.Config.Model = "tts-1"
	}

	if a.Config.Voice == "" {
		a.Config.Voice = "alloy"
	}

	// 验证语音选项
	validVoices := map[string]bool{
		"alloy":   true,
		"echo":    true,
		"fable":   true,
		"onyx":    true,
		"nova":    true,
		"shimmer": true,
	}
	if !validVoices[a.Config.Voice] {
		return fmt.Errorf("无效的语音选项: %s", a.Config.Voice)
	}

	// 验证格式选项
	validFormats := map[string]bool{
		"mp3":  true,
		"opus": true,
		"aac":  true,
		"flac": true,
	}
	if a.Config.Format == "" {
		a.Config.Format = "mp3"
	} else if !validFormats[a.Config.Format] {
		return fmt.Errorf("无效的音频格式: %s", a.Config.Format)
	}

	// 验证语速范围
	if a.Config.Speed < 0.25 || a.Config.Speed > 4.0 {
		return fmt.Errorf("语速必须在0.25-4.0之间: %f", a.Config.Speed)
	}

	return nil
}

// ConvertRequest 将通用TTS请求转换为OpenAI TTS请求
func (a *OpenAIAdapter) ConvertRequest(request *models.TTSRequest) (*models.OpenAIRequest, error) {
	if request == nil {
		return nil, ErrInvalidRequest
	}

	if request.Text == "" {
		return nil, fmt.Errorf("文本内容不能为空")
	}

	// 创建OpenAI请求
	openaiReq := &models.OpenAIRequest{
		Model: a.Config.Model,
		Input: request.Text,
		Voice: a.Config.Voice,
	}

	// 设置选项
	openaiReq.Options.Format = a.Config.Format
	openaiReq.Options.SampleRate = a.Config.SampleRate
	openaiReq.Options.Speed = a.Config.Speed

	// 应用请求中指定的参数（优先级高于配置）
	if request.Voice != "" {
		// 验证语音选项
		validVoices := map[string]bool{
			"alloy":   true,
			"echo":    true,
			"fable":   true,
			"onyx":    true,
			"nova":    true,
			"shimmer": true,
		}
		if !validVoices[request.Voice] {
			return nil, fmt.Errorf("无效的语音选项: %s", request.Voice)
		}
		openaiReq.Voice = request.Voice
	}

	if request.Format != "" {
		// 验证格式选项
		validFormats := map[string]bool{
			"mp3":  true,
			"opus": true,
			"aac":  true,
			"flac": true,
		}
		if !validFormats[request.Format] {
			return nil, fmt.Errorf("无效的音频格式: %s", request.Format)
		}
		openaiReq.Options.Format = request.Format
	}

	if request.SampleRate > 0 {
		openaiReq.Options.SampleRate = request.SampleRate
	}

	if request.Speed > 0 {
		if request.Speed < 0.25 || request.Speed > 4.0 {
			return nil, fmt.Errorf("语速必须在0.25-4.0之间: %f", request.Speed)
		}
		openaiReq.Options.Speed = request.Speed
	}

	return openaiReq, nil
}