package impl

import (
	"aicode/constant"
	"aicode/internal/exception"
	"aicode/internal/mapper"
	"aicode/internal/model/dto/user"
	"aicode/internal/model/entity"
	"aicode/internal/model/enums"
	"aicode/internal/model/vo"
	"aicode/internal/service"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// UserServiceImpl 用户服务实现
type UserServiceImpl struct {
	userMapper *mapper.UserMapper
}

// NewUserService 创建用户服务实例
func NewUserService(userMapper *mapper.UserMapper) service.UserService {
	return &UserServiceImpl{
		userMapper: userMapper,
	}
}

// UserRegister 用户注册
func (s *UserServiceImpl) UserRegister(userAccount, userPassword, checkPassword string) (int64, error) {
	// 1. 校验参数
	if strings.TrimSpace(userAccount) == "" ||
		strings.TrimSpace(userPassword) == "" ||
		strings.TrimSpace(checkPassword) == "" {
		return 0, exception.NewBusinessErrorWithMessage(exception.ParamsError, "参数为空")
	}
	if len(userAccount) < 4 {
		return 0, exception.NewBusinessErrorWithMessage(exception.ParamsError, "账号长度过短")
	}
	if len(userPassword) < 8 || len(checkPassword) < 8 {
		return 0, exception.NewBusinessErrorWithMessage(exception.ParamsError, "密码长度过短")
	}
	if userPassword != checkPassword {
		return 0, exception.NewBusinessErrorWithMessage(exception.ParamsError, "两次输入的密码不一致")
	}

	// 2. 查询用户是否已存在
	count, err := s.userMapper.CountByAccount(userAccount)
	if err != nil {
		return 0, exception.NewBusinessErrorWithMessage(exception.SystemError, "查询用户失败")
	}
	if count > 0 {
		return 0, exception.NewBusinessErrorWithMessage(exception.ParamsError, "账号重复")
	}

	// 3. 加密密码
	encryptPassword := s.GetEncryptPassword(userPassword)

	// 4. 创建用户，插入数据库
	newUser := &entity.User{
		UserAccount:  userAccount,
		UserPassword: encryptPassword,
		UserName:     userAccount,
		UserRole:     enums.USER.Value(),
	}

	err = s.userMapper.Save(newUser)
	if err != nil {
		return 0, exception.NewBusinessErrorWithMessage(exception.OperationError, "注册失败，数据库错误")
	}

	return newUser.ID, nil
}

// GetLoginUserVO 获取脱敏的已登录用户信息
func (s *UserServiceImpl) GetLoginUserVO(user *entity.User) *vo.LoginUserVO {
	if user == nil {
		return nil
	}
	return &vo.LoginUserVO{
		ID:          user.ID,
		UserAccount: user.UserAccount,
		UserName:    user.UserName,
		UserAvatar:  user.UserAvatar,
		UserProfile: user.UserProfile,
		UserRole:    user.UserRole,
		CreateTime:  user.CreateTime,
		UpdateTime:  user.UpdateTime,
	}
}

// UserLogin 用户登录
func (s *UserServiceImpl) UserLogin(userAccount, userPassword string, c *gin.Context) (*vo.LoginUserVO, error) {
	// 1. 校验参数
	if strings.TrimSpace(userAccount) == "" || strings.TrimSpace(userPassword) == "" {
		return nil, exception.NewBusinessErrorWithMessage(exception.ParamsError, "参数为空")
	}
	if len(userAccount) < 4 {
		return nil, exception.NewBusinessErrorWithMessage(exception.ParamsError, "账号长度过短")
	}
	if len(userPassword) < 8 {
		return nil, exception.NewBusinessErrorWithMessage(exception.ParamsError, "密码长度过短")
	}

	// 2. 加密
	encryptPassword := s.GetEncryptPassword(userPassword)

	// 3. 查询用户是否存在
	loginUser, err := s.userMapper.GetByAccountAndPassword(userAccount, encryptPassword)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, exception.NewBusinessErrorWithMessage(exception.ParamsError, "用户不存在或密码错误")
		}
		return nil, exception.NewBusinessErrorWithMessage(exception.SystemError, "查询用户失败")
	}

	// 4. 记录用户的登录态（使用 session）
	session := gin.H{
		constant.UserLoginState: loginUser,
	}
	c.Set(constant.UserLoginState, loginUser)
	// 可以使用 gin-contrib/sessions 来实现真正的 session 管理
	_ = session

	// 5. 返回脱敏的用户信息
	return s.GetLoginUserVO(loginUser), nil
}

// GetLoginUser 获取当前登录用户
func (s *UserServiceImpl) GetLoginUser(c *gin.Context) (*entity.User, error) {
	// 先判断用户是否登录
	userObj, exists := c.Get(constant.UserLoginState)
	if !exists {
		return nil, exception.NewBusinessErrorFromCode(exception.NotLoginError)
	}

	currentUser, ok := userObj.(*entity.User)
	if !ok || currentUser == nil || currentUser.ID == 0 {
		return nil, exception.NewBusinessErrorFromCode(exception.NotLoginError)
	}

	// 从数据库查询当前用户信息
	userId := currentUser.ID
	currentUser, err := s.userMapper.GetById(userId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, exception.NewBusinessErrorFromCode(exception.NotLoginError)
		}
		return nil, exception.NewBusinessErrorWithMessage(exception.SystemError, "查询用户失败")
	}

	return currentUser, nil
}

// GetUserVO 获取脱敏后的用户信息
func (s *UserServiceImpl) GetUserVO(user *entity.User) *vo.UserVO {
	if user == nil {
		return nil
	}
	return &vo.UserVO{
		ID:          user.ID,
		UserAccount: user.UserAccount,
		UserName:    user.UserName,
		UserAvatar:  user.UserAvatar,
		UserProfile: user.UserProfile,
		UserRole:    user.UserRole,
		CreateTime:  user.CreateTime,
	}
}

// GetUserVOList 获取脱敏后的用户信息列表
func (s *UserServiceImpl) GetUserVOList(users []entity.User) []vo.UserVO {
	if len(users) == 0 {
		return []vo.UserVO{}
	}
	userVOList := make([]vo.UserVO, 0, len(users))
	for _, user := range users {
		userVO := s.GetUserVO(&user)
		if userVO != nil {
			userVOList = append(userVOList, *userVO)
		}
	}
	return userVOList
}

// UserLogout 用户注销
func (s *UserServiceImpl) UserLogout(c *gin.Context) (bool, error) {
	// 先判断用户是否登录
	_, exists := c.Get(constant.UserLoginState)
	if !exists {
		return false, exception.NewBusinessErrorWithMessage(exception.OperationError, "用户未登录")
	}
	// 移除登录态
	c.Set(constant.UserLoginState, nil)
	return true, nil
}

// AddUser 创建用户（管理员）
func (s *UserServiceImpl) AddUser(req *user.UserAddRequest) (int64, error) {
	if req == nil {
		return 0, exception.NewBusinessErrorFromCode(exception.ParamsError)
	}

	// 默认密码
	encryptPassword := s.GetEncryptPassword(constant.DefaultPassword)

	newUser := &entity.User{
		UserAccount:  req.UserAccount,
		UserPassword: encryptPassword,
		UserName:     req.UserName,
		UserAvatar:   req.UserAvatar,
		UserProfile:  req.UserProfile,
		UserRole:     req.UserRole,
		EditTime:     time.Now(),
	}

	err := s.userMapper.Save(newUser)
	if err != nil {
		return 0, exception.NewBusinessErrorFromCode(exception.OperationError)
	}

	return newUser.ID, nil
}

// GetById 根据ID获取用户
func (s *UserServiceImpl) GetById(id int64) (*entity.User, error) {
	if id <= 0 {
		return nil, exception.NewBusinessErrorFromCode(exception.ParamsError)
	}
	user, err := s.userMapper.GetById(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, exception.NewBusinessErrorFromCode(exception.NotFoundError)
		}
		return nil, exception.NewBusinessErrorWithMessage(exception.SystemError, "查询用户失败")
	}
	return user, nil
}

// DeleteById 删除用户
func (s *UserServiceImpl) DeleteById(id int64) (bool, error) {
	if id <= 0 {
		return false, exception.NewBusinessErrorFromCode(exception.ParamsError)
	}
	err := s.userMapper.DeleteById(id)
	if err != nil {
		return false, exception.NewBusinessErrorFromCode(exception.OperationError)
	}
	return true, nil
}

// UpdateById 更新用户
func (s *UserServiceImpl) UpdateById(req *user.UserUpdateRequest) (bool, error) {
	if req == nil || req.ID == 0 {
		return false, exception.NewBusinessErrorFromCode(exception.ParamsError)
	}

	updateUser := &entity.User{
		ID:          req.ID,
		UserName:    req.UserName,
		UserAvatar:  req.UserAvatar,
		UserProfile: req.UserProfile,
		UserRole:    req.UserRole,
	}

	err := s.userMapper.UpdateById(updateUser)
	if err != nil {
		return false, exception.NewBusinessErrorFromCode(exception.OperationError)
	}

	return true, nil
}

// ListUserVOByPage 分页获取用户封装列表
func (s *UserServiceImpl) ListUserVOByPage(req *user.UserQueryRequest) ([]vo.UserVO, int64, error) {
	if req == nil {
		return nil, 0, exception.NewBusinessErrorWithMessage(exception.ParamsError, "请求参数为空")
	}

	// 构建查询条件
	query := s.userMapper.DB.Model(&entity.User{}).Where("is_delete = 0")

	if req.ID != nil {
		query = query.Where("id = ?", *req.ID)
	}
	if req.UserRole != "" {
		query = query.Where("user_role = ?", req.UserRole)
	}
	if req.UserAccount != "" {
		query = query.Where("user_account LIKE ?", "%"+req.UserAccount+"%")
	}
	if req.UserName != "" {
		query = query.Where("user_name LIKE ?", "%"+req.UserName+"%")
	}
	if req.UserProfile != "" {
		query = query.Where("user_profile LIKE ?", "%"+req.UserProfile+"%")
	}

	// 排序
	if req.SortField != "" {
		order := "DESC"
		if req.SortOrder == "ascend" {
			order = "ASC"
		}
		query = query.Order(fmt.Sprintf("%s %s", req.SortField, order))
	}

	// 计算分页
	pageNum := req.PageNum
	pageSize := req.PageSize
	if pageNum <= 0 {
		pageNum = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	offset := (pageNum - 1) * pageSize

	// 分页查询
	users, total, err := s.userMapper.Page(offset, pageSize, query)
	if err != nil {
		return nil, 0, exception.NewBusinessErrorWithMessage(exception.SystemError, "查询用户失败")
	}

	// 数据脱敏
	userVOList := s.GetUserVOList(users)
	return userVOList, total, nil
}

// GetEncryptPassword 加密密码
func (s *UserServiceImpl) GetEncryptPassword(userPassword string) string {
	// 盐值，混淆密码
	saltedPassword := userPassword + constant.PasswordSalt
	hash := md5.Sum([]byte(saltedPassword))
	return hex.EncodeToString(hash[:])
}
