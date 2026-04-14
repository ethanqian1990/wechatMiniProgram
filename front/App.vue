<script setup>
import { onLaunch, onShow, onHide } from '@dcloudio/uni-app';
import { useUserStore } from '@/store/user';
import { useRegionStore } from '@/store/region';
import { useCategoryStore } from '@/store/category';

onLaunch(() => {
  console.log('App Launch');
  initApp();
});

onShow(() => {
  console.log('App Show');
});

onHide(() => {
  console.log('App Hide');
});

async function initApp() {
  const userStore = useUserStore();
  const regionStore = useRegionStore();
  const categoryStore = useCategoryStore();
  
  // 初始化基础数据
  await Promise.all([
    regionStore.fetchRegions(),
    categoryStore.fetchCategories()
  ]);
  
  // 如果已登录，获取用户信息
  if (userStore.isLoggedIn) {
    userStore.fetchUserInfo();
  }
}
</script>

<style>
@import '@/uni.scss';

/* 全局样式 */
page {
  background-color: #F5F5F5;
  font-size: 13px;
  color: #333;
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
}

.container {
  min-height: 100vh;
}

/* 通用按钮 */
.btn-primary {
  background-color: #E4393C;
  color: #fff;
  border: none;
  border-radius: 8rpx;
  padding: 20rpx 40rpx;
  font-size: 14px;
}

.btn-default {
  background-color: #fff;
  color: #333;
  border: 1px solid #ddd;
  border-radius: 8rpx;
  padding: 20rpx 40rpx;
  font-size: 14px;
}

/* 通用卡片 */
.card {
  background: #fff;
  border-radius: 12rpx;
  box-shadow: 0 2px 12px rgba(0,0,0,0.06);
  margin: 24rpx;
  padding: 24rpx;
}

/* 价格样式 */
.price {
  color: #E4393C;
  font-weight: 700;
}

.price-original {
  color: #999;
  text-decoration: line-through;
  font-size: 12px;
}

/* 安全区域 */
.safe-area-bottom {
  padding-bottom: constant(safe-area-inset-bottom);
  padding-bottom: env(safe-area-inset-bottom);
}
</style>
