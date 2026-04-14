package service

import (
	github.com/ethanqian1990/wechat-mall-backend/internal/model
	github.com/ethanqian1990/wechat-mall-backend/internal/repository
	github.com/google/uuid
	errors errors
)

type AddressService struct {
	addressRepo *repo.AddressRepository
}

func NewAddressService() *AddressService {
	return &AddressService{
		addressRepo: repo.NewAddressRepository(),
	}
}

// GetAddresses 获取用户地址列表
func (s *AddressService) GetAddresses(userID string) ([]model.Address, error) {
	return s.addressRepo.FindByUserID(userID)
}

// GetAddressByID 获取地址详情
func (s *AddressService) GetAddressByID(id, userID string) (*model.Address, error) {
	return s.addressRepo.FindByIDAndUserID(id, userID)
}

// CreateAddress 创建地址
func (s *AddressService) CreateAddress(userID, receiver, phone, province, city, district, detail string, isDefault bool) (*model.Address, error) {
	if isDefault {
		s.addressRepo.ClearDefault(userID)
	}

	address := &model.Address{
		ID:        'addr_' + uuid.New().String()[:8],
		UserID:    userID,
		Receiver:  receiver,
		Phone:     phone,
		Province:  province,
		City:      city,
		District:  district,
		Detail:    detail,
		IsDefault: 0,
	}
	if isDefault {
		address.IsDefault = 1
	}

	err := s.addressRepo.Create(address)
	return address, err
}

// UpdateAddress 更新地址
func (s *AddressService) UpdateAddress(id, userID, receiver, phone, province, city, district, detail string, isDefault bool) error {
	if isDefault {
		s.addressRepo.ClearDefault(userID)
	}

	updates := map[string]interface{}{
		'receiver': receiver,
		'phone':    phone,
		'province': province,
		'city':     city,
		'district': district,
		'detail':   detail,
	}
	if isDefault {
		updates['is_default'] = 1
	}

	return s.addressRepo.Update(id, userID, updates)
}

// DeleteAddress 删除地址
func (s *AddressService) DeleteAddress(id, userID string) error {
	return s.addressRepo.Delete(id, userID)
}

// SetDefault 设为默认地址
func (s *AddressService) SetDefault(id, userID string) error {
	s.addressRepo.ClearDefault(userID)
	return s.addressRepo.SetDefault(id, userID)
}

// GetDefaultAddress 获取默认地址
func (s *AddressService) GetDefaultAddress(userID string) (*model.Address, error) {
	return s.addressRepo.FindDefault(userID)
}