package handlers

import (
	"net/http"
	"ttshub/internal/config"
	"ttshub/internal/repository"
	"ttshub/internal/service"
	"ttshub/internal/utils"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// SetupRouter 设置路由
func SetupRouter() *gin.Engine {
	// 设置运行模式
	if config.GetConfig().Server.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	// 创建路由
	router := gin.New()

	// 添加中间件
	router.Use(gin.Recovery())
	router.Use(loggingMiddleware())
	router.Use(corsMiddleware())

	// 初始化服务和处理器
	db := config.GetDB()
	channelRepo := repository.NewChannelRepository(db)
	ttsService := service.NewTTSService(channelRepo)
	ttsHandler := NewTTSHandler(ttsService)

	// 预加载配置缓存
	if err := repository.PreloadConfigCache(channelRepo); err != nil {
		utils.Warn("配置缓存预加载失败", zap.Error(err))
	}

	// 健康检查路由
	router.GET("/ping", ttsHandler.Ping)
	router.GET("/health", ttsHandler.Health)

	// API v1 路由组
	v1 := router.Group("/api/v1")
	{
		// TTS 转换相关路由
		v1.POST("/tts/synthesize", ttsHandler.Synthesize)

		// 渠道管理相关路由
		v1.GET("/channels", ttsHandler.ListChannels)
		v1.GET("/channels/:id", ttsHandler.GetChannel)
		v1.POST("/channels", ttsHandler.CreateChannel)
		v1.PUT("/channels/:id", ttsHandler.UpdateChannel)
		v1.DELETE("/channels/:id", ttsHandler.DeleteChannel)

		// 配置验证路由
		v1.POST("/channels/validate", ttsHandler.ValidateChannelConfig)

		// 适配器类型路由
		v1.GET("/adapters/types", ttsHandler.GetAvailableAdapterTypes)
	}

	// 静态文件服务（用于网页控制台）
	router.Static("/static", "./static")

	// 前端网页入口（SPA）
	router.GET("/*any", func(c *gin.Context) {
		// 检查是否为API请求
		if c.Request.URL.Path == "/" || !isAPIRequest(c.Request.URL.Path) {
			// 返回前端入口页面
			c.File("./static/index.html")
		} else {
			// 404 处理
			c.JSON(http.StatusNotFound, gin.H{
				"code":    404,
				"message": "接口不存在",
			})
		}
	})

	return router
}

// loggingMiddleware 日志中间件
func loggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 记录请求开始时间
		start := utils.Now()

		// 处理请求
		c.Next()

		// 计算请求耗时
		latency := utils.Since(start)

		// 记录日志
		utils.Info("HTTP请求",
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.Int("status", c.Writer.Status()),
			zap.Duration("latency", latency),
			zap.String("client_ip", c.ClientIP()),
			zap.String("user_agent", c.Request.UserAgent()),
		)

		// 记录错误（如果有）
		if len(c.Errors) > 0 {
			for _, e := range c.Errors {
				utils.Error("请求处理错误",
					zap.String("method", c.Request.Method),
					zap.String("path", c.Request.URL.Path),
					zap.Error(e.Err),
				)
			}
		}
	}
}

// corsMiddleware CORS中间件
func corsMiddleware() gin.HandlerFunc {
	cfg := config.GetConfig()
	
	corsConfig := cors.Config{
		AllowOrigins:     cfg.Server.AllowedOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization"},
		ExposeHeaders:    []string{"Content-Length", "X-TTS-Duration", "X-TTS-Channel"},
		AllowCredentials: true,
		MaxAge:           86400,
	}

	return cors.New(corsConfig)
}

// isAPIRequest 检查是否为API请求
func isAPIRequest(path string) bool {
	// 检查是否以/api开头
	return len(path) >= 5 && path[:5] == "/api/"
}