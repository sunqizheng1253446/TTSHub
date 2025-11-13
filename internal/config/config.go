package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"ttshub/pkg/common"

	"github.com/spf13/viper"
)

// Config 应用配置结构
type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	Log      LogConfig      `mapstructure:"log"`
	OpenAI   OpenAIConfig   `mapstructure:"openai"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Port string `mapstructure:"port"`
	Host string `mapstructure:"host"`
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Path string `mapstructure:"path"`
}

// LogConfig 日志配置
type LogConfig struct {
	Level  string `mapstructure:"level"`
	Path   string `mapstructure:"path"`
	Format string `mapstructure:"format"`
}

// OpenAIConfig OpenAI配置
type OpenAIConfig struct {
	APIKey    string `mapstructure:"api_key"`
	BaseURL   string `mapstructure:"base_url"`
	Model     string `mapstructure:"model"`
	Timeout   int    `mapstructure:"timeout"`
	MaxTokens int    `mapstructure:"max_tokens"`
}

// 全局配置实例
var AppConfig *Config

// LoadConfig 加载配置
func LoadConfig(configPath string) (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	// 设置默认值
	setDefaultConfig()

	// 添加配置路径
	viper.AddConfigPath(configPath)
	viper.AddConfigPath("./configs")
	viper.AddConfigPath(".")

	// 读取环境变量
	viper.AutomaticEnv()

	// 读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Printf("配置文件未找到，使用默认配置和环境变量")
			// 创建默认配置文件
			if err := createDefaultConfigFile(); err != nil {
				log.Printf("创建默认配置文件失败: %v", err)
			}
		} else {
			return nil, fmt.Errorf("读取配置文件错误: %v", err)
		}
	}

	// 解析配置
	config := &Config{}
	if err := viper.Unmarshal(config); err != nil {
		return nil, fmt.Errorf("解析配置错误: %v", err)
	}

	// 确保数据目录存在
	ensureDataDir(config.Database.Path)

	AppConfig = config
	return config, nil
}

// 设置默认配置
func setDefaultConfig() {
	viper.SetDefault("server.port", common.DefaultServerPort)
	viper.SetDefault("server.host", "0.0.0.0")
	viper.SetDefault("database.path", common.DefaultDatabasePath)
	viper.SetDefault("log.level", common.DefaultLogLevel)
	viper.SetDefault("log.path", "./logs")
	viper.SetDefault("log.format", "json")
	viper.SetDefault("openai.base_url", "https://api.openai.com")
	viper.SetDefault("openai.model", "tts-1")
	viper.SetDefault("openai.timeout", 30)
	viper.SetDefault("openai.max_tokens", 4096)
}

// 创建默认配置文件
func createDefaultConfigFile() error {
	// 确保configs目录存在
	if err := os.MkdirAll("./configs", 0755); err != nil {
		return err
	}

	// 默认配置内容
	defaultConfig := `# TTSHub 配置文件
server:
  port: "8080"
  host: "0.0.0.0"

database:
  path: "./data/ttshub.db"

log:
  level: "info"
  path: "./logs"
  format: "json"

openai:
  api_key: "your_api_key_here"
  base_url: "https://api.openai.com"
  model: "tts-1"
  timeout: 30
  max_tokens: 4096
`

	// 写入配置文件
	return os.WriteFile("./configs/config.yaml", []byte(defaultConfig), 0644)
}

// 确保数据目录存在
func ensureDataDir(dbPath string) {
	dir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		log.Printf("创建数据目录失败: %v", err)
	}
}

// GetConfig 获取应用配置
func GetConfig() *Config {
	if AppConfig == nil {
		// 如果配置未加载，尝试加载默认配置
		config, err := LoadConfig(".")
		if err != nil {
			log.Printf("加载配置失败: %v", err)
			return &Config{}
		}
		return config
	}
	return AppConfig
}