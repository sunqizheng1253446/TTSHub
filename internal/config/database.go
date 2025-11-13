package config

import (
	"log"

	"ttshub/internal/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DB 全局数据库实例
var DB *gorm.DB

// InitDatabase 初始化数据库连接
func InitDatabase(dbPath string) error {
	// 配置GORM日志
	logConfig := logger.Config{
		SlowThreshold: 200, // 慢查询阈值，单位毫秒
		LogLevel:      logger.Info,
	}

	// 连接SQLite数据库
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		Logger: logger.New(
			log.New(log.Writer(), "[DB] ", log.LstdFlags),
			logConfig,
		),
	})
	if err != nil {
		return err
	}

	// 设置连接池参数
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}

	// 设置最大空闲连接数
	sqlDB.SetMaxIdleConns(10)
	// 设置最大打开连接数
	sqlDB.SetMaxOpenConns(100)

	// 执行数据库迁移
	if err := migrateDatabase(db); err != nil {
		return err
	}

	DB = db
	log.Println("数据库连接成功")
	return nil
}

// 执行数据库迁移
func migrateDatabase(db *gorm.DB) error {
	// 自动迁移表结构
	if err := db.AutoMigrate(
		&models.ChannelConfig{},
	); err != nil {
		return err
	}

	// 初始化默认数据
	if err := initDefaultData(db); err != nil {
		return err
	}

	log.Println("数据库迁移完成")
	return nil
}

// 初始化默认数据
func initDefaultData(db *gorm.DB) error {
	// 检查是否已有渠道配置
	var count int64
	db.Model(&models.ChannelConfig{}).Count(&count)
	if count > 0 {
		return nil
	}

	// 创建默认OpenAI渠道配置
	defaultChannel := &models.ChannelConfig{
		Name:   "默认OpenAI渠道",
		Type:   "openai",
		Config: `{"api_key":"","model":"tts-1","voice":"alloy"}`,
		Status: 1,
	}

	// 创建默认自定义渠道配置模板
	customChannel := &models.ChannelConfig{
		Name:   "自定义渠道模板",
		Type:   "custom",
		Config: `{"param1":"value1","param2":"value2"}`,
		Status: 0,
	}

	// 保存默认配置
	if err := db.Create(defaultChannel).Error; err != nil {
		return err
	}
	if err := db.Create(customChannel).Error; err != nil {
		return err
	}

	log.Println("默认数据初始化完成")
	return nil
}

// GetDB 获取数据库实例
func GetDB() *gorm.DB {
	return DB
}