package handler

import (
	github.com/gin-gonic/gin
	github.com/ethanqian1990/wechat-mall-backend/pkg/response
	github.com/ethanqian1990/wechat-mall-backend/internal/service
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler() *UserHandler {
	return &UserHandler{
		userService: service.NewUserService(),
	}
}

// WxLogin 微信登录
type WxLoginRequest struct {
	Code string `json: code binding:required`
}

func (h *UserHandler) WxLogin(c *gin.Context) {
	var req WxLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, '参数错误')
		return
	}

	user, token, err := h.userService.WxLogin(req.Code)
	if err != nil {
		response.Error(c, 500, '登录失败', err.Error())
		return
	}

	response.Success(c, gin.H{
		'token': token,
		'user': gin.H{
			'id':       user.ID,
			'nickname': user.Nickname,
			'avatar':   user.Avatar,
		},
	})
}

// GetInfo 获取用户信息
func (h *UserHandler) GetInfo(c *gin.Context) {
	userID := c.GetString('user_id')
	
	user, err := h.userService.GetUserInfo(userID)
	if err != nil {
		response.Error(c, 500, '获取用户信息失败')
		return
	}

	response.Success(c, gin.H{
		'id':       user.ID,
		'nickname': user.Nickname,
		'avatar':   user.Avatar,
		'phone':    user.Phone,
	})
}

// UpdateProfile 更新用户资料
type UpdateProfileRequest struct {
	Nickname string `json: nickname`
	Avatar   string `json: avatar`
}

func (h *UserHandler) UpdateProfile(c *gin.Context) {
	userID := c.GetString('user_id')
	
	var req UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, '参数错误')
		return
	}

	err := h.userService.UpdateProfile(userID, req.Nickname, req.Avatar)
	if err != nil {
		response.Error(c, 500, '更新失败')
		return
	}

	response.Success(c, nil)
}