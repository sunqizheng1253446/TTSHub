package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"ttshub/internal/models"
	"ttshub/internal/service"
	"ttshub/internal/utils"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// TTSHandler TTS处理器
type TTSHandler struct {
	ttsService service.TTSService
}

// NewTTSHandler 创建TTS处理器实例
func NewTTSHandler(ttsService service.TTSService) *TTSHandler {
	return &TTSHandler{
		ttsService: ttsService,
	}
}

// Synthesize 文本转语音接口
func (h *TTSHandler) Synthesize(c *gin.Context) {
	// 解析请求参数
	var request models.TTSRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		utils.Error("请求参数解析失败", zap.Error(err))
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Code:    400,
			Message: "请求参数无效: " + err.Error(),
		})
		return
	}

	// 执行TTS合成
	response, err := h.ttsService.Synthesize(c.Request.Context(), &request)
	if err != nil {
		// 处理错误
		appErr, ok := err.(*utils.AppError)
		if ok {
			c.JSON(appErr.StatusCode(), models.APIResponse{
				Code:    appErr.Code,
				Message: appErr.Message,
			})
		} else {
			c.JSON(http.StatusInternalServerError, models.APIResponse{
				Code:    500,
				Message: "内部服务器错误: " + err.Error(),
			})
		}
		return
	}

	// 设置响应头
	contentType := "audio/mp3"
	switch response.Format {
	case "mp3":
		contentType = "audio/mp3"
	case "opus":
		contentType = "audio/ogg"
	case "aac":
		contentType = "audio/aac"
	case "flac":
		contentType = "audio/flac"
	}

	// 如果是直接返回音频数据的请求
	if c.Query("raw") == "true" {
		c.Header("Content-Type", contentType)
		c.Header("Content-Length", strconv.Itoa(len(response.AudioData)))
		c.Header("X-TTS-Duration", strconv.FormatFloat(response.Duration, 'f', 2, 64))
		c.Header("X-TTS-Channel", response.ChannelName)
		c.Data(http.StatusOK, contentType, response.AudioData)
		return
	}

	// 返回JSON响应（默认）
	c.JSON(http.StatusOK, models.APIResponse{
		Code:    200,
		Message: "success",
		Data: map[string]interface{}{
			"audio_data":   response.AudioData, // 二进制数据，客户端需要base64解码
			"format":       response.Format,
			"duration":     response.Duration,
			"cost":         response.Cost,
			"channel_id":   response.ChannelID,
			"channel_name": response.ChannelName,
		},
	})
}

// ListChannels 获取渠道列表
func (h *TTSHandler) ListChannels(c *gin.Context) {
	channels, err := h.ttsService.ListChannels()
	if err != nil {
		utils.Error("获取渠道列表失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Code:    500,
			Message: "获取渠道列表失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Code:    200,
		Message: "success",
		Data:    channels,
	})
}

// GetChannel 获取单个渠道
func (h *TTSHandler) GetChannel(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Code:    400,
			Message: "无效的渠道ID",
		})
		return
	}

	channel, err := h.ttsService.GetChannel(uint(id))
	if err != nil {
		appErr, ok := err.(*utils.AppError)
		if ok {
			c.JSON(appErr.StatusCode(), models.APIResponse{
				Code:    appErr.Code,
				Message: appErr.Message,
			})
		} else {
			c.JSON(http.StatusNotFound, models.APIResponse{
				Code:    404,
				Message: "渠道不存在",
			})
		}
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Code:    200,
		Message: "success",
		Data:    channel,
	})
}

// CreateChannel 创建新渠道
func (h *TTSHandler) CreateChannel(c *gin.Context) {
	var channel models.ChannelConfig
	if err := c.ShouldBindJSON(&channel); err != nil {
		utils.Error("创建渠道请求参数解析失败", zap.Error(err))
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Code:    400,
			Message: "请求参数无效: " + err.Error(),
		})
		return
	}

	// 验证配置格式
	if channel.Config != "" {
		var configMap map[string]interface{}
		if err := json.Unmarshal([]byte(channel.Config), &configMap); err != nil {
			c.JSON(http.StatusBadRequest, models.APIResponse{
				Code:    400,
				Message: "配置格式无效，必须是有效的JSON字符串",
			})
			return
		}
	}

	// 创建渠道
	if err := h.ttsService.CreateChannel(&channel); err != nil {
		appErr, ok := err.(*utils.AppError)
		if ok {
			c.JSON(appErr.StatusCode(), models.APIResponse{
				Code:    appErr.Code,
				Message: appErr.Message,
			})
		} else {
			c.JSON(http.StatusInternalServerError, models.APIResponse{
				Code:    500,
				Message: "创建渠道失败: " + err.Error(),
			})
		}
		return
	}

	c.JSON(http.StatusCreated, models.APIResponse{
		Code:    201,
		Message: "渠道创建成功",
		Data:    channel,
	})
}

// UpdateChannel 更新渠道
func (h *TTSHandler) UpdateChannel(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Code:    400,
			Message: "无效的渠道ID",
		})
		return
	}

	var channel models.ChannelConfig
	if err := c.ShouldBindJSON(&channel); err != nil {
		utils.Error("更新渠道请求参数解析失败", zap.Error(err))
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Code:    400,
			Message: "请求参数无效: " + err.Error(),
		})
		return
	}

	// 设置ID
	channel.ID = uint(id)

	// 验证配置格式
	if channel.Config != "" {
		var configMap map[string]interface{}
		if err := json.Unmarshal([]byte(channel.Config), &configMap); err != nil {
			c.JSON(http.StatusBadRequest, models.APIResponse{
				Code:    400,
				Message: "配置格式无效，必须是有效的JSON字符串",
			})
			return
		}
	}

	// 更新渠道
	if err := h.ttsService.UpdateChannel(&channel); err != nil {
		appErr, ok := err.(*utils.AppError)
		if ok {
			c.JSON(appErr.StatusCode(), models.APIResponse{
				Code:    appErr.Code,
				Message: appErr.Message,
			})
		} else {
			c.JSON(http.StatusInternalServerError, models.APIResponse{
				Code:    500,
				Message: "更新渠道失败: " + err.Error(),
			})
		}
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Code:    200,
		Message: "渠道更新成功",
		Data:    channel,
	})
}

// DeleteChannel 删除渠道
func (h *TTSHandler) DeleteChannel(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Code:    400,
			Message: "无效的渠道ID",
		})
		return
	}

	// 删除渠道
	if err := h.ttsService.DeleteChannel(uint(id)); err != nil {
		appErr, ok := err.(*utils.AppError)
		if ok {
			c.JSON(appErr.StatusCode(), models.APIResponse{
				Code:    appErr.Code,
				Message: appErr.Message,
			})
		} else {
			c.JSON(http.StatusInternalServerError, models.APIResponse{
				Code:    500,
				Message: "删除渠道失败: " + err.Error(),
			})
		}
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Code:    200,
		Message: "渠道删除成功",
	})
}

// ValidateChannelConfig 验证渠道配置
func (h *TTSHandler) ValidateChannelConfig(c *gin.Context) {
	type ValidateRequest struct {
		Type   string `json:"type" binding:"required"`
		Config string `json:"config"`
	}

	var req ValidateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Code:    400,
			Message: "请求参数无效: " + err.Error(),
		})
		return
	}

	// 验证配置
	if err := h.ttsService.ValidateChannelConfig(req.Type, req.Config); err != nil {
		appErr, ok := err.(*utils.AppError)
		if ok {
			c.JSON(appErr.StatusCode(), models.APIResponse{
				Code:    appErr.Code,
				Message: appErr.Message,
			})
		} else {
			c.JSON(http.StatusBadRequest, models.APIResponse{
				Code:    400,
				Message: "配置验证失败: " + err.Error(),
			})
		}
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Code:    200,
		Message: "配置验证通过",
	})
}

// GetAvailableAdapterTypes 获取可用的适配器类型
func (h *TTSHandler) GetAvailableAdapterTypes(c *gin.Context) {
	types := h.ttsService.GetAvailableAdapterTypes()
	c.JSON(http.StatusOK, models.APIResponse{
		Code:    200,
		Message: "success",
		Data:    types,
	})
}

// Ping 服务健康检查
func (h *TTSHandler) Ping(c *gin.Context) {
	if err := h.ttsService.Ping(); err != nil {
		utils.Error("服务健康检查失败", zap.Error(err))
		c.JSON(http.StatusServiceUnavailable, models.APIResponse{
			Code:    503,
			Message: "服务暂时不可用: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Code:    200,
		Message: "pong",
	})
}

// Health 健康检查（详细信息）
func (h *TTSHandler) Health(c *gin.Context) {
	healthInfo := map[string]interface{}{
		"status":  "healthy",
		"service": "ttshub",
		"version": "1.0.0",
	}

	// 检查服务是否可用
	if err := h.ttsService.Ping(); err != nil {
		healthInfo["status"] = "unhealthy"
		healthInfo["error"] = err.Error()
		c.JSON(http.StatusServiceUnavailable, healthInfo)
		return
	}

	c.JSON(http.StatusOK, healthInfo)
}