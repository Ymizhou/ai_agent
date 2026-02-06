package user

// UserRegisterRequest 用户注册请求
type UserRegisterRequest struct {
	UserAccount   string `json:"userAccount" binding:"required"`   // 账号
	UserPassword  string `json:"userPassword" binding:"required"`  // 密码
	CheckPassword string `json:"checkPassword" binding:"required"` // 确认密码
}
