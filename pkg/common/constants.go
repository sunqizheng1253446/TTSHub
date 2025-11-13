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

// 常见错误信息
const (
	ErrMsgSuccess           = "操作成功"
	ErrMsgInvalidRequest    = "请求参数无效"
	ErrMsgUnauthorized      = "未授权访问"
	ErrMsgForbidden         = "禁止访问"
	ErrMsgNotFound          = "资源不存在"
	ErrMsgInternalError     = "内部服务器错误"
	ErrMsgChannelNotFound   = "TTS渠道不存在"
	ErrMsgChannelDisabled   = "TTS渠道已禁用"
	ErrMsgInvalidConfig     = "配置无效"
	ErrMsgOpenAIError       = "OpenAI API错误"
)