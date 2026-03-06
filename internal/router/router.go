package router

import (
	"aicode/config"
	"aicode/docs"
	"aicode/internal/controller"
	"aicode/internal/router/middleware"

	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type HttpRouter struct {
	healthController *controller.HealthController
	userController   *controller.UserController
	aiController     *controller.AIController
	aiCodeController *controller.AICodeController
}

// SetupRouter 设置路由
func SetupRouter(
	healthController *controller.HealthController,
	userController *controller.UserController,
	aiController *controller.AIController,
	aiCodeController *controller.AICodeController,
) *gin.Engine {
	cfg := config.GetConfig()
	hr := &HttpRouter{
		healthController: healthController,
		userController:   userController,
		aiController:     aiController,
		aiCodeController: aiCodeController,
	}
	// 创建 Gin 引擎
	r := gin.New()

	// 添加全局中间件
	r.Use(middleware.CORSMiddleware(), // CORS 跨域
		middleware.FormatLog(),          // 格式化日志
		middleware.GlobalErrorHandler()) // 全局异常处理

	// 获取配置
	rootPath := cfg.Server.RootPath
	if rootPath == "" {
		rootPath = "/"
	}

	// docs 接口文档（注册到根路由，不受 rootPath 影响）
	if gin.IsDebugging() {
		docs.SwaggerInfo.BasePath = rootPath
		r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	}

	// 创建主路由组
	apiGroup := r.Group(rootPath)

	// 注册健康检查路由
	{
		health := apiGroup.Group("/health")
		hr.healthController.RegisterRoutes(health)
	}

	// 注册用户路由
	{
		user := apiGroup.Group("/user")
		hr.userController.RegisterRoutes(user)
	}

	// 注册ai交互路由
	{
		aiChat := apiGroup.Group("/ai_chat")
		hr.aiController.RegisterRoutes(aiChat)
	}
	// 代码生成
	{
		aiCode := apiGroup.Group("/ai_code")
		hr.aiCodeController.RegisterRoutes(aiCode)
	}

	return r
}
