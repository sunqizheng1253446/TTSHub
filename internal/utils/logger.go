package utils

import (
	"os"
	"path/filepath"
	"ttshub/internal/config"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger 全局日志实例
var Logger *zap.Logger

// InitLogger 初始化日志系统
func InitLogger(cfg *config.LogConfig) error {
	// 创建日志目录
	if err := os.MkdirAll(cfg.Path, 0755); err != nil {
		return err
	}

	// 配置编码器
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "timestamp"
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	var encoder zapcore.Encoder
	if cfg.Format == "console" {
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	} else {
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	}

	// 设置日志级别
	var level zapcore.Level
	switch cfg.Level {
	case "debug":
		level = zapcore.DebugLevel
	case "info":
		level = zapcore.InfoLevel
	case "warn":
		level = zapcore.WarnLevel
	case "error":
		level = zapcore.ErrorLevel
	default:
		level = zapcore.InfoLevel
	}

	// 创建文件输出
	logFilePath := filepath.Join(cfg.Path, "ttshub.log")
	fileWriter, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	// 配置Core
	core := zapcore.NewCore(
		encoder,
		zapcore.NewMultiWriteSyncer(zapcore.AddSync(fileWriter), zapcore.AddSync(os.Stdout)),
		level,
	)

	// 创建Logger
	Logger = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))

	// 替换全局Logger
	zap.ReplaceGlobals(Logger)

	Logger.Info("日志系统初始化完成", zap.String("log_path", logFilePath), zap.String("level", cfg.Level))
	return nil
}

// GetLogger 获取日志实例
func GetLogger() *zap.Logger {
	if Logger == nil {
		// 如果日志未初始化，使用默认配置
		defaultLogConfig := &config.LogConfig{
			Level:  "info",
			Path:   "./logs",
			Format: "json",
		}
		if err := InitLogger(defaultLogConfig); err != nil {
			// 如果初始化失败，使用默认的Zap Logger
			Logger, _ = zap.NewProduction()
		}
	}
	return Logger
}

// Debug 记录调试日志
func Debug(msg string, fields ...zapcore.Field) {
	GetLogger().Debug(msg, fields...)
}

// Info 记录信息日志
func Info(msg string, fields ...zapcore.Field) {
	GetLogger().Info(msg, fields...)
}

// Warn 记录警告日志
func Warn(msg string, fields ...zapcore.Field) {
	GetLogger().Warn(msg, fields...)
}

// Error 记录错误日志
func Error(msg string, fields ...zapcore.Field) {
	GetLogger().Error(msg, fields...)
}

// Fatal 记录致命错误日志
func Fatal(msg string, fields ...zapcore.Field) {
	GetLogger().Fatal(msg, fields...)
}