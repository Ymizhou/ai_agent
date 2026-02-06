package user

// UserLoginRequest 用户登录请求
type UserLoginRequest struct {
	UserAccount  string `json:"userAccount" binding:"required"`  // 账号
	UserPassword string `json:"userPassword" binding:"required"` // 密码
}
