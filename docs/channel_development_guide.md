# TTSHub 渠道开发指南

本文档提供了在 TTSHub 中添加新的 TTS 渠道的详细指南。通过遵循这些步骤，您可以快速开发并集成新的 TTS 服务提供商。

## 目录

1. [概述](#概述)
2. [渠道适配器架构](#渠道适配器架构)
3. [开发步骤](#开发步骤)
4. [配置管理](#配置管理)
5. [请求转换](#请求转换)
6. [错误处理](#错误处理)
7. [测试](#测试)
8. [注册和部署](#注册和部署)
9. [示例：添加百度TTS渠道](#示例添加百度TTS渠道)

## 概述

TTSHub 采用适配器模式设计，通过统一的接口抽象不同的 TTS 服务提供商。每个渠道适配器需要实现 `TTSAdapter` 接口，负责：

- 初始化渠道配置
- 验证配置有效性
- 将通用请求转换为 OpenAI TTS 格式
- 可选：直接调用渠道 API（如需）
- 提供渠道信息
- 测试连接

## 渠道适配器架构

核心接口定义在 `adapter/adapter.go` 中：

```go
type TTSAdapter interface {
	Init(configJSON string) error
	ValidateConfig(configJSON string) error
	ConvertRequest(req *TTSSynthesizeRequest) (*OpenAIRequest, error)
	Synthesize(req *TTSSynthesizeRequest) ([]byte, string, error)
	GetInfo() *AdapterInfo
	GetConfig() interface{}
	TestConnection() error
}
```

## 开发步骤

### 1. 创建适配器文件

在 `internal/adapter` 目录下创建新的适配器文件，例如 `baidu.go`。可以基于 `template.go` 模板进行修改。

### 2. 定义配置结构体

根据渠道 API 的要求，定义配置结构体：

```go
type BaiduConfig struct {
	APIKey      string `json:"api_key"`
	SecretKey   string `json:"secret_key"`
	BaseURL     string `json:"base_url"`
	DefaultVoice string `json:"default_voice"`
	// 其他必要配置字段
}
```

### 3. 实现适配器结构体

创建适配器结构体并实现 `TTSAdapter` 接口：

```go
type BaiduAdapter struct {
	config *BaiduConfig
	client *http.Client
	name   string
	typeID string
}

func NewBaiduAdapter() TTSAdapter {
	return &BaiduAdapter{
		name:   "百度语音合成",
		typeID: "baidu",
	}
}
```

### 4. 实现接口方法

#### 4.1 初始化方法 (Init)

```go
func (a *BaiduAdapter) Init(configJSON string) error {
	// 解析配置
	config := &BaiduConfig{}
	if err := json.Unmarshal([]byte(configJSON), config); err != nil {
		return utils.NewBadRequestError("配置解析失败: " + err.Error())
	}

	// 验证配置
	if err := a.ValidateConfig(configJSON); err != nil {
		return err
	}

	// 设置HTTP客户端和配置
	a.client = createHTTPClient(config.Timeout, config.MaxRetries, config.EnableProxy, config.ProxyURL)
	a.config = config

	utils.Info("百度适配器初始化成功", zap.String("channel_type", a.typeID))

	return nil
}
```

#### 4.2 配置验证方法 (ValidateConfig)

```go
func (a *BaiduAdapter) ValidateConfig(configJSON string) error {
	config := &BaiduConfig{}
	if err := json.Unmarshal([]byte(configJSON), config); err != nil {
		return utils.NewBadRequestError("配置解析失败: " + err.Error())
	}

	// 验证必填字段
	if config.APIKey == "" {
		return utils.NewBadRequestError("API密钥不能为空")
	}

	if config.SecretKey == "" {
		return utils.NewBadRequestError("Secret Key不能为空")
	}

	// 其他验证逻辑

	return nil
}
```

#### 4.3 请求转换方法 (ConvertRequest)

```go
func (a *BaiduAdapter) ConvertRequest(req *TTSSynthesizeRequest) (*OpenAIRequest, error) {
	if a.config == nil {
		return nil, utils.NewInternalError("适配器未初始化")
	}

	// 映射语音和参数
	voice := req.Voice
	if voice == "" {
		voice = a.config.DefaultVoice
	}

	// 创建OpenAI格式请求
	openAIReq := &OpenAIRequest{
		Model:  "tts-1",
		Input:  req.Text,
		Voice:  mapBaiduVoiceToOpenAI(voice),
		Speed:  req.Speed,
		Format: "mp3",
	}

	return openAIReq, nil
}
```

#### 4.4 其他必要方法

实现 `Synthesize`, `GetInfo`, `GetConfig` 和 `TestConnection` 方法。

### 5. 注册适配器

在 `adapter/registry.go` 中注册新的适配器：

```go
func RegisterBuiltinAdapters() {
	// 已有适配器
	Register("openai", NewOpenAIAdapter())
	Register("custom", NewCustomChannelAdapter())
	// 注册新适配器
	Register("baidu", NewBaiduAdapter())
}
```

## 配置管理

每个渠道适配器都有自己的配置结构。配置以 JSON 格式存储在 SQLite 数据库中。确保：

- 提供默认配置值
- 验证所有必要字段
- 支持常见配置如超时、代理等

## 请求转换

关键功能是将渠道特定的请求参数映射到 OpenAI TTS 格式。需要处理：

- 语音映射（渠道特定的语音ID到OpenAI语音ID）
- 速度调整
- 格式转换
- 文本预处理（如编码、长度限制等）

## 错误处理

使用 `utils` 包中的错误创建函数返回标准化错误：

- `NewBadRequestError`: 参数错误
- `NewUnauthorizedError`: 认证失败
- `NewExternalError`: 外部API错误
- `NewInternalError`: 内部错误
- `NewNotImplementedError`: 未实现功能

## 测试

为新适配器编写测试：

1. 单元测试：测试各个方法的功能
2. 集成测试：测试端到端的请求流程
3. 手动测试：使用 API 或控制台测试实际请求

## 注册和部署

1. 注册适配器到注册表
2. 重新构建应用
3. 启动服务
4. 通过 API 创建新渠道配置

## 示例：添加百度TTS渠道

以下是添加百度TTS渠道的完整示例：

### 1. 创建 `baidu.go`

```go
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

type BaiduConfig struct {
	APIKey      string `json:"api_key"`
	SecretKey   string `json:"secret_key"`
	BaseURL     string `json:"base_url"`
	Timeout     int    `json:"timeout"`
	DefaultVoice string `json:"default_voice"`
}

type BaiduAdapter struct {
	config *BaiduConfig
	client *http.Client
	name   string
	typeID string
}

func NewBaiduAdapter() TTSAdapter {
	return &BaiduAdapter{
		name:   "百度语音合成",
		typeID: "baidu",
	}
}

func (a *BaiduAdapter) Init(configJSON string) error {
	config := &BaiduConfig{}
	if err := json.Unmarshal([]byte(configJSON), config); err != nil {
		return utils.NewBadRequestError("配置解析失败: " + err.Error())
	}

	if err := a.ValidateConfig(configJSON); err != nil {
		return err
	}

	a.client = &http.Client{
		Timeout: time.Duration(config.Timeout) * time.Second,
	}
	a.config = config

	utils.Info("百度适配器初始化成功", zap.String("channel_type", a.typeID))

	return nil
}

func (a *BaiduAdapter) ValidateConfig(configJSON string) error {
	config := &BaiduConfig{}
	if err := json.Unmarshal([]byte(configJSON), config); err != nil {
		return utils.NewBadRequestError("配置解析失败: " + err.Error())
	}

	if config.APIKey == "" {
		return utils.NewBadRequestError("API密钥不能为空")
	}

	if config.SecretKey == "" {
		return utils.NewBadRequestError("Secret Key不能为空")
	}

	if config.BaseURL == "" {
		config.BaseURL = "https://tsn.baidu.com/text2audio"
	}

	return nil
}

func (a *BaiduAdapter) ConvertRequest(req *TTSSynthesizeRequest) (*OpenAIRequest, error) {
	if a.config == nil {
		return nil, utils.NewInternalError("适配器未初始化")
	}

	voice := req.Voice
	if voice == "" {
		voice = a.config.DefaultVoice
	}

	speed := req.Speed
	if speed == 0 {
		speed = 1.0
	}

	openAIReq := &OpenAIRequest{
		Model:  "tts-1",
		Input:  req.Text,
		Voice:  mapBaiduVoiceToOpenAI(voice),
		Speed:  speed,
		Format: "mp3",
	}

	return openAIReq, nil
}

func mapBaiduVoiceToOpenAI(voice string) string {
	voiceMap := map[string]string{
		"度小宇": "alloy",
		"度小美": "echo",
		"度逍遥": "fable",
		"度丫丫": "onyx",
	}

	if mappedVoice, ok := voiceMap[voice]; ok {
		return mappedVoice
	}

	return "alloy"
}

func (a *BaiduAdapter) Synthesize(req *TTSSynthesizeRequest) ([]byte, string, error) {
	utils.Warn("百度适配器的Synthesize方法未实现")
	return nil, "", utils.NewNotImplementedError("直接合成功能未实现")
}

func (a *BaiduAdapter) GetInfo() *AdapterInfo {
	return &AdapterInfo{
		TypeID:   a.typeID,
		Name:     a.name,
		Version:  "1.0.0",
		Provider: "百度AI",
		Features: []string{
			"中文语音合成",
			"多音色支持",
		},
		Status: "active",
	}
}

func (a *BaiduAdapter) GetConfig() interface{} {
	return a.config
}

func (a *BaiduAdapter) TestConnection() error {
	if a.config == nil {
		return utils.NewInternalError("适配器未初始化")
	}

	// 实现连接测试逻辑
	// ...

	return nil
}
```

### 2. 注册适配器

在 `registry.go` 中添加：

```go
func RegisterBuiltinAdapters() {
	Register("openai", NewOpenAIAdapter())
	Register("custom", NewCustomChannelAdapter())
	Register("baidu", NewBaiduAdapter())
}
```

## 常见问题解答

### Q: 如何处理渠道特定的参数？
A: 在配置结构体中添加渠道特定字段，并在 `ConvertRequest` 方法中处理这些参数。

### Q: 如何处理语音映射？
A: 创建一个映射函数，将渠道特定的语音ID映射到OpenAI支持的语音ID。

### Q: 如何实现重试机制？
A: 可以使用第三方库如 `github.com/avast/retry-go` 或在HTTP客户端中实现自定义重试逻辑。

### Q: 如何支持代理？
A: 在配置中添加代理相关字段，并在创建HTTP客户端时配置代理设置。

## 最佳实践

1. 提供详细的配置文档和示例
2. 实现全面的错误处理和日志记录
3. 确保线程安全
4. 性能优化（如连接池、缓存等）
5. 遵循Go语言代码规范

---

通过遵循本指南，您可以轻松地为 TTSHub 添加新的 TTS 渠道支持。如有任何问题，请参考现有适配器的实现或创建issue。