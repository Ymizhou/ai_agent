package controller

import (
	"aicode/common"
	"net/http"

	"github.com/gin-gonic/gin"
)

// HealthController 健康检查控制器
type HealthController struct{}

// NewHealthController 创建健康检查控制器
func NewHealthController() *HealthController {
	return &HealthController{}
}

// RegisterRoutes 注册路由
func (h *HealthController) RegisterRoutes(r *gin.RouterGroup) {
	{
		r.GET("/", h.HealthCheck)
	}
}

// HealthCheck 健康检查接口
// @Summary 健康检查
// @Description 检查服务是否正常运行
// @Tags 健康检查
// @Accept json
// @Produce json
// @Success 200 {object} common.BaseResponse[string]
// @Router /health/ [get]
func (h *HealthController) HealthCheck(c *gin.Context) {
	response := common.Success("ok")
	c.JSON(http.StatusOK, response)
}
