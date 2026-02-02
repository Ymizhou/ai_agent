//go:build wireinject
// +build wireinject

package cmd

import (
	"aicode/controller"
	"aicode/mapper"
	"aicode/router"
	"aicode/service/impl"

	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

// App 应用程序结构体，包含所有需要的组件
type App struct {
	Router           *gin.Engine
	UserController   *controller.UserController
	HealthController *controller.HealthController
}

// wireSet 定义所有的provider集合
var wireSet = wire.NewSet(
	MustProvideConfig,
	MustProvideDB,
	router.SetupRouter,
	mapper.NewUserMapper,
	impl.NewUserService,
	controller.NewUserController,
	controller.NewHealthController,
)

// InitializeApp 初始化应用程序（此函数会被wire生成）
func InitializeApp() (*App, error) {
	wire.Build(
		wireSet,
		wire.Struct(new(App), "*"),
	)
	return nil, nil
}
