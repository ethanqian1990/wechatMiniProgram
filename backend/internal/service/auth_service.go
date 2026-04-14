package service

import (
	"fmt"
	"time"

	wechatmall "github.com/ethanqian1990/wechat-mall-backend/internal/model"
	wechatrepo "github.com/ethanqian1990/wechat-mall-backend/internal/repository"
	"github.com/ethanqian1990/wechat-mall-backend/pkg/utils"
	"github.com/google/uuid"
)

type AuthService struct {
	userRepo     *wechatrepo.UserRepository
	wechatUtil   *utils.WechatUtil
	jwtSecret    []byte
}

func NewAuthService(wechatUtil *utils.WechatUtil, jwtSecret string) *AuthService {
	return &AuthService{
		userRepo:   wechatrepo.NewUserRepository(),
		wechatUtil: wechatUtil,
		jwtSecret:  []byte(jwtSecret),
	}
}

// WxLogin 微信登录
func (s *AuthService) WxLogin(code string) (*wechatmall.User, string, error) {
	// 1. 调用微信接口获取 openid
	result, err := s.wechatUtil.Code2Session(code)
	if err != nil {
		return nil, "", fmt.Errorf("微信登录失败: %v", err)
	}

	// 2. 查询或创建用户
	openID := result.OpenID
	user, err := s.userRepo.FindByOpenID(openID)
	if err != nil {
		return nil, "", err
	}

	if user == nil {
		// 创建新用户
		user = &wechatmall.User{
			ID:       "u_" + uuid.New().String()[:8],
			OpenID:   openID,
			Nickname: "微信用户",
			Status:   1,
		}
		if err := s.userRepo.Create(user); err != nil {
			return nil, "", fmt.Errorf("创建用户失败: %v", err)
		}
	}

	// 3. 更新最后登录时间
	s.userRepo.UpdateLastLogin(user.ID)

	// 4. 生成 JWT Token
	token, err := s.generateToken(user.ID)
	if err != nil {
		return nil, "", fmt.Errorf("生成Token失败: %v", err)
	}

	return user, token, nil
}

// generateToken 生成JWT Token
func (s *AuthService) generateToken(userID string) (string, error) {
	// 简化实现，实际应使用 JWT 库
	return fmt.Sprintf("jwt_%s_%d", userID, time.Now().Unix()), nil
}

// GetUserInfo 获取用户信息
func (s *AuthService) GetUserInfo(userID string) (*wechatmall.User, error) {
	return s.userRepo.FindByID(userID)
}

// UpdateProfile 更新用户资料
func (s *AuthService) UpdateProfile(userID string, nickname, avatar string) error {
	return s.userRepo.Update(userID, nickname, avatar)
}

// BindPhone 绑定手机号
func (s *AuthService) BindPhone(userID string, code string) error {
	// 1. 调用微信获取手机号接口
	// phone, errCheck := s.wechatUtil.GetPhoneNumber("", code)
	if false {
		return nil
	}

	// 2. 更新用户手机号
	// 这里需要实现 UserRepository 的 UpdatePhone 方法
	return nil
}
