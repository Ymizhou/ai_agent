package controller

import (
	"aicode/common"
	"aicode/constant"
	"aicode/exception"
	"aicode/model/dto/user"
	_ "aicode/model/entity"
	_ "aicode/model/vo"
	"aicode/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// UserController 用户控制层
type UserController struct {
	userService service.UserService
}

// NewUserController 创建用户控制器
func NewUserController(userService service.UserService) *UserController {
	return &UserController{
		userService: userService,
	}
}

// UserRegister 用户注册
// @Summary 用户注册
// @Description 用户注册接口
// @Tags 用户模块
// @Accept json
// @Produce json
// @Param request body user.UserRegisterRequest true "用户注册请求"
// @Success 200 {object} common.BaseResponse[int64]
// @Router /user/register [post]
func (ctrl *UserController) UserRegister(c *gin.Context) {
	var req user.UserRegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, common.Error(exception.ParamsError))
		return
	}

	result, err := ctrl.userService.UserRegister(req.UserAccount, req.UserPassword, req.CheckPassword)
	if err != nil {
		if bizErr, ok := err.(*exception.BusinessError); ok {
			c.JSON(http.StatusOK, common.ErrorWithCode(bizErr.Code(), bizErr.Message()))
			return
		}
		c.JSON(http.StatusOK, common.Error(exception.SystemError))
		return
	}

	c.JSON(http.StatusOK, common.Success(result))
}

// UserLogin 用户登录
// @Summary 用户登录
// @Description 用户登录接口
// @Tags 用户模块
// @Accept json
// @Produce json
// @Param request body user.UserLoginRequest true "用户登录请求"
// @Success 200 {object} common.BaseResponse[vo.LoginUserVO]
// @Router /user/login [post]
func (ctrl *UserController) UserLogin(c *gin.Context) {
	var req user.UserLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, common.Error(exception.ParamsError))
		return
	}

	loginUserVO, err := ctrl.userService.UserLogin(req.UserAccount, req.UserPassword, c)
	if err != nil {
		if bizErr, ok := err.(*exception.BusinessError); ok {
			c.JSON(http.StatusOK, common.ErrorWithCode(bizErr.Code(), bizErr.Message()))
			return
		}
		c.JSON(http.StatusOK, common.Error(exception.SystemError))
		return
	}

	c.JSON(http.StatusOK, common.Success(loginUserVO))
}

// GetLoginUser 获取当前登录用户
// @Summary 获取当前登录用户
// @Description 获取当前登录用户信息
// @Tags 用户模块
// @Accept json
// @Produce json
// @Success 200 {object} common.BaseResponse[vo.LoginUserVO]
// @Router /user/get/login [get]
func (ctrl *UserController) GetLoginUser(c *gin.Context) {
	loginUser, err := ctrl.userService.GetLoginUser(c)
	if err != nil {
		if bizErr, ok := err.(*exception.BusinessError); ok {
			c.JSON(http.StatusOK, common.ErrorWithCode(bizErr.Code(), bizErr.Message()))
			return
		}
		c.JSON(http.StatusOK, common.Error(exception.SystemError))
		return
	}

	c.JSON(http.StatusOK, common.Success(ctrl.userService.GetLoginUserVO(loginUser)))
}

// UserLogout 用户注销
// @Summary 用户注销
// @Description 用户注销登录
// @Tags 用户模块
// @Accept json
// @Produce json
// @Success 200 {object} common.BaseResponse[bool]
// @Router /user/logout [post]
func (ctrl *UserController) UserLogout(c *gin.Context) {
	result, err := ctrl.userService.UserLogout(c)
	if err != nil {
		if bizErr, ok := err.(*exception.BusinessError); ok {
			c.JSON(http.StatusOK, common.ErrorWithCode(bizErr.Code(), bizErr.Message()))
			return
		}
		c.JSON(http.StatusOK, common.Error(exception.SystemError))
		return
	}

	c.JSON(http.StatusOK, common.Success(result))
}

// AddUser 创建用户（仅管理员）
// @Summary 创建用户
// @Description 管理员创建用户接口
// @Tags 用户模块
// @Accept json
// @Produce json
// @Param request body user.UserAddRequest true "用户创建请求"
// @Success 200 {object} common.BaseResponse[int64]
// @Router /user/add [post]
func (ctrl *UserController) AddUser(c *gin.Context) {
	// TODO: 添加管理员权限检查中间件
	var req user.UserAddRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, common.Error(exception.ParamsError))
		return
	}

	result, err := ctrl.userService.AddUser(&req)
	if err != nil {
		if bizErr, ok := err.(*exception.BusinessError); ok {
			c.JSON(http.StatusOK, common.ErrorWithCode(bizErr.Code(), bizErr.Message()))
			return
		}
		c.JSON(http.StatusOK, common.Error(exception.SystemError))
		return
	}

	c.JSON(http.StatusOK, common.Success(result))
}

// GetUserById 根据 id 获取用户（仅管理员）
// @Summary 根据ID获取用户
// @Description 管理员根据ID获取用户详细信息
// @Tags 用户模块
// @Accept json
// @Produce json
// @Param id query int64 true "用户ID"
// @Success 200 {object} common.BaseResponse[entity.User]
// @Router /user/get [get]
func (ctrl *UserController) GetUserById(c *gin.Context) {
	// TODO: 添加管理员权限检查中间件
	idStr := c.Query("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 {
		c.JSON(http.StatusOK, common.Error(exception.ParamsError))
		return
	}

	user, err := ctrl.userService.GetById(id)
	if err != nil {
		if bizErr, ok := err.(*exception.BusinessError); ok {
			c.JSON(http.StatusOK, common.ErrorWithCode(bizErr.Code(), bizErr.Message()))
			return
		}
		c.JSON(http.StatusOK, common.Error(exception.SystemError))
		return
	}

	c.JSON(http.StatusOK, common.Success(user))
}

// GetUserVOById 根据 id 获取包装类
// @Summary 根据ID获取用户VO
// @Description 根据ID获取脱敏后的用户信息
// @Tags 用户模块
// @Accept json
// @Produce json
// @Param id query int64 true "用户ID"
// @Success 200 {object} common.BaseResponse[vo.UserVO]
// @Router /user/get/vo [get]
func (ctrl *UserController) GetUserVOById(c *gin.Context) {
	idStr := c.Query("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 {
		c.JSON(http.StatusOK, common.Error(exception.ParamsError))
		return
	}

	user, err := ctrl.userService.GetById(id)
	if err != nil {
		if bizErr, ok := err.(*exception.BusinessError); ok {
			c.JSON(http.StatusOK, common.ErrorWithCode(bizErr.Code(), bizErr.Message()))
			return
		}
		c.JSON(http.StatusOK, common.Error(exception.SystemError))
		return
	}

	c.JSON(http.StatusOK, common.Success(ctrl.userService.GetUserVO(user)))
}

// DeleteUser 删除用户（仅管理员）
// @Summary 删除用户
// @Description 管理员删除用户接口
// @Tags 用户模块
// @Accept json
// @Produce json
// @Param request body common.DeleteRequest true "删除请求"
// @Success 200 {object} common.BaseResponse[bool]
// @Router /user/delete [post]
func (ctrl *UserController) DeleteUser(c *gin.Context) {
	// TODO: 添加管理员权限检查中间件
	var req common.DeleteRequest
	if err := c.ShouldBindJSON(&req); err != nil || req.ID <= 0 {
		c.JSON(http.StatusOK, common.Error(exception.ParamsError))
		return
	}

	result, err := ctrl.userService.DeleteById(req.ID)
	if err != nil {
		if bizErr, ok := err.(*exception.BusinessError); ok {
			c.JSON(http.StatusOK, common.ErrorWithCode(bizErr.Code(), bizErr.Message()))
			return
		}
		c.JSON(http.StatusOK, common.Error(exception.SystemError))
		return
	}

	c.JSON(http.StatusOK, common.Success(result))
}

// UpdateUser 更新用户（仅管理员）
// @Summary 更新用户
// @Description 管理员更新用户信息
// @Tags 用户模块
// @Accept json
// @Produce json
// @Param request body user.UserUpdateRequest true "用户更新请求"
// @Success 200 {object} common.BaseResponse[bool]
// @Router /user/update [post]
func (ctrl *UserController) UpdateUser(c *gin.Context) {
	// TODO: 添加管理员权限检查中间件
	var req user.UserUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil || req.ID == 0 {
		c.JSON(http.StatusOK, common.Error(exception.ParamsError))
		return
	}

	result, err := ctrl.userService.UpdateById(&req)
	if err != nil {
		if bizErr, ok := err.(*exception.BusinessError); ok {
			c.JSON(http.StatusOK, common.ErrorWithCode(bizErr.Code(), bizErr.Message()))
			return
		}
		c.JSON(http.StatusOK, common.Error(exception.SystemError))
		return
	}

	c.JSON(http.StatusOK, common.Success(result))
}

// ListUserVOByPage 分页获取用户封装列表（仅管理员）
// @Summary 分页获取用户列表
// @Description 管理员分页查询用户列表
// @Tags 用户模块
// @Accept json
// @Produce json
// @Param request body user.UserQueryRequest true "用户查询请求"
// @Success 200 {object} common.BaseResponse[common.MapResponse]
// @Router /user/list/page/vo [post]
func (ctrl *UserController) ListUserVOByPage(c *gin.Context) {
	// TODO: 添加管理员权限检查中间件
	var req user.UserQueryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, common.Error(exception.ParamsError))
		return
	}

	userVOList, total, err := ctrl.userService.ListUserVOByPage(&req)
	if err != nil {
		if bizErr, ok := err.(*exception.BusinessError); ok {
			c.JSON(http.StatusOK, common.ErrorWithCode(bizErr.Code(), bizErr.Message()))
			return
		}
		c.JSON(http.StatusOK, common.Error(exception.SystemError))
		return
	}

	// 构造分页响应
	pageResponse := map[string]interface{}{
		"records":  userVOList,
		"total":    total,
		"pageNum":  req.PageNum,
		"pageSize": req.PageSize,
	}

	c.JSON(http.StatusOK, common.Success(pageResponse))
}

// RegisterRoutes 注册路由
func (ctrl *UserController) RegisterRoutes(r *gin.RouterGroup) {
	{
		r.POST("/register", ctrl.UserRegister)
		r.POST("/login", ctrl.UserLogin)
		r.GET("/get/login", ctrl.GetLoginUser)
		r.POST("/logout", ctrl.UserLogout)

		// 管理员接口（后续需要添加权限验证中间件）
		r.POST("/add", ctrl.AddUser)
		r.GET("/get", ctrl.GetUserById)
		r.GET("/get/vo", ctrl.GetUserVOById)
		r.POST("/delete", ctrl.DeleteUser)
		r.POST("/update", ctrl.UpdateUser)
		r.POST("/list/page/vo", ctrl.ListUserVOByPage)
	}
}

// CheckAdminAuth 检查管理员权限的中间件（待实现）
func CheckAdminAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从上下文获取登录用户
		userObj, exists := c.Get(constant.UserLoginState)
		if !exists {
			c.JSON(http.StatusOK, common.Error(exception.NotLoginError))
			c.Abort()
			return
		}

		// 类型断言获取用户实体
		// loginUser, ok := userObj.(*entity.User)
		// if !ok || loginUser == nil {
		// 	c.JSON(http.StatusOK, common.Error(exception.NotLoginError))
		// 	c.Abort()
		// 	return
		// }

		// // 检查用户角色
		// if loginUser.UserRole != constant.AdminRole {
		// 	c.JSON(http.StatusOK, common.Error(exception.NoAuthError))
		// 	c.Abort()
		// 	return
		// }

		// TODO: 完善权限检查逻辑
		_ = userObj
		c.Next()
	}
}
