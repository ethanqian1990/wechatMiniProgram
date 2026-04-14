// 商品 API
import { get } from '@/utils/request';

export const getProducts = (params) => get('/products', params);
export const getProduct = (id) => get(`/products/${id}`);
