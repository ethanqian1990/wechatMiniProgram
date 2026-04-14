// 支付 API
import { post, get } from '@/utils/request';

export const createPay = (orderId) => post('/pay/create', { order_id: orderId });
export const getPayStatus = (orderId) => get(`/pay/status/${orderId}`);
