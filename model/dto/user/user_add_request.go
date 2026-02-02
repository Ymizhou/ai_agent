package user

// UserAddRequest 用户创建请求
type UserAddRequest struct {
	UserName    string `json:"userName" binding:"omitempty"`    // 用户昵称
	UserAccount string `json:"userAccount" binding:"required"`  // 账号
	UserAvatar  string `json:"userAvatar" binding:"omitempty"`  // 用户头像
	UserProfile string `json:"userProfile" binding:"omitempty"` // 用户简介
	UserRole    string `json:"userRole" binding:"omitempty"`    // 用户角色: user, admin
}
