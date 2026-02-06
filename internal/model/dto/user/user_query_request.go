package user

import "aicode/internal/common"

// UserQueryRequest 用户查询请求
type UserQueryRequest struct {
	common.PageRequest
	ID          *int64 `json:"id" form:"id"`                   // id
	UserName    string `json:"userName" form:"userName"`       // 用户昵称
	UserAccount string `json:"userAccount" form:"userAccount"` // 账号
	UserProfile string `json:"userProfile" form:"userProfile"` // 简介
	UserRole    string `json:"userRole" form:"userRole"`       // 用户角色：user/admin/ban
}
