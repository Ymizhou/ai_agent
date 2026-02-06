package user

// UserUpdateRequest 用户更新请求
type UserUpdateRequest struct {
	ID          int64  `json:"id" binding:"required"`           // id
	UserName    string `json:"userName" binding:"omitempty"`    // 用户昵称
	UserAvatar  string `json:"userAvatar" binding:"omitempty"`  // 用户头像
	UserProfile string `json:"userProfile" binding:"omitempty"` // 简介
	UserRole    string `json:"userRole" binding:"omitempty"`    // 用户角色：user/admin
}
