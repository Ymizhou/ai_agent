package entity

import (
	"time"
)

// User 用户实体类
type User struct {
	ID           int64     `json:"id" gorm:"primaryKey;autoIncrement;comment:id"`
	UserAccount  string    `json:"userAccount" gorm:"column:user_account;type:varchar(256);not null;uniqueIndex;comment:账号"`
	UserPassword string    `json:"userPassword" gorm:"column:user_password;type:varchar(512);not null;comment:密码"`
	UserName     string    `json:"userName" gorm:"column:user_name;type:varchar(256);comment:用户昵称"`
	UserAvatar   string    `json:"userAvatar" gorm:"column:user_avatar;type:varchar(1024);comment:用户头像"`
	UserProfile  string    `json:"userProfile" gorm:"column:user_profile;type:varchar(512);comment:用户简介"`
	UserRole     string    `json:"userRole" gorm:"column:user_role;type:varchar(256);default:user;not null;comment:用户角色：user/admin"`
	EditTime     time.Time `json:"editTime" gorm:"column:edit_time;comment:编辑时间"`
	CreateTime   time.Time `json:"createTime" gorm:"column:create_time;autoCreateTime;comment:创建时间"`
	UpdateTime   time.Time `json:"updateTime" gorm:"column:update_time;autoUpdateTime;comment:更新时间"`
	IsDelete     int       `json:"isDelete" gorm:"column:is_delete;type:tinyint;default:0;not null;comment:是否删除(0-未删除，1-已删除)"`
}

// TableName 指定表名
func (User) TableName() string {
	return "user"
}
