// 用户 API
import { post, get } from '@/utils/request';

export const wxLogin = (code) => post('/user/wx_login', { code });
export const getUserInfo = () => get('/user/info');
export const updateProfile = (data) => post('/user/profile', data);
