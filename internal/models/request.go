package models

// TTSRequest 统一TTS请求模型
type TTSRequest struct {
	ChannelID  string            `json:"channel_id" binding:"required"`
	Text       string            `json:"text" binding:"required,max=4096"`
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
	AudioData []byte `json:"audio_data,omitempty"`
	AudioURL  string `json:"audio_url,omitempty"`
	Format    string `json:"format"`
	Duration  float64 `json:"duration,omitempty"`
	Cost      float64 `json:"cost,omitempty"`
}

// APIResponse 统一API响应格式
type APIResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}