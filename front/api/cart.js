// 购物车 API
import { get, post, put, del } from '@/utils/request';

export const getCartList = () => get('/cart');
export const addToCart = (data) => post('/cart', data);
export const updateCartItem = (id, data) => put(`/cart/${id}`, data);
export const removeCartItem = (id) => del(`/cart/${id}`);
export const clearCart = () => del('/cart/clear');
