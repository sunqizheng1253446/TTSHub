package adapter

import (
	"fmt"
	"sync"
	"ttshub/internal/utils"

	"go.uber.org/zap"
)

// AdapterRegistry 适配器注册表接口
type AdapterRegistry interface {
	Register(name string, factory AdapterFactory)
	Unregister(name string)
	Get(name string) (TTSAdapter, error)
	List() []string
	Has(name string) bool
}

// AdapterFactory 适配器工厂函数类型
type AdapterFactory func() TTSAdapter

// adapterRegistry 适配器注册表实现
type adapterRegistry struct {
	adapters map[string]AdapterFactory
	mu       sync.RWMutex
}

// GlobalRegistry 全局适配器注册表实例
var GlobalRegistry AdapterRegistry

// init 初始化全局注册表
func init() {
	GlobalRegistry = NewAdapterRegistry()
	
	// 注册内置适配器
	RegisterBuiltinAdapters(GlobalRegistry)
}

// NewAdapterRegistry 创建适配器注册表实例
func NewAdapterRegistry() AdapterRegistry {
	return &adapterRegistry{
		adapters: make(map[string]AdapterFactory),
	}
}

// Register 注册适配器
func (r *adapterRegistry) Register(name string, factory AdapterFactory) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if name == "" || factory == nil {
		utils.Error("无效的适配器注册参数", zap.String("name", name))
		return
	}

	_, exists := r.adapters[name]
	r.adapters[name] = factory

	if exists {
		utils.Warn("适配器已被覆盖", zap.String("name", name))
	} else {
		utils.Info("适配器注册成功", zap.String("name", name))
	}
}

// Unregister 注销适配器
func (r *adapterRegistry) Unregister(name string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.adapters[name]; exists {
		delete(r.adapters, name)
		utils.Info("适配器注销成功", zap.String("name", name))
	} else {
		utils.Warn("尝试注销不存在的适配器", zap.String("name", name))
	}
}

// Get 获取适配器实例
func (r *adapterRegistry) Get(name string) (TTSAdapter, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	factory, exists := r.adapters[name]
	if !exists {
		return nil, fmt.Errorf("适配器不存在: %s", name)
	}

	adapter := factory()
	if adapter == nil {
		return nil, fmt.Errorf("适配器创建失败: %s", name)
	}

	return adapter, nil
}

// List 列出所有注册的适配器名称
func (r *adapterRegistry) List() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	names := make([]string, 0, len(r.adapters))
	for name := range r.adapters {
		names = append(names, name)
	}

	return names
}

// Has 检查适配器是否已注册
func (r *adapterRegistry) Has(name string) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()

	_, exists := r.adapters[name]
	return exists
}

// RegisterBuiltinAdapters 注册内置适配器
func RegisterBuiltinAdapters(registry AdapterRegistry) {
	// 注册OpenAI适配器
	registry.Register("openai", func() TTSAdapter {
		return NewOpenAIAdapter()
	})

	// 注册默认的自定义渠道适配器
	registry.Register("custom", func() TTSAdapter {
		return NewCustomChannelAdapter("custom")
	})

	utils.Info("内置适配器注册完成", zap.Int("count", 2))
}

// CreateAdapterByName 根据名称创建适配器实例
func CreateAdapterByName(name string) (TTSAdapter, error) {
	return GlobalRegistry.Get(name)
}

// GetAllAdapterNames 获取所有适配器名称
func GetAllAdapterNames() []string {
	return GlobalRegistry.List()
}

// IsAdapterRegistered 检查适配器是否已注册
func IsAdapterRegistered(name string) bool {
	return GlobalRegistry.Has(name)
}