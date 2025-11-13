package repository

import (
	"encoding/json"
	"fmt"
	"ttshub/internal/models"
	"ttshub/internal/utils"

	"gorm.io/gorm"
	"go.uber.org/zap"
)

// ChannelRepository 渠道配置仓库接口
type ChannelRepository interface {
	Create(channel *models.ChannelConfig) error
	Update(channel *models.ChannelConfig) error
	Delete(id uint) error
	GetByID(id uint) (*models.ChannelConfig, error)
	GetByName(name string) (*models.ChannelConfig, error)
	GetByType(channelType string) ([]*models.ChannelConfig, error)
	ListAll() ([]*models.ChannelConfig, error)
	ListEnabled() ([]*models.ChannelConfig, error)
	UpdateStatus(id uint, status int) error
	ValidateUniqueName(name string, excludeID uint) error
}

// channelRepository 渠道配置仓库实现
type channelRepository struct {
	db *gorm.DB
}

// NewChannelRepository 创建渠道配置仓库实例
func NewChannelRepository(db *gorm.DB) ChannelRepository {
	return &channelRepository{
		db: db,
	}
}

// Create 创建新的渠道配置
func (r *channelRepository) Create(channel *models.ChannelConfig) error {
	if channel == nil {
		return fmt.Errorf("渠道配置不能为空")
	}

	// 验证名称唯一性
	if err := r.ValidateUniqueName(channel.Name, 0); err != nil {
		return err
	}

	// 验证配置JSON格式
	if channel.Config != "" {
		var configMap map[string]interface{}
		if err := json.Unmarshal([]byte(channel.Config), &configMap); err != nil {
			return fmt.Errorf("配置格式无效: %w", err)
		}
	}

	// 设置默认值
	if channel.Status == 0 {
		channel.Status = 1 // 1: 启用
	}

	if channel.Type == "" {
		channel.Type = "custom"
	}

	// 执行创建
	result := r.db.Create(channel)
	if result.Error != nil {
		utils.Error("创建渠道配置失败", zap.Error(result.Error), zap.String("name", channel.Name))
		return fmt.Errorf("创建渠道配置失败: %w", result.Error)
	}

	utils.Info("渠道配置创建成功", zap.Uint("id", channel.ID), zap.String("name", channel.Name))
	return nil
}

// Update 更新渠道配置
func (r *channelRepository) Update(channel *models.ChannelConfig) error {
	if channel == nil || channel.ID == 0 {
		return fmt.Errorf("无效的渠道配置或ID")
	}

	// 验证名称唯一性
	if err := r.ValidateUniqueName(channel.Name, channel.ID); err != nil {
		return err
	}

	// 验证配置JSON格式
	if channel.Config != "" {
		var configMap map[string]interface{}
		if err := json.Unmarshal([]byte(channel.Config), &configMap); err != nil {
			return fmt.Errorf("配置格式无效: %w", err)
		}
	}

	// 检查是否存在
	existing := &models.ChannelConfig{}
	if err := r.db.First(existing, channel.ID).Error; err != nil {
		return fmt.Errorf("渠道配置不存在")
	}

	// 执行更新
	result := r.db.Save(channel)
	if result.Error != nil {
		utils.Error("更新渠道配置失败", zap.Error(result.Error), zap.Uint("id", channel.ID))
		return fmt.Errorf("更新渠道配置失败: %w", result.Error)
	}

	utils.Info("渠道配置更新成功", zap.Uint("id", channel.ID), zap.String("name", channel.Name))
	return nil
}

// Delete 删除渠道配置
func (r *channelRepository) Delete(id uint) error {
	if id == 0 {
		return fmt.Errorf("无效的渠道ID")
	}

	// 执行删除
	result := r.db.Delete(&models.ChannelConfig{}, id)
	if result.Error != nil {
		utils.Error("删除渠道配置失败", zap.Error(result.Error), zap.Uint("id", id))
		return fmt.Errorf("删除渠道配置失败: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("渠道配置不存在")
	}

	utils.Info("渠道配置删除成功", zap.Uint("id", id))
	return nil
}

// GetByID 根据ID获取渠道配置
func (r *channelRepository) GetByID(id uint) (*models.ChannelConfig, error) {
	if id == 0 {
		return nil, fmt.Errorf("无效的渠道ID")
	}

	channel := &models.ChannelConfig{}
	result := r.db.First(channel, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("渠道配置不存在")
		}
		utils.Error("获取渠道配置失败", zap.Error(result.Error), zap.Uint("id", id))
		return nil, fmt.Errorf("获取渠道配置失败: %w", result.Error)
	}

	return channel, nil
}

// GetByName 根据名称获取渠道配置
func (r *channelRepository) GetByName(name string) (*models.ChannelConfig, error) {
	if name == "" {
		return nil, fmt.Errorf("渠道名称不能为空")
	}

	channel := &models.ChannelConfig{}
	result := r.db.Where("name = ?", name).First(channel)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("渠道配置不存在")
		}
		utils.Error("获取渠道配置失败", zap.Error(result.Error), zap.String("name", name))
		return nil, fmt.Errorf("获取渠道配置失败: %w", result.Error)
	}

	return channel, nil
}

// GetByType 根据类型获取渠道配置列表
func (r *channelRepository) GetByType(channelType string) ([]*models.ChannelConfig, error) {
	channels := []*models.ChannelConfig{}
	result := r.db.Where("type = ?", channelType).Find(&channels)
	if result.Error != nil {
		utils.Error("获取渠道配置列表失败", zap.Error(result.Error), zap.String("type", channelType))
		return nil, fmt.Errorf("获取渠道配置列表失败: %w", result.Error)
	}

	return channels, nil
}

// ListAll 获取所有渠道配置
func (r *channelRepository) ListAll() ([]*models.ChannelConfig, error) {
	channels := []*models.ChannelConfig{}
	result := r.db.Order("created_at DESC").Find(&channels)
	if result.Error != nil {
		utils.Error("获取所有渠道配置失败", zap.Error(result.Error))
		return nil, fmt.Errorf("获取所有渠道配置失败: %w", result.Error)
	}

	return channels, nil
}

// ListEnabled 获取所有启用的渠道配置
func (r *channelRepository) ListEnabled() ([]*models.ChannelConfig, error) {
	channels := []*models.ChannelConfig{}
	result := r.db.Where("status = ?", 1).Order("created_at DESC").Find(&channels)
	if result.Error != nil {
		utils.Error("获取启用的渠道配置失败", zap.Error(result.Error))
		return nil, fmt.Errorf("获取启用的渠道配置失败: %w", result.Error)
	}

	return channels, nil
}

// UpdateStatus 更新渠道状态
func (r *channelRepository) UpdateStatus(id uint, status int) error {
	if id == 0 {
		return fmt.Errorf("无效的渠道ID")
	}

	// 验证状态值
	validStatus := map[int]bool{
		1: true, // 启用
		0: true, // 禁用
	}
	if !validStatus[status] {
		return fmt.Errorf("无效的状态值: %d", status)
	}

	// 执行更新
	result := r.db.Model(&models.ChannelConfig{}).Where("id = ?", id).Update("status", status)
	if result.Error != nil {
		utils.Error("更新渠道状态失败", zap.Error(result.Error), zap.Uint("id", id), zap.Int("status", status))
		return fmt.Errorf("更新渠道状态失败: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("渠道配置不存在")
	}

	utils.Info("渠道状态更新成功", zap.Uint("id", id), zap.Int("status", status))
	return nil
}

// ValidateUniqueName 验证渠道名称唯一性
func (r *channelRepository) ValidateUniqueName(name string, excludeID uint) error {
	if name == "" {
		return fmt.Errorf("渠道名称不能为空")
	}

	query := r.db.Model(&models.ChannelConfig{}).Where("name = ?", name)
	if excludeID > 0 {
		query = query.Where("id != ?", excludeID)
	}

	var count int64
	if err := query.Count(&count).Error; err != nil {
		return fmt.Errorf("验证名称唯一性失败: %w", err)
	}

	if count > 0 {
		return fmt.Errorf("渠道名称已存在: %s", name)
	}

	return nil
}