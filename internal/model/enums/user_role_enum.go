package enums

// UserRoleEnum 用户角色枚举
type UserRoleEnum struct {
	text  string
	value string
}

// Text 获取文本描述
func (e UserRoleEnum) Text() string {
	return e.text
}

// Value 获取值
func (e UserRoleEnum) Value() string {
	return e.value
}

var (
	// USER 普通用户
	USER = UserRoleEnum{text: "用户", value: "user"}
	// ADMIN 管理员
	ADMIN = UserRoleEnum{text: "管理员", value: "admin"}
)

// GetEnumByValue 根据 value 获取枚举
func GetEnumByValue(value string) *UserRoleEnum {
	if value == "" {
		return nil
	}
	allRoles := []UserRoleEnum{USER, ADMIN}
	for _, role := range allRoles {
		if role.value == value {
			return &role
		}
	}
	return nil
}
