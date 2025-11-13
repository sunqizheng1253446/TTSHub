package service

import (
	"context"
	"encoding/json"
	"fmt"
	"ttshub/internal/adapter"
	"ttshub/internal/models"
	"ttshub/internal/repository"
	"ttshub/internal/utils"

	"go.uber.org/zap"
)

// TTSServiceImpl TTS服务实现
type TTSServiceImpl struct {
	channelRepo repository.ChannelRepository
	configCache repository.ConfigCache
}

// NewTTSService 创建TTS服务实例
func NewTTSService(channelRepo repository.ChannelRepository) TTSService {
	return &TTSServiceImpl{
		channelRepo: channelRepo,
		configCache: repository.GetConfigCache(),
	}
}

// Synthesize 执行文本转语音
func (s *TTSServiceImpl) Synthesize(ctx context.Context, request *models.TTSRequest) (*models.TTSResponse, error) {
	if request == nil {
		return nil, utils.NewBadRequestError("请求参数不能为空")
	}

	// 验证请求
	if err := request.Validate(); err != nil {
		return nil, utils.NewBadRequestError(err.Error())
	}

	// 获取渠道配置
	channel, err := s.getChannelConfig(request.ChannelID)
	if err != nil {
		return nil, utils.NewNotFoundError(fmt.Sprintf("渠道不存在或已禁用: %v", err))
	}

	// 确保渠道已启用
	if channel.Status != "enabled" {
		return nil, utils.NewBadRequestError("渠道未启用")
	}

	// 创建并初始化适配器
	adapterInstance, err := adapter.CreateAdapterByName(channel.Type)
	if err != nil {
		return nil, utils.NewInternalError(fmt.Sprintf("适配器创建失败: %v", err))
	}

	// 初始化适配器配置
	if err := adapterInstance.Init(channel.Config); err != nil {
		return nil, utils.NewInternalError(fmt.Sprintf("适配器初始化失败: %v", err))
	}

	// 验证适配器配置
	if err := adapterInstance.ValidateConfig(); err != nil {
		return nil, utils.NewBadRequestError(fmt.Sprintf("适配器配置无效: %v", err))
	}

	// 转换请求为OpenAI格式
	openaiRequest, err := adapterInstance.ConvertRequest(request)
	if err != nil {
		return nil, utils.NewBadRequestError(fmt.Sprintf("请求转换失败: %v", err))
	}

	// 创建OpenAI客户端
	openaiClient := adapter.NewOpenAIClient()

	// 如果渠道有特定的API密钥配置，使用渠道配置的密钥
	if channel.Type == "openai" {
		// 解析OpenAI配置
		var openaiConfig adapter.OpenAIAdapterConfig
		if err := json.Unmarshal([]byte(channel.Config), &openaiConfig); err == nil && openaiConfig.APIKey != "" {
			// 使用渠道配置的API密钥
			openaiClient = openaiClient.WithConfig(openaiConfig.APIKey, "https://api.openai.com", 30)
		}
	} else if channel.Type == "custom" {
		// 解析自定义渠道配置
		var customConfig adapter.CustomChannelConfig
		if err := json.Unmarshal([]byte(channel.Config), &customConfig); err == nil && customConfig.APIKey != "" {
			// 使用渠道配置的API密钥
			openaiClient = openaiClient.WithConfig(customConfig.APIKey, "https://api.openai.com", 30)
		}
	}

	// 执行TTS请求
	response, err := openaiClient.Synthesize(ctx, openaiRequest)
	if err != nil {
		utils.Error("TTS合成失败", zap.Error(err), zap.Uint("channel_id", channel.ID), zap.String("channel_name", channel.Name))
		return nil, utils.NewOpenAIError(fmt.Sprintf("TTS合成失败: %v", err), err)
	}

	// 添加渠道信息到响应
	response.ChannelID = channel.ID
	response.ChannelName = channel.Name

	utils.Info("TTS合成成功", zap.Uint("channel_id", channel.ID), zap.String("channel_name", channel.Name), zap.Float64("duration", response.Duration))
	return response, nil
}

// ListChannels 获取所有渠道列表
func (s *TTSServiceImpl) ListChannels() ([]*models.ChannelConfig, error) {
	return s.channelRepo.ListAll()
}

// GetChannel 获取单个渠道配置
func (s *TTSServiceImpl) GetChannel(id uint) (*models.ChannelConfig, error) {
	return s.getChannelConfig(id)
}

// CreateChannel 创建新渠道
func (s *TTSServiceImpl) CreateChannel(channel *models.ChannelConfig) error {
	// 验证渠道类型
	if !adapter.IsAdapterRegistered(channel.Type) {
		return utils.NewBadRequestError(fmt.Sprintf("不支持的渠道类型: %s", channel.Type))
	}

	// 创建渠道
	if err := s.channelRepo.Create(channel); err != nil {
		return err
	}

	// 更新缓存
	s.configCache.Set(channel.ID, channel)

	return nil
}

// UpdateChannel 更新渠道配置
func (s *TTSServiceImpl) UpdateChannel(channel *models.ChannelConfig) error {
	// 验证渠道存在
	existing, err := s.getChannelConfig(channel.ID)
	if err != nil {
		return err
	}

	// 不允许更改已存在渠道的类型
	if existing.Type != channel.Type {
		return utils.NewBadRequestError("渠道类型不允许更改")
	}

	// 更新渠道
	if err := s.channelRepo.Update(channel); err != nil {
		return err
	}

	// 更新缓存
	s.configCache.Update(channel.ID, channel)

	return nil
}

// DeleteChannel 删除渠道
func (s *TTSServiceImpl) DeleteChannel(id uint) error {
	// 验证渠道存在
	_, err := s.getChannelConfig(id)
	if err != nil {
		return err
	}

	// 删除渠道
	if err := s.channelRepo.Delete(id); err != nil {
		return err
	}

	// 从缓存中删除
	s.configCache.Delete(id)

	return nil
}

// RegisterAdapter 注册新的适配器
func (s *TTSServiceImpl) RegisterAdapter(name string, factory adapter.AdapterFactory) error {
	if name == "" || factory == nil {
		return utils.NewBadRequestError("适配器名称和工厂函数不能为空")
	}

	adapter.GlobalRegistry.Register(name, factory)
	return nil
}

// getChannelConfig 获取渠道配置（优先从缓存）
func (s *TTSServiceImpl) getChannelConfig(id uint) (*models.ChannelConfig, error) {
	// 尝试从缓存获取
	if config, found := s.configCache.Get(id); found {
		return config, nil
	}

	// 从数据库获取
	config, err := s.channelRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// 更新缓存
	s.configCache.Set(id, config)

	return config, nil
}

// ValidateChannelConfig 验证渠道配置
func (s *TTSServiceImpl) ValidateChannelConfig(channelType string, config string) error {
	// 验证渠道类型
	if !adapter.IsAdapterRegistered(channelType) {
		return utils.NewBadRequestError(fmt.Sprintf("不支持的渠道类型: %s", channelType))
	}

	// 创建适配器实例
	adapterInstance, err := adapter.CreateAdapterByName(channelType)
	if err != nil {
		return err
	}

	// 初始化并验证配置
	if err := adapterInstance.Init(config); err != nil {
		return fmt.Errorf("配置初始化失败: %w", err)
	}

	return adapterInstance.ValidateConfig()
}

// GetAvailableAdapterTypes 获取可用的适配器类型
func (s *TTSServiceImpl) GetAvailableAdapterTypes() []string {
	return adapter.GetAllAdapterNames()
}

// Ping 测试服务可用性
func (s *TTSServiceImpl) Ping() error {
	// 简单测试，确保服务可以正常访问数据库
	_, err := s.channelRepo.ListAll()
	return err
}