package handler

import (
	"github.com/ethanqian1990/wechat-mall-backend/internal/config"
	"github.com/ethanqian1990/wechat-mall-backend/internal/middleware"
	"github.com/ethanqian1990/wechat-mall-backend/pkg/response"
	"github.com/gin-gonic/gin"
)

type AdminAuthHandler struct {
	cfg *config.Config
}

func NewAdminAuthHandler(cfg *config.Config) *AdminAuthHandler {
	return &AdminAuthHandler{cfg: cfg}
}

func (h *AdminAuthHandler) Login(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, "参数错误")
		return
	}

	if req.Username != h.cfg.Admin.Username || req.Password != h.cfg.Admin.Password {
		response.Error(c, 401, "用户名或密码错误")
		return
	}

	token, err := middleware.GenerateToken("admin", "admin")
	if err != nil {
		response.Error(c, 500, "生成token失败", err.Error())
		return
	}
	response.Success(c, gin.H{"token": token})
}

