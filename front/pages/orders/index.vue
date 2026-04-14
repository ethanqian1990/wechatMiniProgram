<template>
  <view class="page">
    <!-- 状态 Tab -->
    <view class="order-tabs">
      <view 
        v-for="tab in tabs" 
        :key="tab.value"
        class="order-tab"
        :class="{ active: currentStatus === tab.value }"
        @click="onTabChange(tab.value)"
      >
        {{ tab.label }}
      </view>
    </view>
    
    <!-- 订单列表 -->
    <view v-if="loading" class="loading">加载中...</view>
    <view v-else-if="!orders.length" class="empty">暂无订单</view>
    <view v-else class="order-list">
      <view v-for="order in orders" :key="order.id" class="order-card" @click="goDetail(order.id)">
        <view class="order-header">
          <text class="order-no">{{ order.order_no }}</text>
          <text class="order-status" :class="order.status">{{ statusText(order.status) }}</text>
        </view>
        
        <view class="order-goods">
          <image 
            v-for="item in order.items?.slice(0, 3)" 
            :key="item.id"
            :src="item.product_image" 
            class="goods-thumb" 
          />
          <text v-if="order.items?.length > 3" class="more-goods">+{{ order.items.length - 3 }}</text>
        </view>
        
        <view class="order-footer">
          <text class="order-amount">¥{{ order.final_amount }}</text>
          <view class="order-actions">
            <button v-if="order.status === 'PENDING_PAY'" class="btn-pay" @click.stop="goPay(order)">支付</button>
            <button v-if="order.status === 'SHIPPED'" class="btn-confirm" @click.stop="confirmReceive(order.id)">确认收货</button>
            <button v-if="order.status === 'PENDING_PAY'" class="btn-cancel" @click.stop="cancelOrder(order.id)">取消</button>
          </view>
        </view>
      </view>
    </view>
  </view>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import { useOrderStore } from '@/store/order';

const orderStore = useOrderStore();

const tabs = [
  { label: '全部', value: 'ALL' },
  { label: '待付款', value: 'PENDING_PAY' },
  { label: '待发货', value: 'PENDING_SHIP' },
  { label: '待收货', value: 'SHIPPED' },
  { label: '已完成', value: 'FINISHED' },
  { label: '已取消', value: 'CANCELED' }
];

const currentStatus = ref('ALL');
const orders = ref([]);
const loading = ref(false);

onMounted(() => {
  const pages = getCurrentPages();
  const options = pages[pages.length - 1].options;
  if (options.status && options.status !== 'ALL') {
    currentStatus.value = options.status;
  }
  loadOrders();
});

async function loadOrders() {
  loading.value = true;
  try {
    const params = currentStatus.value !== 'ALL' ? { status: currentStatus.value } : {};
    orders.value = await orderStore.fetchOrders(params);
  } finally {
    loading.value = false;
  }
}

function onTabChange(status) {
  currentStatus.value = status;
  loadOrders();
}

function statusText(status) {
  const map = {
    'PENDING_PAY': '待付款',
    'PENDING_SHIP': '待发货',
    'SHIPPED': '待收货',
    'FINISHED': '已完成',
    'CANCELED': '已取消'
  };
  return map[status] || status;
}

function goDetail(id) {
  uni.navigateTo({ url: `/pages/order-detail/index?id=${id}` });
}

async function goPay(order) {
  const { continuePay } = require('@/api/order');
  try {
    const payParams = await continuePay(order.id);
    await new Promise((resolve, reject) => {
      uni.requestPayment({
        provider: 'wxpay',
        ...payParams,
        success: resolve,
        fail: reject
      });
    });
    uni.showLoading({ title: '支付结果确认中...' });
    const { pollPayStatus } = require('@/utils/pay');
    await pollPayStatus(order.id);
    uni.hideLoading();
    uni.showToast({ title: '支付成功', icon: 'success' });
    loadOrders();
  } catch (e) {
    if (!e.errMsg?.includes('cancel')) {
      uni.showToast({ title: '支付失败', icon: 'none' });
    }
  } finally {
    uni.hideLoading();
  }
}

async function confirmReceive(id) {
  try {
    await orderStore.confirm(id);
    uni.showToast({ title: '确认收货成功', icon: 'success' });
  } catch (e) {
    console.error(e);
  }
}

async function cancelOrder(id) {
  uni.showModal({
    title: '提示',
    content: '确定取消该订单？',
    success: async (res) => {
      if (res.confirm) {
        try {
          await orderStore.cancel(id);
          uni.showToast({ title: '已取消', icon: 'success' });
          loadOrders();
        } catch (e) {
          uni.showToast({ title: '取消失败', icon: 'none' });
        }
      }
    }
  });
}
</script>

<style lang="scss" scoped>
.page {
  min-height: 100vh;
  background: #F5F5F5;
}

.order-tabs {
  display: flex;
  background: #fff;
  padding: 20rpx 0;
  position: sticky;
  top: 0;
  z-index: 10;
  
  .order-tab {
    flex: 1;
    text-align: center;
    font-size: 26rpx;
    color: #666;
    padding: 16rpx 0;
    
    &.active {
      color: #E4393C;
      font-weight: 600;
    }
  }
}

.order-list {
  padding: 20rpx;
}

.order-card {
  background: #fff;
  border-radius: 16rpx;
  padding: 24rpx;
  margin-bottom: 20rpx;
  
  .order-header {
    display: flex;
    justify-content: space-between;
    margin-bottom: 20rpx;
    
    .order-no {
      font-size: 24rpx;
      color: #999;
    }
    
    .order-status {
      font-size: 26rpx;
      
      &.PENDING_PAY { color: #E4393C; }
      &.PENDING_SHIP { color: #11998E; }
      &.SHIPPED { color: #E4393C; }
      &.FINISHED { color: #11998E; }
      &.CANCELED { color: #999; }
    }
  }
  
  .order-goods {
    display: flex;
    gap: 12rpx;
    margin-bottom: 20rpx;
    
    .goods-thumb {
      width: 120rpx;
      height: 120rpx;
      border-radius: 8rpx;
      background: #EEE;
    }
    
    .more-goods {
      font-size: 24rpx;
      color: #999;
      line-height: 120rpx;
    }
  }
  
  .order-footer {
    display: flex;
    justify-content: space-between;
    align-items: center;
    
    .order-amount {
      font-size: 32rpx;
      font-weight: 700;
      color: #E4393C;
    }
    
    .order-actions {
      display: flex;
      gap: 16rpx;
      
      button {
        padding: 12rpx 24rpx;
        border-radius: 8rpx;
        font-size: 24rpx;
      }
      
      .btn-pay {
        background: #E4393C;
        color: #fff;
      }
      
      .btn-confirm {
        background: #11998E;
        color: #fff;
      }
      
      .btn-cancel {
        background: #F5F5F5;
        color: #666;
      }
    }
  }
}

.loading, .empty {
  text-align: center;
  padding: 200rpx;
  color: #999;
}
</style>
