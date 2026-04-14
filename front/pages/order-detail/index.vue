<template>
  <view class="page">
    <!-- 订单信息 -->
    <view class="card">
      <view class="section-title">订单信息</view>
      <view class="info-row">
        <text class="label">订单号</text>
        <text class="value">{{ order.order_no }}</text>
      </view>
      <view class="info-row">
        <text class="label">下单时间</text>
        <text class="value">{{ formatTime(order.created_at) }}</text>
      </view>
      <view class="info-row">
        <text class="label">状态</text>
        <text class="value status" :class="order.status">{{ statusText(order.status) }}</text>
      </view>
    </view>
    
    <!-- 收货信息 -->
    <view class="card">
      <view class="section-title">收货信息</view>
      <view class="info-row">
        <text class="label">收货人</text>
        <text class="value">{{ order.receiver }}</text>
      </view>
      <view class="info-row">
        <text class="label">手机号</text>
        <text class="value">{{ order.phone }}</text>
      </view>
      <view class="info-row">
        <text class="label">收货地址</text>
        <text class="value">{{ fullAddress }}</text>
      </view>
    </view>
    
    <!-- 商品信息 -->
    <view class="card">
      <view class="section-title">商品信息</view>
      <view v-for="item in order.items" :key="item.id" class="goods-item">
        <image :src="item.product_image" class="goods-img" />
        <view class="goods-info">
          <text class="goods-name">{{ item.product_name }}</text>
          <text class="goods-sku">{{ item.sku_name }}</text>
          <view class="goods-price-row">
            <text class="goods-price">¥{{ item.price }}</text>
            <text class="goods-qty">×{{ item.quantity }}</text>
          </view>
        </view>
      </view>
    </view>
    
    <!-- 费用明细 -->
    <view class="card">
      <view class="section-title">费用明细</view>
      <view class="fee-row">
        <text>商品总额</text>
        <text>¥{{ order.total_amount }}</text>
      </view>
      <view class="fee-row">
        <text>运费</text>
        <text>¥{{ order.freight_amount || 0 }}</text>
      </view>
      <view class="fee-row">
        <text>优惠</text>
        <text>-¥{{ order.discount_amount || 0 }}</text>
      </view>
      <view class="fee-total">
        <text>实付金额</text>
        <text class="total-price">¥{{ order.final_amount }}</text>
      </view>
    </view>
    
    <!-- 操作按钮 -->
    <view class="action-bar">
      <view class="action-left">
        <button v-if="order.status === 'PENDING_PAY'" class="btn-cancel" @click="onCancel">取消订单</button>
        <button v-if="order.status === 'SHIPPED'" class="btn-confirm" @click="onConfirm">确认收货</button>
      </view>
      <view class="action-right">
        <button v-if="order.status === 'PENDING_PAY'" class="btn-pay" @click="onPay">立即支付</button>
      </view>
    </view>
  </view>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue';
import { useOrderStore } from '@/store/order';
import { createPay } from '@/api/pay';

const orderStore = useOrderStore();
const order = ref({});

const fullAddress = computed(() => {
  if (!order.value.province) return '';
  return `${order.value.province}${order.value.city}${order.value.district}${order.value.detail}`;
});

const statusText = (status) => {
  const map = {
    'PENDING_PAY': '待付款',
    'PENDING_SHIP': '待发货',
    'SHIPPED': '待收货',
    'FINISHED': '已完成',
    'CANCELED': '已取消'
  };
  return map[status] || status;
};

const formatTime = (time) => {
  if (!time) return '';
  return new Date(time).toLocaleString('zh-CN');
};

onMounted(async () => {
  const pages = getCurrentPages();
  const id = pages[pages.length - 1].options.id;
  if (id) {
    order.value = await orderStore.fetchOrder(id);
  }
});

async function onPay() {
  try {
    const payParams = await createPay(order.value.id);
    await new Promise((resolve, reject) => {
      uni.requestPayment({
        provider: 'wxpay',
        ...payParams,
        success: resolve,
        fail: reject
      });
    });
    uni.showToast({ title: '支付成功', icon: 'success' });
    order.value = await orderStore.fetchOrder(order.value.id);
  } catch (e) {
    if (!e.errMsg?.includes('cancel')) {
      uni.showToast({ title: '支付失败', icon: 'none' });
    }
  }
}

async function onCancel() {
  uni.showModal({
    title: '提示',
    content: '确定取消该订单？',
    success: async (res) => {
      if (res.confirm) {
        try {
          await orderStore.cancel(order.value.id);
          uni.showToast({ title: '已取消', icon: 'success' });
          order.value = await orderStore.fetchOrder(order.value.id);
        } catch (e) {
          uni.showToast({ title: '取消失败', icon: 'none' });
        }
      }
    }
  });
}

async function onConfirm() {
  try {
    await orderStore.confirm(order.value.id);
    uni.showToast({ title: '确认收货成功', icon: 'success' });
    order.value = await orderStore.fetchOrder(order.value.id);
  } catch (e) {
    uni.showToast({ title: '操作失败', icon: 'none' });
  }
}
</script>

<style lang="scss" scoped>
.page {
  min-height: 100vh;
  background: #F5F5F5;
  padding-bottom: 120rpx;
}

.card {
  margin: 24rpx;
  background: #fff;
  border-radius: 16rpx;
  padding: 24rpx;
}

.section-title {
  font-size: 28rpx;
  font-weight: 600;
  margin-bottom: 20rpx;
}

.info-row {
  display: flex;
  justify-content: space-between;
  padding: 12rpx 0;
  
  .label {
    font-size: 26rpx;
    color: #999;
  }
  
  .value {
    font-size: 26rpx;
    
    &.status {
      &.PENDING_PAY { color: #E4393C; }
      &.PENDING_SHIP { color: #11998E; }
      &.SHIPPED { color: #E4393C; }
      &.FINISHED { color: #11998E; }
      &.CANCELED { color: #999; }
    }
  }
}

.goods-item {
  display: flex;
  padding: 16rpx 0;
  border-bottom: 1px solid #F5F5F5;
  
  &:last-child {
    border-bottom: none;
  }
  
  .goods-img {
    width: 140rpx;
    height: 140rpx;
    border-radius: 8rpx;
    background: #EEE;
    flex-shrink: 0;
  }
  
  .goods-info {
    flex: 1;
    margin-left: 16rpx;
    
    .goods-name {
      font-size: 26rpx;
      display: -webkit-box;
      -webkit-line-clamp: 2;
      -webkit-box-orient: vertical;
      overflow: hidden;
    }
    
    .goods-sku {
      font-size: 22rpx;
      color: #999;
      display: block;
      margin-top: 8rpx;
    }
    
    .goods-price-row {
      display: flex;
      justify-content: space-between;
      margin-top: 12rpx;
      
      .goods-price {
        font-size: 26rpx;
        color: #E4393C;
        font-weight: 600;
      }
      
      .goods-qty {
        font-size: 24rpx;
        color: #999;
      }
    }
  }
}

.fee-row {
  display: flex;
  justify-content: space-between;
  font-size: 26rpx;
  color: #666;
  padding: 10rpx 0;
}

.fee-total {
  display: flex;
  justify-content: space-between;
  font-size: 28rpx;
  font-weight: 600;
  padding-top: 16rpx;
  border-top: 1px solid #EEE;
  margin-top: 12rpx;
  
  .total-price {
    font-size: 36rpx;
    color: #E4393C;
  }
}

.action-bar {
  position: fixed;
  bottom: 0;
  left: 0;
  right: 0;
  height: 100rpx;
  background: #fff;
  border-top: 1px solid #EEE;
  display: flex;
  align-items: center;
  padding: 0 24rpx;
  
  .action-left {
    flex: 1;
    display: flex;
    gap: 16rpx;
    
    button {
      padding: 16rpx 32rpx;
      border-radius: 999rpx;
      font-size: 26rpx;
    }
    
    .btn-cancel {
      background: #F5F5F5;
      color: #666;
    }
    
    .btn-confirm {
      background: #11998E;
      color: #fff;
    }
  }
  
  .action-right {
    .btn-pay {
      background: #E4393C;
      color: #fff;
      padding: 16rpx 40rpx;
      border-radius: 999rpx;
      font-size: 28rpx;
    }
  }
}
</style>
