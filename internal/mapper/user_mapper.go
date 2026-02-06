package mapper

import (
	"gorm.io/gorm"

	"aicode/internal/model/entity"
)

// UserMapper 用户数据访问层
type UserMapper struct {
	DB *gorm.DB
}

// NewUserMapper 创建用户Mapper
func NewUserMapper(db *gorm.DB) *UserMapper {
	return &UserMapper{DB: db}
}

// Save 保存用户
func (m *UserMapper) Save(user *entity.User) error {
	return m.DB.Create(user).Error
}

// GetById 根据ID查询用户
func (m *UserMapper) GetById(id int64) (*entity.User, error) {
	var user entity.User
	err := m.DB.Where("id = ? AND is_delete = 0", id).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetByAccount 根据账号查询用户
func (m *UserMapper) GetByAccount(account string) (*entity.User, error) {
	var user entity.User
	err := m.DB.Where("user_account = ? AND is_delete = 0", account).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetByAccountAndPassword 根据账号和密码查询用户
func (m *UserMapper) GetByAccountAndPassword(account, password string) (*entity.User, error) {
	var user entity.User
	err := m.DB.Where("user_account = ? AND user_password = ? AND is_delete = 0", account, password).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// CountByAccount 根据账号统计数量
func (m *UserMapper) CountByAccount(account string) (int64, error) {
	var count int64
	err := m.DB.Model(&entity.User{}).Where("user_account = ? AND is_delete = 0", account).Count(&count).Error
	return count, err
}

// UpdateById 根据ID更新用户
func (m *UserMapper) UpdateById(user *entity.User) error {
	return m.DB.Model(&entity.User{}).Where("id = ?", user.ID).Updates(user).Error
}

// DeleteById 根据ID删除用户（逻辑删除）
func (m *UserMapper) DeleteById(id int64) error {
	return m.DB.Model(&entity.User{}).Where("id = ?", id).Update("is_delete", 1).Error
}

// Page 分页查询用户
func (m *UserMapper) Page(offset, limit int, query *gorm.DB) ([]entity.User, int64, error) {
	var users []entity.User
	var total int64

	// 统计总数
	if err := query.Model(&entity.User{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	if err := query.Offset(offset).Limit(limit).Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, total, nil
}
