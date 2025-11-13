package adapter

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
	"ttshub/internal/config"
	"ttshub/internal/models"
	"ttshub/internal/utils"

	"go.uber.org/zap"
)

// OpenAIClient OpenAI TTS客户端
type OpenAIClient struct {
	HTTPClient *http.Client
	BaseURL    string
	APIKey     string
	Timeout    time.Duration
}

// NewOpenAIClient 创建OpenAI客户端实例
func NewOpenAIClient() *OpenAIClient {
	cfg := config.GetConfig()
	return &OpenAIClient{
		HTTPClient: &http.Client{},
		BaseURL:    cfg.OpenAI.BaseURL,
		APIKey:     cfg.OpenAI.APIKey,
		Timeout:    time.Duration(cfg.OpenAI.Timeout) * time.Second,
	}
}

// WithConfig 使用指定配置创建客户端
func (c *OpenAIClient) WithConfig(apiKey string, baseURL string, timeout time.Duration) *OpenAIClient {
	c.APIKey = apiKey
	c.BaseURL = baseURL
	c.Timeout = timeout
	return c
}

// Synthesize 执行文本转语音
func (c *OpenAIClient) Synthesize(ctx context.Context, request *models.OpenAIRequest) (*models.TTSResponse, error) {
	if request == nil {
		return nil, ErrInvalidRequest
	}

	if c.APIKey == "" {
		return nil, fmt.Errorf("OpenAI API密钥未配置")
	}

	// 序列化请求
	requestBody, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("序列化请求失败: %w", err)
	}

	// 创建HTTP请求
	req, err := http.NewRequestWithContext(ctx, "POST", fmt.Sprintf("%s/v1/audio/speech", c.BaseURL), bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.APIKey))

	// 设置超时
	timeoutCtx, cancel := context.WithTimeout(ctx, c.Timeout)
	defer cancel()
	req = req.WithContext(timeoutCtx)

	// 发送请求
	startTime := time.Now()
	utils.Debug("发送OpenAI TTS请求", zap.String("model", request.Model), zap.String("voice", request.Voice))
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		utils.Error("OpenAI TTS请求失败", zap.Error(err))
		return nil, utils.NewOpenAIError("请求失败", err)
	}
	defer resp.Body.Close()

	// 记录请求耗时
	duration := time.Since(startTime)
	utils.Info("OpenAI TTS请求完成", zap.Int("status_code", resp.StatusCode), zap.Duration("duration", duration))

	// 检查响应状态
	if resp.StatusCode != http.StatusOK {
		errBody, _ := io.ReadAll(resp.Body)
		errMsg := fmt.Sprintf("请求失败，状态码: %d, 响应: %s", resp.StatusCode, string(errBody))
		utils.Error(errMsg)
		return nil, utils.NewOpenAIError(errMsg, nil)
	}

	// 读取音频数据
	audioData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取音频数据失败: %w", err)
	}

	// 估算成本（示例计算，实际成本以OpenAI官方为准）
	cost := c.calculateCost(request.Input)

	// 估算音频时长（基于平均语速）
	// 假设平均每秒15个汉字或5个英文单词
	durationEstimate := c.estimateAudioDuration(request.Input)

	// 返回响应
	return &models.TTSResponse{
		AudioData: audioData,
		Format:    request.Options.Format,
		Duration:  durationEstimate,
		Cost:      cost,
	}, nil
}

// calculateCost 估算API调用成本
// 注意：这是示例计算，实际成本请参考OpenAI官方文档
func (c *OpenAIClient) calculateCost(text string) float64 {
	// tts-1模型: $0.015/1K characters
	// tts-1-hd模型: $0.03/1K characters
	charCount := float64(len([]rune(text)))
	costPer1KChars := 0.015 // 默认tts-1
	return (charCount / 1000.0) * costPer1KChars
}

// estimateAudioDuration 估算音频时长
func (c *OpenAIClient) estimateAudioDuration(text string) float64 {
	// 简单估算：平均每秒15个汉字或5个英文单词
	runes := []rune(text)
	charCount := len(runes)
	
	// 假设平均语速为每秒15个字符
	return float64(charCount) / 15.0
}

// ValidateAPIKey 验证API密钥是否有效
func (c *OpenAIClient) ValidateAPIKey() error {
	// 创建一个简单的验证请求
	request := &models.OpenAIRequest{
		Model: "tts-1",
		Input: "test",
		Voice: "alloy",
	}

	// 发送验证请求
	_, err := c.Synthesize(context.Background(), request)
	return err
}