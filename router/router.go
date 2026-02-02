package router

import (
	"aicode/config"
	"aicode/controller"
	"aicode/docs"
	"aicode/router/middleware"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type HttpRouter struct {
	engine           *gin.Engine
	healthController *controller.HealthController
	userController   *controller.UserController
}

// SetupRouter 设置路由
func SetupRouter(
	healthController *controller.HealthController,
	userController *controller.UserController,
) *gin.Engine {
	cfg := config.GetConfig()
	hr := &HttpRouter{
		healthController: healthController,
		userController:   userController,
	}
	// 创建 Gin 引擎
	r := gin.New()
	hr.engine = r

	logrus.SetFormatter(&logrus.TextFormatter{})
	logLevel, err := logrus.ParseLevel(cfg.Server.LogLevel)
	if err != nil {
		logrus.Panicf("解析日志级别失败: %s, 使用默认级别: %s", err.Error(), logrus.DebugLevel)
		logrus.SetLevel(logrus.DebugLevel)
	} else {
		logrus.SetLevel(logLevel)
	}

	// 添加全局中间件
	r.Use(middleware.CORSMiddleware(), // CORS 跨域
		middleware.GlobalErrorHandler()) // 全局异常处理

	// 获取配置
	contextPath := cfg.Server.ContextPath
	if contextPath == "" {
		contextPath = "/"
	}

	// docs 接口文档（注册到根路由，不受 contextPath 影响）
	if gin.IsDebugging() {
		docs.SwaggerInfo.BasePath = contextPath
		r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	}

	// 创建主路由组
	apiGroup := r.Group(contextPath)

	{
		// 注册健康检查路由
		health := apiGroup.Group("/health")
		hr.healthController.RegisterRoutes(health)
	}
	{
		// 注册用户控制器
		user := apiGroup.Group("/user")
		hr.userController.RegisterRoutes(user)
	}

	return r
}
