package handler

import (
	"strings"

	"github.com/ethanqian1990/wechat-mall-backend/internal/config"
	"github.com/ethanqian1990/wechat-mall-backend/internal/middleware"
	"github.com/ethanqian1990/wechat-mall-backend/pkg/response"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
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

	if req.Username != h.cfg.Admin.Username {
		response.Error(c, 401, "用户名或密码错误")
		return
	}

	cfgPwd := h.cfg.Admin.Password
	// 支持 bcrypt hash（$2a/$2b/$2y），否则兼容明文
	if strings.HasPrefix(cfgPwd, "$2") {
		if err := bcrypt.CompareHashAndPassword([]byte(cfgPwd), []byte(req.Password)); err != nil {
			response.Error(c, 401, "用户名或密码错误")
			return
		}
	} else {
		if req.Password != cfgPwd {
			response.Error(c, 401, "用户名或密码错误")
			return
		}
	}

	token, err := middleware.GenerateToken("admin", "admin")
	if err != nil {
		response.Error(c, 500, "生成token失败", err.Error())
		return
	}
	response.Success(c, gin.H{"token": token})
}

