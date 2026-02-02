package vo

import "time"

// UserVO 脱敏后的用户信息
type UserVO struct {
	ID          int64     `json:"id"`          // id
	UserAccount string    `json:"userAccount"` // 账号
	UserName    string    `json:"userName"`    // 用户昵称
	UserAvatar  string    `json:"userAvatar"`  // 用户头像
	UserProfile string    `json:"userProfile"` // 用户简介
	UserRole    string    `json:"userRole"`    // 用户角色：user/admin
	CreateTime  time.Time `json:"createTime"`  // 创建时间
}
