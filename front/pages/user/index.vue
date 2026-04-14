<template>
  <view class="page">
    <!-- 用户头部 -->
    <view class="user-header">
      <button class="avatar-btn" open-type="getUserProfile" @click="onLogin">
        <image v-if="userStore.userInfo?.avatar" :src="userStore.userInfo.avatar" class="avatar" />
        <view v-else class="avatar-placeholder">👤</view>
      </button>
      <text class="user-name">{{ userStore.userInfo?.nickname || '微信用户' }}</text>
    </view>
    
    <!-- 菜单列表 -->
    <view class="menu-list">
      <view class="menu-item" @click="goOrders('ALL')">
        <text>我的订单</text>
        <text class="arrow">›</text>
      </view>
      
      <view class="menu-row">
        <view class="menu-tab" v-for="status in orderTabs" :key="status.value" @click="goOrders(status.value)">
          <text class="tab-icon">{{ status.icon }}</text>
          <text class="tab-text">{{ status.label }}</text>
        </view>
      </view>
      
      <view class="menu-item" @click="goAddress">
        <text>收货地址</text>
        <text class="arrow">›</text>
      </view>
    </view>
    
    <!-- 退出登录 -->
    <view v-if="userStore.isLoggedIn" class="logout" @click="onLogout">
      <text>退出登录</text>
    </view>
  </view>
</template>

<script setup>
import { useUserStore } from '@/store/user';

const userStore = useUserStore();

const orderTabs = [
  { label: '待付款', value: 'PENDING_PAY', icon: '⏱' },
  { label: '待发货', value: 'PENDING_SHIP', icon: '📦' },
  { label: '待收货', value: 'SHIPPED', icon: '🚚' },
  { label: '已完成', value: 'FINISHED', icon: '✅' }
];

async function onLogin() {
  if (userStore.isLoggedIn) return;
  
  try {
    const res = await new Promise((resolve, reject) => {
      uni.getUserProfile({
        desc: '用于完善用户资料',
        success: resolve,
        fail: reject
      });
    });
    
    // 获取 code 登录
    const loginRes = await new Promise((resolve, reject) => {
      uni.login({
        provider: 'weixin',
        success: resolve,
        fail: reject
      });
    });
    
    await userStore.login(loginRes.code);
    
    // 更新用户信息
    if (res.userInfo) {
      userStore.userInfo = {
        ...userStore.userInfo,
        nickname: res.userInfo.nickName,
        avatar: res.userInfo.avatarUrl
      };
    }
    
    uni.showToast({ title: '登录成功', icon: 'success' });
  } catch (e) {
    console.error('登录失败', e);
    uni.showToast({ title: '登录失败', icon: 'none' });
  }
}

function onLogout() {
  uni.showModal({
    title: '提示',
    content: '确定退出登录？',
    success: (res) => {
      if (res.confirm) {
        userStore.logout();
        uni.showToast({ title: '已退出', icon: 'success' });
      }
    }
  });
}

function goOrders(status) {
  if (!userStore.isLoggedIn) {
    uni.showToast({ title: '请先登录', icon: 'none' });
    return;
  }
  uni.navigateTo({ url: `/pages/orders/index?status=${status}` });
}

function goAddress() {
  if (!userStore.isLoggedIn) {
    uni.showToast({ title: '请先登录', icon: 'none' });
    return;
  }
  uni.navigateTo({ url: '/pages/address/index' });
}
</script>

<style lang="scss" scoped>
.page {
  min-height: 100vh;
  background: #F5F5F5;
}

.user-header {
  background: linear-gradient(135deg, #E4393C, #FF6B6B);
  padding: 60rpx 40rpx;
  color: #fff;
  
  .avatar-btn {
    width: 120rpx;
    height: 120rpx;
    border-radius: 50%;
    padding: 0;
    border: none;
    background: transparent;
    
    .avatar {
      width: 120rpx;
      height: 120rpx;
      border-radius: 50%;
    }
    
    .avatar-placeholder {
      width: 120rpx;
      height: 120rpx;
      border-radius: 50%;
      background: rgba(255,255,255,0.3);
      display: flex;
      align-items: center;
      justify-content: center;
      font-size: 60rpx;
    }
  }
  
  .user-name {
    display: block;
    margin-top: 20rpx;
    font-size: 32rpx;
  }
}

.menu-list {
  margin: 20rpx 0;
  background: #fff;
}

.menu-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 30rpx 40rpx;
  border-bottom: 1px solid #F5F5F5;
  font-size: 28rpx;
  
  .arrow {
    color: #999;
    font-size: 32rpx;
  }
}

.menu-row {
  display: flex;
  justify-content: space-around;
  padding: 30rpx 0;
  
  .menu-tab {
    text-align: center;
    
    .tab-icon {
      display: block;
      font-size: 40rpx;
      margin-bottom: 8rpx;
    }
    
    .tab-text {
      font-size: 24rpx;
      color: #666;
    }
  }
}

.logout {
  text-align: center;
  padding: 30rpx;
  background: #fff;
  color: #666;
  font-size: 28rpx;
  margin-top: 20rpx;
}
</style>
