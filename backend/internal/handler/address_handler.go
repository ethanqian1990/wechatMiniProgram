package handler

import (
	github.com/gin-gonic/gin
	github.com/ethanqian1990/wechat-mall-backend/pkg/response
	github.com/ethanqian1990/wechat-mall-backend/internal/service
)

type AddressHandler struct {
	addressService *service.AddressService
}

func NewAddressHandler() *AddressHandler {
	return &AddressHandler{
		addressService: service.NewAddressService(),
	}
}

// GetAddresses 获取地址列表
func (h *AddressHandler) GetAddresses(c *gin.Context) {
	userID := c.GetString('user_id')

	addresses, err := h.addressService.GetAddresses(userID)
	if err != nil {
		response.Error(c, 500, '获取地址列表失败')
		return
	}

	var result []gin.H
	for _, addr := range addresses {
		result = append(result, gin.H{
			'id':        addr.ID,
			'receiver':  addr.Receiver,
			'phone':     addr.Phone,
			'province':  addr.Province,
			'city':      addr.City,
			'district':  addr.District,
			'detail':    addr.Detail,
			'is_default': addr.IsDefault == 1,
		})
	}

	response.Success(c, gin.H{'list': result})
}

// CreateAddress 创建地址
type CreateAddressRequest struct {
	Receiver   string `json: receiver binding:required`
	Phone      string `json: phone binding:required`
	Province   string `json: province binding:required`
	City       string `json: city binding:required`
	District   string `json: district binding:required`
	Detail     string `json: detail binding:required`
	IsDefault  bool   `json: is_default`
}

func (h *AddressHandler) CreateAddress(c *gin.Context) {
	userID := c.GetString('user_id')

	var req CreateAddressRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, '参数错误')
		return
	}

	address, err := h.addressService.CreateAddress(
		userID, req.Receiver, req.Phone, req.Province, req.City, req.District, req.Detail, req.IsDefault,
	)
	if err != nil {
		response.Error(c, 500, '创建失败')
		return
	}

	response.Success(c, gin.H{'id': address.ID})
}

// UpdateAddress 更新地址
func (h *AddressHandler) UpdateAddress(c *gin.Context) {
	userID := c.GetString('user_id')
	addressID := c.Param('id')

	var req CreateAddressRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, '参数错误')
		return
	}

	err := h.addressService.UpdateAddress(
		addressID, userID, req.Receiver, req.Phone, req.Province, req.City, req.District, req.Detail, req.IsDefault,
	)
	if err != nil {
		response.Error(c, 500, '更新失败')
		return
	}

	response.Success(c, nil)
}

// DeleteAddress 删除地址
func (h *AddressHandler) DeleteAddress(c *gin.Context) {
	userID := c.GetString('user_id')
	addressID := c.Param('id')

	err := h.addressService.DeleteAddress(addressID, userID)
	if err != nil {
		response.Error(c, 500, '删除失败')
		return
	}

	response.Success(c, nil)
}

// SetDefault 设为默认
func (h *AddressHandler) SetDefault(c *gin.Context) {
	userID := c.GetString('user_id')
	addressID := c.Param('id')

	err := h.addressService.SetDefault(addressID, userID)
	if err != nil {
		response.Error(c, 500, '设置失败')
		return
	}

	response.Success(c, nil)
}