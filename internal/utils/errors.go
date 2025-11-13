package utils

import (
	"fmt"
	"ttshub/pkg/common"
	"go.uber.org/zap"
)

// AppError 应用错误结构
type AppError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Err     error  `json:"-"`
}

// Error 实现error接口
func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

// Unwrap 实现errors.Unwrap接口
func (e *AppError) Unwrap() error {
	return e.Err
}

// NewAppError 创建应用错误
func NewAppError(code int, message string, err error) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

// NewBadRequestError 创建请求错误
func NewBadRequestError(message string) *AppError {
	return NewAppError(common.ErrCodeInvalidRequest, message, nil)
}

// NewNotFoundError 创建资源不存在错误
func NewNotFoundError(message string) *AppError {
	return NewAppError(common.ErrCodeNotFound, message, nil)
}

// NewInternalError 创建内部错误
func NewInternalError(message string, err error) *AppError {
	return NewAppError(common.ErrCodeInternalError, message, err)
}

// NewChannelNotFoundError 创建渠道不存在错误
func NewChannelNotFoundError(channelID string) *AppError {
	return NewAppError(common.ErrCodeChannelNotFound, fmt.Sprintf("渠道不存在: %s", channelID), nil)
}

// NewChannelDisabledError 创建渠道禁用错误
func NewChannelDisabledError(channelID string) *AppError {
	return NewAppError(common.ErrCodeChannelDisabled, fmt.Sprintf("渠道已禁用: %s", channelID), nil)
}

// NewInvalidConfigError 创建配置无效错误
func NewInvalidConfigError(message string) *AppError {
	return NewAppError(common.ErrCodeInvalidConfig, message, nil)
}

// NewOpenAIError 创建OpenAI错误
func NewOpenAIError(message string, err error) *AppError {
	return NewAppError(common.ErrCodeOpenAIError, message, err)
}

// NewNotImplementedError 创建未实现错误
func NewNotImplementedError(message string) *AppError {
	return NewAppError(common.ErrCodeInternalError, message, nil)
}

// NewExternalError 创建外部服务错误
func NewExternalError(message string) *AppError {
	return NewAppError(common.ErrCodeOpenAIError, message, nil)
}

// StatusCode 返回HTTP状态码
func (e *AppError) StatusCode() int {
	switch e.Code {
	case common.ErrCodeInvalidRequest, common.ErrCodeChannelDisabled, common.ErrCodeInvalidConfig:
		return 400
	case common.ErrCodeUnauthorized:
		return 401
	case common.ErrCodeForbidden:
		return 403
	case common.ErrCodeNotFound, common.ErrCodeChannelNotFound:
		return 404
	case common.ErrCodeInternalError, common.ErrCodeOpenAIError:
		return 500
	default:
		return 500
	}
}

// IsAppError 检查是否为应用错误
func IsAppError(err error) bool {
	_, ok := err.(*AppError)
	return ok
}

// GetAppError 获取应用错误
func GetAppError(err error) *AppError {
	if appErr, ok := err.(*AppError); ok {
		return appErr
	}
	return NewInternalError("未知错误", err)
}

// HandleError 处理错误并记录日志
func HandleError(err error, message string) *AppError {
	appErr := GetAppError(err)
	
	// 根据错误级别记录日志
	switch appErr.Code {
	case common.ErrCodeInvalidRequest, common.ErrCodeNotFound, common.ErrCodeChannelNotFound, common.ErrCodeChannelDisabled, common.ErrCodeInvalidConfig:
		Warn(message, zap.Error(appErr))
	default:
		Error(message, zap.Error(appErr))
	}
	
	return appErr
}