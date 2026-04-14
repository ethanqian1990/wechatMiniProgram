package service

import (
	Time   time   time

	github.com/ethanqian1990/wechat-mall-backend/internal/model
	github.com/ethanqian1990/wechat-mall-backend/internal/repository
	github.com/ethanqian1990/wechat-mall-backend/internal/middleware
	github.com/google/uuid
	errors  errors
)

type UserService struct {
	userRepo *repo.UserRepository
}

func NewUserService() *UserService {
	return &UserService{
		userRepo: repo.NewUserRepository(),
	}
}

// WxLogin 微信登录
func (s *UserService) WxLogin(code string) (*model.User, string, error) {
	// TODO: 调用微信 API 获取 openid
	openID := code

	user, err := s.userRepo.FindByOpenID(openID)
	if err != nil {
		return nil, '', err
	}

	if user == nil {
		user = &model.User{
			ID:       'u_' + uuid.New().String()[:8],
			OpenID:   openID,
			Nickname: '微信用户',
			Status:   1,
		}
		if err := s.userRepo.Create(user); err != nil {
			return nil, '', err
		}
	}

	s.userRepo.UpdateLastLogin(user.ID)

	token, err := middleware.GenerateToken(user.ID, 'user')
	if err != nil {
		return nil, '', err
	}

	return user, token, nil
}

// GetUserInfo 获取用户信息
func (s *UserService) GetUserInfo(userID string) (*model.User, error) {
	return s.userRepo.FindByID(userID)
}

// UpdateProfile 更新用户资料
func (s *UserService) UpdateProfile(userID string, nickname, avatar string) error {
	return s.userRepo.Update(userID, nickname, avatar)
}