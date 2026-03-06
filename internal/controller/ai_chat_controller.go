package controller

import (
	"aicode/internal/common"
	"aicode/internal/exception"
	"aicode/internal/model/vo"
	"aicode/internal/service"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

// AIController ai控制层
type AIController struct {
	aiChatService service.AIChatService
}

// NewAIController 创建ai控制器
func NewAIController(aiChatService service.AIChatService) *AIController {
	return &AIController{
		aiChatService: aiChatService,
	}
}

// RegisterRoutes 注册路由
func (ctrl *AIController) RegisterRoutes(r *gin.RouterGroup) {
	{
		// 聊天
		r.POST("/chat/generate", ctrl.AIChatGenerate)
		r.POST("/chat/stream", ctrl.AIChatStream)
	}
}

// AIChatGenerate聊天
// @Summary 聊天接口(非流式)
// @Description 聊天接口(非流式)
// @Tags ai_chat模块
// @Accept json
// @Produce json
// @Param request body vo.AIChatRequest true "聊天请求(非流式)"
// @Success 200 {object} common.BaseResponse[string]
// @Router /ai_chat/chat/generate [post]
func (ctrl *AIController) AIChatGenerate(c *gin.Context) {
	req := &vo.AIChatRequest{}
	if err := c.ShouldBindJSON(req); err != nil {
		c.JSON(http.StatusBadRequest, common.Error(exception.ParamsError))
		return
	}
	// 使用请求的上下文
	ctx := c.Request.Context()
	message, err := ctrl.aiChatService.ChatGenerate(ctx, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			common.ErrorWithMessage(exception.OperationError, err.Error()))
		return
	}
	c.JSON(http.StatusOK, common.Success(message))
}

// AIChatStream聊天
// @Summary 聊天接口(流式)
// @Description 聊天接口(流式)
// @Tags ai_chat模块
// @Accept json
// @Produce json
// @Param request body vo.AIChatRequest true "聊天请求(流式)"
// @Success 200 {object} common.BaseResponse[string]
// @Router /ai_chat/chat/stream [post]
func (ctrl *AIController) AIChatStream(c *gin.Context) {
	req := &vo.AIChatRequest{}
	if err := c.ShouldBindJSON(req); err != nil {
		c.JSON(http.StatusBadRequest, common.Error(exception.ParamsError))
		return
	}
	// 检查是否支持流式输出
	flusher, ok := c.Writer.(http.Flusher)
	if !ok {
		c.JSON(http.StatusInternalServerError,
			common.ErrorWithMessage(exception.SystemError, "Streaming not supported"))
		return
	}
	// 设置 SSE 响应头（必须在写入任何内容之前设置）
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Headers", "Content-Type")

	c.Writer.WriteHeader(http.StatusOK)

	ctx := c.Request.Context()
	streamReader, err := ctrl.aiChatService.ChatStream(ctx, req)
	if err != nil {
		c.SSEvent("error", gin.H{"error": err.Error()})
		flusher.Flush()
		return
	}
	// 发送开始事件
	c.SSEvent("message", gin.H{
		"type":    "start",
		"content": "",
	})
	flusher.Flush()

	// 读取大模型流式返回的数据，并实时发送给客户端
	for {
		msg, err := streamReader.Recv()

		if err != nil {
			if err == io.EOF {
				// 流结束
				c.SSEvent("message", gin.H{
					"type":    "end",
					"content": "",
				})
				flusher.Flush()
				return
			}
			// 发生错误
			c.SSEvent("error", gin.H{"error": err.Error()})
			flusher.Flush()
			return
		}

		// 发送接收到的增量内容
		if msg != nil && msg.Content != "" {
			c.SSEvent("message", gin.H{
				"type":    "data",
				"content": msg.Content,
			})
			flusher.Flush()
		}
	}
}
