// 订单 API
import { get, post, put } from '@/utils/request';

export const getOrders = (params) => get('/orders', params);
export const getOrder = (id) => get(`/orders/${id}`);
export const createOrder = (data) => post('/orders', data);
export const cancelOrder = (id) => put(`/orders/${id}/cancel`, {});
export const confirmOrder = (id) => put(`/orders/${id}/confirm`, {});
export const continuePay = (id) => put(`/orders/${id}/pay`, {});
