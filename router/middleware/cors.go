package middleware

import "github.com/gin-gonic/gin"

// CORSMiddleware CORS 跨域中间件
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 允许所有域名
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		// 允许携带凭证
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		// 允许的请求头
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		// 允许的请求方法
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")
		// 暴露的响应头
		c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")

		// 处理 OPTIONS 预检请求
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
