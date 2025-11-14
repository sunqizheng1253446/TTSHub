package repository

import (
	"sync"
	"ttshub/internal/models"
	"ttshub/internal/utils"

	"go.uber.org/zap"
)

// ConfigCache 配置缓存接口
type ConfigCache interface {
	Set(id uint, config *models.ChannelConfig)
	Get(id uint) (*models.ChannelConfig, bool)
	GetByName(name string) (*models.ChannelConfig, bool)
	Delete(id uint)
	Update(id uint, config *models.ChannelConfig)
	Clear()
	Size() int
	RebuildFromList(channels []*models.ChannelConfig)
}

// configCache 配置缓存实现
type configCache struct {
	byID   map[uint]*models.ChannelConfig
	byName map[string]*models.ChannelConfig
	mu     sync.RWMutex
}

// NewConfigCache 创建配置缓存实例
func NewConfigCache() ConfigCache {
	return &configCache{
		byID:   make(map[uint]*models.ChannelConfig),
		byName: make(map[string]*models.ChannelConfig),
	}
}

// GlobalConfigCache 全局配置缓存实例
var GlobalConfigCache ConfigCache

// init 初始化全局配置缓存
func init() {
	GlobalConfigCache = NewConfigCache()
}

// Set 设置配置缓存
func (c *configCache) Set(id uint, config *models.ChannelConfig) {
	if id == 0 || config == nil {
		return
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	// 复制配置，避免外部修改
	configCopy := *config
	c.byID[id] = &configCopy
	c.byName[config.Name] = &configCopy

	utils.Debug("配置缓存已设置", zap.Uint("id", id), zap.String("name", config.Name))
}

// Get 根据ID获取配置
func (c *configCache) Get(id uint) (*models.ChannelConfig, bool) {
	if id == 0 {
		return nil, false
	}

	c.mu.RLock()
	defer c.mu.RUnlock()

	config, exists := c.byID[id]
	if !exists {
		utils.Debug("配置缓存未命中", zap.Uint("id", id))
		return nil, false
	}

	// 返回副本，避免外部修改
	configCopy := *config
	return &configCopy, true
}

// GetByName 根据名称获取配置
func (c *configCache) GetByName(name string) (*models.ChannelConfig, bool) {
	if name == "" {
		return nil, false
	}

	c.mu.RLock()
	defer c.mu.RUnlock()

	config, exists := c.byName[name]
	if !exists {
		utils.Debug("配置缓存未命中", zap.String("name", name))
		return nil, false
	}

	// 返回副本，避免外部修改
	configCopy := *config
	return &configCopy, true
}

// Delete 删除配置缓存
func (c *configCache) Delete(id uint) {
	if id == 0 {
		return
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	config, exists := c.byID[id]
	if !exists {
		return
	}

	delete(c.byID, id)
	delete(c.byName, config.Name)

	utils.Debug("配置缓存已删除", zap.Uint("id", id), zap.String("name", config.Name))
}

// Update 更新配置缓存
func (c *configCache) Update(id uint, config *models.ChannelConfig) {
	if id == 0 || config == nil || config.ID != id {
		return
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	// 检查旧配置是否存在
	oldConfig, exists := c.byID[id]
	if exists {
		// 如果名称改变了，需要从byName中删除旧名称的映射
		if oldConfig.Name != config.Name {
			delete(c.byName, oldConfig.Name)
		}
	}

	// 设置新配置（复制一份）
	configCopy := *config
	c.byID[id] = &configCopy
	c.byName[config.Name] = &configCopy

	utils.Debug("配置缓存已更新", zap.Uint("id", id), zap.String("name", config.Name))
}

// Clear 清空所有缓存
func (c *configCache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.byID = make(map[uint]*models.ChannelConfig)
	c.byName = make(map[string]*models.ChannelConfig)

	utils.Info("配置缓存已清空")
}

// Size 获取缓存大小
func (c *configCache) Size() int {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return len(c.byID)
}

// RebuildFromList 从列表重建缓存
func (c *configCache) RebuildFromList(channels []*models.ChannelConfig) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// 清空现有缓存
	c.byID = make(map[uint]*models.ChannelConfig)
	c.byName = make(map[string]*models.ChannelConfig)

	// 添加新配置
	for _, channel := range channels {
		if channel != nil {
			configCopy := *channel
			c.byID[channel.ID] = &configCopy
			c.byName[channel.Name] = &configCopy
		}
	}

	utils.Info("配置缓存已从列表重建", zap.Int("size", len(channels)))
}

// GetConfigCache 获取全局配置缓存实例
func GetConfigCache() ConfigCache {
	return GlobalConfigCache
}

// PreloadConfigCache 预加载配置缓存
func PreloadConfigCache(repo ChannelRepository) error {
	// 如果处于无数据库模式，跳过数据库操作
	if config.NoDBMode {
		utils.Info("无数据库模式，跳过预加载配置缓存")
		GlobalConfigCache.RebuildFromList([]*models.ChannelConfig{})
		return nil
	}

	channels, err := repo.ListAll()
	if err != nil {
		return err
	}

	GlobalConfigCache.RebuildFromList(channels)
	return nil
}