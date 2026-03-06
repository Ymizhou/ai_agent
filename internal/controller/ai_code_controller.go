package controller

import (
	"net/http"

	"aicode/internal/common"
	"aicode/internal/exception"
	"aicode/internal/model/vo"
	"aicode/internal/service"

	"github.com/gin-gonic/gin"
)

// AIController ai控制层
type AICodeController struct {
	aiCodeService service.AICodeService
}

// NewAIController 创建ai控制器
func NewAICodeController(aiCodeService service.AICodeService) *AICodeController {
	return &AICodeController{
		aiCodeService: aiCodeService,
	}
}

// RegisterRoutes 注册路由
func (ctrl *AICodeController) RegisterRoutes(r *gin.RouterGroup) {
	{
		// 代码生成
		r.POST("/gen", ctrl.CodeGenerate)
		r.POST("/gen/stream", ctrl.CodeGenerateStream)
	}
}

// CodeGenerateStream 代码生成流式
// @Summary 代码生成流式
// @Description 代码生成流式
// @Tags ai_code模块
// @Accept json
// @Produce json
// @Param request body vo.AICodeRequest true "代码生成请求"
// @Success 200 {object} common.BaseResponse[string]
// @Router /ai_code/gen/stream [post]
func (ctrl *AICodeController) CodeGenerateStream(c *gin.Context) {
	// 绑定请求参数
	req := &vo.AICodeRequest{}
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

	// 初始化 channel，service 层异步将流数据写入该 channel
	ch := make(chan vo.CodeStreamResult, 32)
	if err := ctrl.aiCodeService.CodeGenerateStream(ctx, req, ch); err != nil {
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

	// 监听 channel，将增量内容实时推送给客户端
	for result := range ch {
		if result.Err != nil {
			c.SSEvent("error", gin.H{"error": result.Err.Error()})
			flusher.Flush()
			return
		}
		c.SSEvent("message", gin.H{
			"type":    "data",
			"content": result.Content,
		})
		flusher.Flush()
	}

	// channel 关闭即代表 stream 已全部消费完毕
	c.SSEvent("message", gin.H{
		"type":    "end",
		"content": "",
	})
	flusher.Flush()
}

// CodeGenerate 代码生成
// @Summary 代码生成
// @Description 代码生成
// @Tags ai_code模块
// @Accept json
// @Produce json
// @Param request body vo.AICodeRequest true "代码生成请求"
// @Success 200 {object} common.BaseResponse[string]
// @Router /ai_code/gen [post]
func (ctrl *AICodeController) CodeGenerate(c *gin.Context) {
	// 绑定请求参数
	req := &vo.AICodeRequest{}
	if err := c.ShouldBindJSON(req); err != nil {
		c.JSON(http.StatusBadRequest, common.Error(exception.ParamsError))
		return
	}
	ctx := c.Request.Context()
	// 调用服务
	message, err := ctrl.aiCodeService.CodeGenerate(ctx, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			common.ErrorWithMessage(exception.OperationError, err.Error()))
		return
	}
	c.JSON(http.StatusOK, common.Success(message))
}
