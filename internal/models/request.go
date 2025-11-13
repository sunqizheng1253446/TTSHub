package models

import "fmt"

// TTSRequest 统一TTS请求模型
type TTSRequest struct {
	ChannelID  uint              `json:"channel_id" binding:"required"`
	Text       string            `json:"text" binding:"required,max=4096"`
	Model      string            `json:"model,omitempty"`
	Voice      string            `json:"voice,omitempty"`
	Speed      float64           `json:"speed,omitempty"`
	Pitch      float64           `json:"pitch,omitempty"`
	Volume     float64           `json:"volume,omitempty"`
	Format     string            `json:"format,omitempty"`
	SampleRate int               `json:"sample_rate,omitempty"`
	Extra      map[string]interface{} `json:"extra,omitempty"`
}

// OpenAIRequest OpenAI TTS请求模型
type OpenAIRequest struct {
	Model string `json:"model"` // 模型名称，如 "tts-1"
	Input string `json:"input"` // 要转换的文本
	Voice string `json:"voice"` // 语音音色
	Options struct {
		Format     string  `json:"format,omitempty"`     // 输出格式，如 "mp3", "opus", "aac", "flac"
		Speed      float64 `json:"speed,omitempty"`      // 语速，0.25-4.0
		SampleRate int     `json:"sample_rate,omitempty"` // 采样率
	} `json:"options,omitempty"`
}

// TTSResponse 统一TTS响应模型
type TTSResponse struct {
	AudioData   []byte `json:"audio_data,omitempty"`
	AudioURL    string `json:"audio_url,omitempty"`
	Format      string `json:"format"`
	Duration    float64 `json:"duration,omitempty"`
	Cost        float64 `json:"cost,omitempty"`
	ChannelID   uint   `json:"channel_id,omitempty"`
	ChannelName string `json:"channel_name,omitempty"`
}

// APIResponse 统一API响应格式
type APIResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// Validate 验证TTS请求参数
func (r *TTSRequest) Validate() error {
	if r.ChannelID == 0 {
		return fmt.Errorf("渠道ID不能为空")
	}
	if r.Text == "" {
		return fmt.Errorf("文本内容不能为空")
	}
	if len(r.Text) > 4096 {
		return fmt.Errorf("文本长度不能超过4096个字符")
	}
	if r.Speed > 0 && (r.Speed < 0.25 || r.Speed > 4.0) {
		return fmt.Errorf("语速必须在0.25-4.0之间")
	}
	return nil
}