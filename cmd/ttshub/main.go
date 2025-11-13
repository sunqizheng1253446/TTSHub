package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"ttshub/internal/config"
	"ttshub/internal/handlers"
	"ttshub/internal/utils"

	"go.uber.org/zap"
)

func main() {
	// 初始化应用
	initApp()

	// 加载配置
	cfg := config.GetConfig()

	// 初始化数据库
	if err := config.InitDatabase(); err != nil {
		utils.Fatal("数据库初始化失败", zap.Error(err))
	}

	// 设置路由
	router := handlers.SetupRouter()

	// 创建HTTP服务器
	server := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port),
		Handler:      router,
		ReadTimeout:  time.Duration(cfg.Server.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(cfg.Server.WriteTimeout) * time.Second,
		IdleTimeout:  time.Duration(cfg.Server.IdleTimeout) * time.Second,
	}

	// 启动服务器（非阻塞）
	go func() {
		utils.Info("服务器启动中", zap.String("addr", server.Addr))
		utils.Info("服务访问地址", zap.String("url", fmt.Sprintf("http://%s:%d", cfg.Server.Host, cfg.Server.Port)))
		utils.Info("API文档地址", zap.String("url", fmt.Sprintf("http://%s:%d/api/v1", cfg.Server.Host, cfg.Server.Port)))

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			utils.Fatal("服务器启动失败", zap.Error(err))
		}
	}()

	// 等待中断信号以优雅地关闭服务器
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	utils.Info("正在关闭服务器...")

	// 设置关闭超时
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 关闭服务器
	if err := server.Shutdown(ctx); err != nil {
		utils.Error("服务器关闭失败", zap.Error(err))
	} else {
		utils.Info("服务器已优雅关闭")
	}

	// 关闭日志
	utils.SyncLogger()
}

// initApp 初始化应用
func initApp() {
	// 初始化日志
	if err := utils.InitLogger(); err != nil {
		log.Fatalf("日志初始化失败: %v", err)
	}

	// 加载配置
	if err := config.LoadConfig(); err != nil {
		utils.Fatal("配置加载失败", zap.Error(err))
	}

	// 确保数据目录存在
	if err := config.EnsureDataDir(); err != nil {
		utils.Fatal("数据目录创建失败", zap.Error(err))
	}

	utils.Info("TTSHub服务启动")
	utils.Info("版本: 1.0.0")
	utils.Info("作者: TTSHub Team")
}