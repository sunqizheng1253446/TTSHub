package common

// 渠道类型常量
const (
	ChannelTypeOpenAI = "openai"
	ChannelTypeCustom = "custom"
)

// 渠道状态常量
const (
	ChannelStatusEnabled  = 1
	ChannelStatusDisabled = 0
)

// 配置默认值
const (
	DefaultDatabasePath = "./data/ttshub.db"
	DefaultServerPort   = "8080"
	DefaultLogLevel     = "info"
)

// 错误码定义
const (
	ErrCodeSuccess           = 0
	ErrCodeInvalidRequest    = 400
	ErrCodeUnauthorized      = 401
	ErrCodeForbidden         = 403
	ErrCodeNotFound          = 404
	ErrCodeInternalError     = 500
	ErrCodeChannelNotFound   = 1001
	ErrCodeChannelDisabled   = 1002
	ErrCodeInvalidConfig     = 1003
	ErrCodeOpenAIError       = 2001
)

// API版本
const (
	APIVersionV1 = "v1"
)