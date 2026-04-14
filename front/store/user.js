// 用户 Store
import { defineStore } from 'pinia';
import { getUserInfo, wxLogin } from '@/api/user';

export const useUserStore = defineStore('user', {
  state: () => ({
    token: uni.getStorageSync('token') || '',
    userInfo: uni.getStorageSync('userInfo') || null,
    isLoggedIn: !!uni.getStorageSync('token')
  }),
  
  actions: {
    async login(code) {
      try {
        const data = await wxLogin(code);
        this.token = data.token;
        this.userInfo = data.user;
        this.isLoggedIn = true;
        uni.setStorageSync('token', data.token);
        uni.setStorageSync('userInfo', data.user);
        return data;
      } catch (e) {
        console.error('登录失败', e);
        throw e;
      }
    },
    
    async fetchUserInfo() {
      if (!this.token) return null;
      try {
        const data = await getUserInfo();
        this.userInfo = data;
        uni.setStorageSync('userInfo', data);
        return data;
      } catch (e) {
        console.error('获取用户信息失败', e);
        return null;
      }
    },
    
    logout() {
      this.token = '';
      this.userInfo = null;
      this.isLoggedIn = false;
      uni.removeStorageSync('token');
      uni.removeStorageSync('userInfo');
    }
  }
});
