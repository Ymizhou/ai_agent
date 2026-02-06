package service

import (
	"aicode/internal/model/dto/user"
	"aicode/internal/model/entity"
	"aicode/internal/model/vo"

	"github.com/gin-gonic/gin"
)

// UserService 用户服务接口
type UserService interface {
	// UserRegister 用户注册
	UserRegister(userAccount, userPassword, checkPassword string) (int64, error)

	// UserLogin 用户登录
	UserLogin(userAccount, userPassword string, c *gin.Context) (*vo.LoginUserVO, error)

	// GetLoginUser 获取当前登录用户
	GetLoginUser(c *gin.Context) (*entity.User, error)

	// GetLoginUserVO 获取脱敏的已登录用户信息
	GetLoginUserVO(user *entity.User) *vo.LoginUserVO

	// UserLogout 用户注销
	UserLogout(c *gin.Context) (bool, error)

	// AddUser 创建用户（管理员）
	AddUser(req *user.UserAddRequest) (int64, error)

	// GetById 根据ID获取用户
	GetById(id int64) (*entity.User, error)

	// GetUserVO 获取脱敏后的用户信息
	GetUserVO(user *entity.User) *vo.UserVO

	// GetUserVOList 获取脱敏后的用户信息列表
	GetUserVOList(users []entity.User) []vo.UserVO

	// DeleteById 删除用户
	DeleteById(id int64) (bool, error)

	// UpdateById 更新用户
	UpdateById(req *user.UserUpdateRequest) (bool, error)

	// ListUserVOByPage 分页获取用户封装列表
	ListUserVOByPage(req *user.UserQueryRequest) ([]vo.UserVO, int64, error)

	// GetEncryptPassword 加密密码
	GetEncryptPassword(userPassword string) string
}
