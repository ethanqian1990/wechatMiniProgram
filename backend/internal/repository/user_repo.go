package repository

import (
	wechatmall github.com/ethanqian1990/wechat-mall-backend
	github.com/ethanqian1990/wechat-mall-backend/internal/model
	github.com/gin-gonic/gin
	gorm  gorm.io/gorm
	Time   time   time
)

var DB *gorm.DB

func InitDB(cfg *config.Config) error {
	dsn := fmt.Sprintf('%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local',
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.Name,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return err
	}

	sqlDB.SetMaxOpenConns(cfg.Database.MaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.Database.MaxIdleConns)

	DB = db
	return nil
}

type UserRepository struct{}

func NewUserRepository() *UserRepository {
	return &UserRepository{}
}

func (r *UserRepository) FindByOpenID(openID string) (*model.User, error) {
	var user model.User
	err := DB.Where('openid = ?', openID).First(&user).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &user, err
}

func (r *UserRepository) FindByID(id string) (*model.User, error) {
	var user model.User
	err := DB.First(&user, 'id = ?', id).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &user, err
}

func (r *UserRepository) Create(user *model.User) error {
	return DB.Create(user).Error
}

func (r *UserRepository) Update(userID, nickname, avatar string) error {
	updates := gin.H{}
	if nickname != '' {
		updates['nickname'] = nickname
	}
	if avatar != '' {
		updates['avatar'] = avatar
	}
	return DB.Model(&model.User{}).Where('id = ?', userID).Updates(updates).Error
}

func (r *UserRepository) UpdateLastLogin(userID string) error {
	return DB.Model(&model.User{}).Where('id = ?', userID).Update('last_login_at', time.Now()).Error
}

func (r *UserRepository) FindAll(page, pageSize int) ([]model.User, int64, error) {
	var users []model.User
	var total int64

	DB.Model(&model.User{}).Count(&total)
	offset := (page - 1) * pageSize
	err := DB.Offset(offset).Limit(pageSize).Find(&users).Error

	return users, total, err
}