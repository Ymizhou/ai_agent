package middleware

import (
	"context"
	"net/http"

	"aicode/constant"
	"aicode/internal/common"
	"aicode/internal/exception"
	"aicode/internal/model/entity"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// skipRoute 描述一条不需要登录校验的路由（方法 + 完整路径）
type skipRoute struct {
	Method string
	Path   string
}

// skipRoutes 白名单列表，路径须含 rootPath 前缀，与 config.yml server.root_path 保持一致
// 后续需要跳过登录校验的接口，直接在此处追加即可
var skipRoutes = []skipRoute{
	{http.MethodPost, "/api/v1/user/register"},
	{http.MethodPost, "/api/v1/user/login"},
	{http.MethodGet, "/swagger/*any"},
}

// AuthMiddleware 登录态校验中间件
//   - 白名单内的路由直接放行
//   - 其余路由从 session 中读取登录用户，未登录则返回 NotLoginError
//   - 登录用户同时写入 gin context（c.Set）和 request context（context.WithValue），
//     后续业务层可通过 c.Get(constant.UserLoginState) 或 ctx.Value(constant.UserLoginState) 取用
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 白名单检查
		for _, route := range skipRoutes {
			if c.Request.Method == route.Method && c.FullPath() == route.Path {
				c.Next()
				return
			}
		}

		// 从 session 中获取登录用户
		session := sessions.Default(c)
		userObj := session.Get(constant.UserLoginState)
		if userObj == nil {
			c.JSON(http.StatusOK, common.Error(exception.NotLoginError))
			c.Abort()
			return
		}

		loginUser, ok := userObj.(*entity.User)
		if !ok || loginUser == nil || loginUser.ID == 0 {
			c.JSON(http.StatusOK, common.Error(exception.NotLoginError))
			c.Abort()
			return
		}

		// 将用户信息写入 gin context
		c.Set(constant.UserLoginState, loginUser)
		// 将用户信息写入 request context，方便 service 层通过 ctx 取用
		ctx := context.WithValue(c.Request.Context(), constant.UserLoginState, loginUser)
		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}
