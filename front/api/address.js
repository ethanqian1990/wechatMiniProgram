// 地址 API
import { get, post, put, del } from '@/utils/request';

export const getAddressList = () => get('/addresses');
export const getAddress = (id) => get(`/addresses/${id}`);
export const addAddress = (data) => post('/addresses', data);
export const updateAddress = (id, data) => put(`/addresses/${id}`, data);
export const deleteAddress = (id) => del(`/addresses/${id}`);
export const setDefaultAddress = (id) => put(`/addresses/${id}/default`, {});
