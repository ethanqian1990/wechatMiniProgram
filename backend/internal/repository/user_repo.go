package repository

import (
	github.com/ethanqian1990/wechat-mall-backend/internal/model
	Time   time   time
)

type UserRepository struct{}

func NewUserRepository() *UserRepository {
	return &UserRepository{}
}

func (r *UserRepository) FindByOpenID(openID string) (*model.User, error) {
	// TODO: 实现数据库查询
	return nil, nil
}

func (r *UserRepository) FindByID(id string) (*model.User, error) {
	// TODO: 实现数据库查询
	return nil, nil
}

func (r *UserRepository) Create(user *model.User) error {
	// TODO: 实现数据库创建
	return nil
}

func (r *UserRepository) Update(userID, nickname, avatar string) error {
	// TODO: 实现数据库更新
	return nil
}

func (r *UserRepository) UpdateLastLogin(userID string) error {
	// TODO: 实现最后登录时间更新
	return nil
}

func (r *UserRepository) FindAll(page, pageSize int) ([]model.User, int64, error) {
	// TODO: 实现分页查询
	return nil, 0, nil
}