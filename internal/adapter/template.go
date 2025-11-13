package adapter

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
	"ttshub/internal/utils"

	"go.uber.org/zap"
)

// ChannelTemplate 是一个示例模板，用于快速创建新的TTS渠道适配器
// 开发者可以复制这个文件，重命名并根据实际渠道的API规范进行修改

// TemplateConfig 定义了示例渠道的配置结构
type TemplateConfig struct {
	// 基础配置字段
	APIKey      string `json:"api_key" validate:"required"`
	BaseURL     string `json:"base_url" validate:"required,url"`
	Timeout     int    `json:"timeout" validate:"min=1,max=60"`
	MaxRetries  int    `json:"max_retries" validate:"min=0,max=5"`
	EnableProxy bool   `json:"enable_proxy"`
	ProxyURL    string `json:"proxy_url" validate:"omitempty,url"`

	// 渠道特定配置字段 - 根据实际渠道的API需求进行调整
	DefaultVoice  string  `json:"default_voice" validate:"required"`
	DefaultSpeed  float64 `json:"default_speed" validate:"min=0.1,max=3.0"`
	DefaultFormat string  `json:"default_format" validate:"required,oneof=mp3 wav ogg"`
	// 可以添加更多特定的配置字段
}

// TemplateAdapter 实现了TTSAdapter接口，是一个示例适配器
type TemplateAdapter struct {
	config *TemplateConfig
	client *http.Client
	name   string
	typeID string
}

// NewTemplateAdapter 创建一个新的示例适配器实例
// 开发者需要为每个新渠道创建对应的构造函数
func NewTemplateAdapter() TTSAdapter {
	return &TemplateAdapter{
		name:   "Template Channel",
		typeID: "template", // 每个渠道必须有唯一的typeID
	}
}

// Init 初始化适配器
func (a *TemplateAdapter) Init(configJSON string) error {
	// 解析配置JSON
	config := &TemplateConfig{}
	if err := json.Unmarshal([]byte(configJSON), config); err != nil {
		return utils.NewBadRequestError("配置解析失败: " + err.Error())
	}

	// 验证配置
	if err := a.ValidateConfig(configJSON); err != nil {
		return err
	}

	// 设置HTTP客户端
	a.client = createHTTPClient(config.Timeout, config.MaxRetries, config.EnableProxy, config.ProxyURL)
	a.config = config

	utils.Info("模板适配器初始化成功", 
		zap.String("channel_type", a.typeID),
		zap.String("base_url", config.BaseURL),
		zap.Int("timeout", config.Timeout),
		zap.Bool("proxy_enabled", config.EnableProxy),
	)

	return nil
}

// ValidateConfig 验证配置是否有效
func (a *TemplateAdapter) ValidateConfig(configJSON string) error {
	config := &TemplateConfig{}
	if err := json.Unmarshal([]byte(configJSON), config); err != nil {
		return utils.NewBadRequestError("配置解析失败: " + err.Error())
	}

	// 基本验证 - 可以使用更复杂的验证库
	if config.APIKey == "" {
		return utils.NewBadRequestError("API密钥不能为空")
	}

	if config.BaseURL == "" {
		return utils.NewBadRequestError("基础URL不能为空")
	}

	// 验证URL格式
	if !strings.HasPrefix(config.BaseURL, "http://") && !strings.HasPrefix(config.BaseURL, "https://") {
		return utils.NewBadRequestError("基础URL必须以http://或https://开头")
	}

	// 验证超时设置
	if config.Timeout < 1 || config.Timeout > 60 {
		return utils.NewBadRequestError("超时设置必须在1-60秒之间")
	}

	// 验证重试设置
	if config.MaxRetries < 0 || config.MaxRetries > 5 {
		return utils.NewBadRequestError("最大重试次数必须在0-5之间")
	}

	// 如果启用代理，验证代理URL
	if config.EnableProxy {
		if config.ProxyURL == "" {
			return utils.NewBadRequestError("启用代理时必须设置代理URL")
		}
		if !strings.HasPrefix(config.ProxyURL, "http://") && !strings.HasPrefix(config.ProxyURL, "https://") {
			return utils.NewBadRequestError("代理URL必须以http://或https://开头")
		}
	}

	// 验证渠道特定配置
	if config.DefaultVoice == "" {
		return utils.NewBadRequestError("默认语音不能为空")
	}

	if config.DefaultSpeed < 0.1 || config.DefaultSpeed > 3.0 {
		return utils.NewBadRequestError("语音速度必须在0.1-3.0之间")
	}

	validFormats := map[string]bool{"mp3": true, "wav": true, "ogg": true}
	if !validFormats[config.DefaultFormat] {
		return utils.NewBadRequestError("无效的音频格式，支持的格式: mp3, wav, ogg")
	}

	return nil
}

// ConvertRequest 将通用请求转换为OpenAI TTS格式
func (a *TemplateAdapter) ConvertRequest(req *TTSSynthesizeRequest) (*OpenAIRequest, error) {
	if a.config == nil {
		return nil, utils.NewInternalError("适配器未初始化")
	}

	// 使用渠道特定配置或请求参数
	voice := req.Voice
	if voice == "" {
		voice = a.config.DefaultVoice
	}

	speed := req.Speed
	if speed == 0 {
		speed = a.config.DefaultSpeed
	}

	// 创建OpenAI格式的请求
	openAIReq := &OpenAIRequest{
		Model:  "tts-1", // 或根据实际需求选择模型
		Input:  req.Text,
		Voice:  mapVoiceToOpenAI(voice), // 需要实现语音映射函数
		Speed:  speed,
		Format: a.config.DefaultFormat,
	}

	return openAIReq, nil
}

// mapVoiceToOpenAI 将渠道特定的语音ID映射到OpenAI的语音ID
// 这是一个示例实现，开发者需要根据实际渠道的语音列表进行映射
func mapVoiceToOpenAI(voice string) string {
	// 语音映射表
	voiceMap := map[string]string{
		// 示例映射，需要根据实际渠道和OpenAI的语音列表进行调整
		"male1":   "alloy",
		"female1": "echo",
		"male2":   "fable",
		"female2": "onyx",
		"neutral": "nova",
	}

	if mappedVoice, ok := voiceMap[voice]; ok {
		return mappedVoice
	}

	// 默认为alloy
	return "alloy"
}

// Synthesize 直接调用渠道API进行文本转语音（可选实现）
// 注意：如果直接调用渠道API，需要确保返回格式与OpenAI兼容
func (a *TemplateAdapter) Synthesize(req *TTSSynthesizeRequest) ([]byte, string, error) {
	// 这里可以实现直接调用渠道API的逻辑
	// 但根据项目需求，我们主要使用ConvertRequest将请求转换为OpenAI格式
	// 然后由OpenAI客户端处理实际的API调用

	// 示例实现（仅供参考）：
	utils.Warn("模板适配器的Synthesize方法未实现完整功能", zap.String("channel_type", a.typeID))
	return nil, "", utils.NewNotImplementedError("直接合成功能未实现")
}

// GetInfo 获取渠道信息
func (a *TemplateAdapter) GetInfo() *AdapterInfo {
	return &AdapterInfo{
		TypeID:   a.typeID,
		Name:     a.name,
		Version:  "1.0.0",
		Provider: "Template Provider",
		Features: []string{
			"文本转语音",
			"多语言支持",
			"多语音选择",
			"速度调整",
		},
		Status: "active",
	}
}

// GetConfig 获取适配器配置
func (a *TemplateAdapter) GetConfig() interface{} {
	return a.config
}

// TestConnection 测试与渠道API的连接
func (a *TemplateAdapter) TestConnection() error {
	if a.config == nil {
		return utils.NewInternalError("适配器未初始化")
	}

	// 构建测试URL - 根据实际API文档调整
	testURL := a.config.BaseURL + "/v1/test"

	// 创建请求
	req, err := http.NewRequest("GET", testURL, nil)
	if err != nil {
		return utils.NewInternalError("创建请求失败: " + err.Error())
	}

	// 设置认证头 - 根据实际API文档调整
	req.Header.Set("Authorization", "Bearer " + a.config.APIKey)
	req.Header.Set("Content-Type", "application/json")

	// 发送请求
	resp, err := a.client.Do(req)
	if err != nil {
		return utils.NewExternalError("连接失败: " + err.Error())
	}
	defer resp.Body.Close()

	// 检查响应状态码
	if resp.StatusCode != http.StatusOK {
		return utils.NewExternalError(fmt.Sprintf("连接测试失败，状态码: %d", resp.StatusCode))
	}

	utils.Info("模板适配器连接测试成功", zap.String("channel_type", a.typeID), zap.String("base_url", a.config.BaseURL))

	return nil
}

// 创建HTTP客户端的辅助函数
func createHTTPClient(timeout, maxRetries int, enableProxy bool, proxyURL string) *http.Client {
	// 这里可以实现创建HTTP客户端的逻辑，包括超时设置、重试机制、代理配置等
	// 这是一个简化的实现
	client := &http.Client{
		Timeout: time.Duration(timeout) * time.Second,
		// 可以添加更多配置，如代理、重试逻辑等
	}

	return client
}

// 如何使用此模板的说明：
// 1. 复制此文件并重命名为新的渠道名称（如：microsoft.go, google.go等）
// 2. 将TemplateAdapter重命名为对应的渠道适配器（如：MicrosoftAdapter, GoogleAdapter等）
// 3. 将TemplateConfig重命名为对应的渠道配置（如：MicrosoftConfig, GoogleConfig等）
// 4. 修改NewTemplateAdapter函数为对应的渠道构造函数
// 5. 根据渠道的API规范修改Init, ValidateConfig, ConvertRequest等方法
// 6. 实现渠道特定的语音映射函数
// 7. 在registry.go中注册新的适配器
// 8. 编写相应的单元测试

// 适配器开发最佳实践：
// 1. 提供详细的配置验证，确保用户输入的配置有效
// 2. 实现完整的错误处理，返回清晰的错误信息
// 3. 记录关键操作的日志，便于调试和监控
// 4. 实现连接测试功能，帮助用户验证配置
// 5. 提供详细的文档注释
// 6. 遵循Go语言的代码规范