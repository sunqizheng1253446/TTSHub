# TTS Hub

TTS Hub 是一个统一的文本转语音(TTS)服务平台，支持多种TTS供应商的集成，提供统一的API接口和管理控制台，实现多渠道语音合成服务的集中管理和调用。

## 功能特性

- **统一API接口**：提供标准化的TTS合成接口，屏蔽不同供应商的API差异
- **多渠道支持**：内置支持OpenAI TTS，可通过适配器框架轻松扩展其他供应商
- **渠道管理**：Web控制台可视化管理TTS渠道配置、密钥等信息
- **语音合成测试**：在线测试不同渠道的合成效果，支持文本输入和音频播放
- **健康监控**：渠道可用性监控和性能统计
- **可扩展架构**：基于适配器模式设计，便于添加新的TTS供应商
- **本地存储**：使用SQLite存储渠道配置和调用历史
- **RESTful API**：完整的RESTful接口，便于集成到各种系统

## 技术栈

- **后端**：Go语言 (1.23+)
- **Web框架**：Gin
- **存储**：SQLite
- **前端**：HTML, CSS, JavaScript
- **配置管理**：Viper
- **日志**：Zap

## 安装指南

### 前提条件

- Go 1.23 或更高版本
- Git

### 克隆项目

```bash
git clone https://your-repository/ttshub.git
cd ttshub
```

### 安装依赖

```bash
go mod tidy
```

### 构建项目

```bash
go build -o ttshub ./cmd/ttshub
```

### 运行服务

```bash
# 直接运行
./ttshub

# 或使用Go命令运行
go run ./cmd/ttshub
```

服务默认在 `http://localhost:8080` 启动。

## 配置说明

TTS Hub 支持通过环境变量或配置文件进行配置。默认配置文件为 `config.yaml`。

### 环境变量

| 环境变量 | 描述 | 默认值 |
|---------|------|--------|
| `PORT` | 服务端口 | `8080` |
| `DB_PATH` | SQLite数据库文件路径 | `ttshub.db` |
| `LOG_LEVEL` | 日志级别 (debug, info, warn, error) | `info` |
| `CORS_ALLOWED_ORIGINS` | 允许的CORS来源 | `*` |

### 配置文件示例 (config.yaml)

```yaml
server:
  port: 8080
  host: "0.0.0.0"

storage:
  db_path: "ttshub.db"

logging:
  level: "info"
  format: "json"

cors:
  allowed_origins: ["*"]
  allowed_methods: ["GET", "POST", "PUT", "DELETE", "OPTIONS"]
  allowed_headers: ["*"]
  allow_credentials: true
```

## Web控制台

服务启动后，可以通过浏览器访问 `http://localhost:8080` 打开Web控制台，控制台提供以下功能：

- **仪表盘**：显示渠道统计和使用情况
- **渠道管理**：添加、编辑、删除TTS渠道配置
- **TTS测试**：在线测试文本转语音功能
- **日志查询**：查看系统和API调用日志
- **系统设置**：配置系统参数

## API文档

### 健康检查

```
GET /health
```

**响应示例**：
```json
{"status": "healthy"}
```

### 获取适配器类型

```
GET /api/v1/adapters
```

**响应示例**：
```json
["openai", "baidu", "google", "azure"]
```

### 渠道管理

#### 列出所有渠道
```
GET /api/v1/channels
```

**响应示例**：
```json
[
  {
    "id": "openai-default",
    "name": "OpenAI默认渠道",
    "adapter_type": "openai",
    "config": {
      "api_key": "sk-...",
      "model": "tts-1",
      "voice": "alloy"
    },
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  }
]
```

#### 获取单个渠道
```
GET /api/v1/channels/{id}
```

#### 创建渠道
```
POST /api/v1/channels
```

**请求体**：
```json
{
  "id": "openai-default",
  "name": "OpenAI默认渠道",
  "adapter_type": "openai",
  "config": {
    "api_key": "sk-...",
    "model": "tts-1",
    "voice": "alloy"
  }
}
```

#### 更新渠道
```
PUT /api/v1/channels/{id}
```

**请求体**：与创建渠道相同

#### 删除渠道
```
DELETE /api/v1/channels/{id}
```

### 文本转语音合成

```
POST /api/v1/synthesize
```

**请求体**：
```json
{
  "channel": "openai-default",
  "text": "你好，这是一段测试文本。",
  "params": {
    "voice": "nova",
    "speed": 1.0
  }
}
```

**响应**：返回二进制音频数据，Content-Type为对应的音频格式。

## 项目结构

```
ttshub/
├── cmd/
│   └── ttshub/              # 主程序入口
│       └── main.go
├── internal/
│   ├── adapter/             # TTS适配器框架
│   │   ├── adapter.go       # 适配器接口定义
│   │   ├── openai.go        # OpenAI适配器实现
│   │   └── test/            # 适配器测试
│   ├── config/              # 配置管理
│   │   └── config.go
│   ├── handlers/            # HTTP处理器
│   │   ├── handlers.go
│   │   ├── tts_handler.go
│   │   └── test/            # 处理器测试
│   ├── service/             # 业务逻辑层
│   │   ├── service.go       # TTS服务接口
│   │   ├── tts_service.go   # 服务实现
│   │   └── test/            # 服务测试
│   ├── storage/             # 存储层
│   │   ├── storage.go       # 存储接口
│   │   ├── sqlite_store.go  # SQLite实现
│   │   └── test/            # 存储测试
│   └── utils/               # 工具函数
│       ├── logger.go
│       └── response.go
├── public/                  # 前端静态文件
│   ├── index.html
│   ├── css/
│   ├── js/
│   └── assets/
├── docs/                    # 文档
├── go.mod
├── go.sum
├── config.yaml.example
└── README.md
```

## 添加新的TTS渠道

要添加新的TTS供应商支持，需要实现 `adapter.TTSAdapter` 接口：

```go
// 在internal/adapter目录下创建新的适配器文件，例如custom.go
package adapter

import (
	"context"
	"ttshub/internal/service"
)

// CustomAdapter 实现TTSAdapter接口
type CustomAdapter struct {
	config map[string]interface{}
}

// NewCustomAdapter 创建新的Custom适配器
func NewCustomAdapter(config map[string]interface{}) (TTSAdapter, error) {
	return &CustomAdapter{config: config}, nil
}

// Synthesize 实现文本转语音功能
func (a *CustomAdapter) Synthesize(ctx context.Context, req *service.SynthesizeRequest) (*service.SynthesizeResponse, error) {
	// 实现与自定义TTS API的集成逻辑
	// ...
	return &service.SynthesizeResponse{
		AudioData: audioBytes,
		Format:    "mp3",
		Duration:  durationMs,
	}, nil
}

// ValidateConfig 验证适配器配置
func (a *CustomAdapter) ValidateConfig() error {
	// 验证配置是否有效
	// ...
	return nil
}
```

然后在 `cmd/ttshub/main.go` 中注册新的适配器：

```go
// 注册新的适配器	service.RegisterAdapter("custom", adapter.NewCustomAdapter)
```

## 开发指南

### 运行测试

```bash
go test ./...
```

### 代码风格

- 遵循Go官方代码风格
- 使用 `go fmt` 格式化代码
- 使用 `go vet` 检查潜在问题

### 提交规范

- 提交消息清晰描述变更内容
- 功能开发完成后先运行测试确保无问题
- 大型变更建议先提交Pull Request进行代码审查

## 故障排除

### 常见问题

1. **服务无法启动**：检查端口是否被占用，数据库文件权限是否正确
2. **API调用失败**：检查渠道配置是否正确，API密钥是否有效
3. **音频合成失败**：检查文本格式，确认渠道支持的语言和字符集

### 日志查看

系统日志默认输出到标准输出，可以通过调整配置文件中的日志级别来获取更详细的信息。

## 许可证

[MIT License](LICENSE)

## 联系方式

如有问题或建议，请通过以下方式联系：

- 项目仓库：https://your-repository/ttshub
- 邮件：your-email@example.com