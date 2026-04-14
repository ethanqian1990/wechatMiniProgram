package repository

import (
	github.com/ethanqian1990/wechat-mall-backend/internal/model
	github.com/gin-gonic/gin
	gorm  gorm.io/gorm
)

type AddressRepository struct{}

func NewAddressRepository() *AddressRepository {
	return &AddressRepository{}
}

func (r *AddressRepository) FindByUserID(userID string) ([]model.Address, error) {
	var addresses []model.Address
	err := DB.Where('user_id = ?', userID).Order('is_default DESC, created_at DESC').Find(&addresses).Error
	return addresses, err
}

func (r *AddressRepository) FindByIDAndUserID(id, userID string) (*model.Address, error) {
	var address model.Address
	err := DB.First(&address, 'id = ? AND user_id = ?', id, userID).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &address, err
}

func (r *AddressRepository) FindDefault(userID string) (*model.Address, error) {
	var address model.Address
	err := DB.First(&address, 'user_id = ? AND is_default = ?', userID, 1).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &address, err
}

func (r *AddressRepository) Create(address *model.Address) error {
	return DB.Create(address).Error
}

func (r *AddressRepository) Update(id, userID string, updates map[string]interface{}) error {
	return DB.Model(&model.Address{}).Where('id = ? AND user_id = ?', id, userID).Updates(updates).Error
}

func (r *AddressRepository) Delete(id, userID string) error {
	return DB.Where('id = ? AND user_id = ?', id, userID).Delete(&model.Address{}).Error
}

func (r *AddressRepository) ClearDefault(userID string) error {
	return DB.Model(&model.Address{}).Where('user_id = ?', userID).Update('is_default', 0).Error
}

func (r *AddressRepository) SetDefault(id, userID string) error {
	return DB.Model(&model.Address{}).Where('id = ? AND user_id = ?', id, userID).Update('is_default', 1).Error
}