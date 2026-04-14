<template>
  <view class="page">
    <!-- 地址选择 -->
    <view class="address-card" @click="goAddressSelect">
      <view v-if="addressStore.selectedAddress" class="address-info">
        <text class="receiver">{{ addressStore.selectedAddress.receiver }} {{ addressStore.selectedAddress.phone }}</text>
        <text class="detail">{{ fullAddress }}</text>
      </view>
      <view v-else class="address-empty">
        <text>请选择收货地址</text>
        <text class="arrow">›</text>
      </view>
    </view>
    
    <!-- 商品信息 -->
    <view class="card goods-card">
      <view class="section-title">商品信息</view>
      <view v-for="item in orderItems" :key="item.id" class="goods-item">
        <image :src="item.product_image || item.image" class="goods-img" />
        <view class="goods-info">
          <text class="goods-name">{{ item.product_name || item.name }}</text>
          <text class="goods-sku">{{ item.sku_name || item.spec }}</text>
          <view class="goods-price-row">
            <text class="goods-price">¥{{ item.price }}</text>
            <text class="goods-qty">×{{ item.quantity }}</text>
          </view>
        </view>
      </view>
    </view>
    
    <!-- 费用明细 -->
    <view class="card fee-card">
      <view class="section-title">费用明细</view>
      <view class="fee-row">
        <text>商品小计</text>
        <text>¥{{ subtotal }}</text>
      </view>
      <view class="fee-row">
        <text>运费</text>
        <text>¥{{ freight }}</text>
      </view>
      <view class="fee-row">
        <text>优惠</text>
        <text>-¥{{ discount }}</text>
      </view>
      <view class="fee-total">
        <text>应付合计</text>
        <text class="total-price">¥{{ total }}</text>
      </view>
    </view>
    
    <!-- 提交按钮 -->
    <view class="submit-bar">
      <view class="total-info">
        合计 <text class="total-price">¥{{ total }}</text>
      </view>
      <button class="btn-submit" :disabled="submitting" @click="onSubmit">提交订单</button>
    </view>
  </view>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue';
import { useCartStore } from '@/store/cart';
import { useAddressStore } from '@/store/address';
import { useRegionStore } from '@/store/region';
import { useOrderStore } from '@/store/order';
import { createPay } from '@/api/pay';

const cartStore = useCartStore();
const addressStore = useAddressStore();
const regionStore = useRegionStore();
const orderStore = useOrderStore();

const orderItems = ref([]);
const submitting = ref(false);

const subtotal = computed(() => {
  return orderItems.value.reduce((sum, item) => sum + item.price * item.quantity, 0);
});

const freight = ref(0);
const discount = ref(0);
const total = computed(() => subtotal.value - discount.value + freight.value);

const fullAddress = computed(() => {
  const addr = addressStore.selectedAddress;
  if (!addr) return '';
  return `${addr.province}${addr.city}${addr.district}${addr.detail}`;
});

onMounted(async () => {
  const pages = getCurrentPages();
  const options = pages[pages.length - 1].options;
  
  // 加载地址
  await addressStore.fetchAddresses();
  
  // 判断是购物车结算还是立即购买
  if (options.goodsId) {
    // 立即购买
    orderItems.value = [{
      product_id: options.goodsId,
      sku_id: options.skuId,
      quantity: parseInt(options.quantity) || 1,
      price: 0 // 后面从商品详情获取
    }];
  } else {
    // 从购物车来
    orderItems.value = cartStore.checkedItems.map(item => ({
      product_id: item.product_id,
      sku_id: item.sku_id,
      quantity: item.quantity,
      price: item.sku?.price || 0,
      product_name: item.product?.name,
      product_image: item.product?.main_image,
      sku_name: item.sku?.sku_name
    }));
  }
});

async function onSubmit() {
  if (!addressStore.selectedAddress) {
    uni.showToast({ title: '请选择收货地址', icon: 'none' });
    return;
  }
  
  if (!orderItems.value.length) {
    uni.showToast({ title: '订单商品不能为空', icon: 'none' });
    return;
  }
  
  submitting.value = true;
  
  try {
    // 创建订单
    const orderData = {
      region_code: regionStore.currentRegion,
      address_id: addressStore.selectedAddress.id,
      client_order_token: `${Date.now()}_${Math.random().toString(36).slice(2)}`,
      items: orderItems.value.map(item => ({
        product_id: item.product_id,
        sku_id: item.sku_id,
        quantity: item.quantity
      }))
    };
    
    const order = await orderStore.create(orderData);
    
    // 发起支付
    const payParams = await createPay(order.id);
    
    // 调起微信支付
    await new Promise((resolve, reject) => {
      uni.requestPayment({
        provider: 'wxpay',
        ...payParams,
        success: resolve,
        fail: reject
      });
    });
    
    // 支付成功
    uni.showToast({ title: '支付成功', icon: 'success' });
    setTimeout(() => {
      uni.navigateTo({ url: `/pages/order-detail/index?id=${order.id}` });
    }, 1500);
    
  } catch (e) {
    console.error('提交订单失败', e);
    if (e.errMsg?.includes('cancel')) {
      // 用户取消支付
      uni.showToast({ title: '已取消支付', icon: 'none' });
    } else {
      uni.showToast({ title: e.message || '提交失败', icon: 'none' });
    }
  } finally {
    submitting.value = false;
  }
}

function goAddressSelect() {
  uni.navigateTo({ url: '/pages/address/index?mode=select' });
}
</script>

<style lang="scss" scoped>
.page {
  min-height: 100vh;
  background: #F5F5F5;
  padding-bottom: 120rpx;
}

.address-card {
  margin: 24rpx;
  background: #fff;
  border-radius: 16rpx;
  padding: 30rpx;
  
  .address-info {
    .receiver {
      font-size: 28rpx;
      font-weight: 600;
      display: block;
      margin-bottom: 8rpx;
    }
    
    .detail {
      font-size: 26rpx;
      color: #666;
    }
  }
  
  .address-empty {
    display: flex;
    justify-content: space-between;
    align-items: center;
    color: #999;
    
    .arrow {
      font-size: 32rpx;
    }
  }
}

.card {
  margin: 0 24rpx 24rpx;
  background: #fff;
  border-radius: 16rpx;
  padding: 24rpx;
}

.section-title {
  font-size: 28rpx;
  font-weight: 600;
  margin-bottom: 20rpx;
}

.goods-item {
  display: flex;
  margin-bottom: 20rpx;
  
  &:last-child {
    margin-bottom: 0;
  }
  
  .goods-img {
    width: 160rpx;
    height: 160rpx;
    border-radius: 8rpx;
    background: #EEE;
    flex-shrink: 0;
  }
  
  .goods-info {
    flex: 1;
    margin-left: 20rpx;
    
    .goods-name {
      font-size: 28rpx;
      display: -webkit-box;
      -webkit-line-clamp: 2;
      -webkit-box-orient: vertical;
      overflow: hidden;
    }
    
    .goods-sku {
      font-size: 24rpx;
      color: #999;
      display: block;
      margin-top: 8rpx;
    }
    
    .goods-price-row {
      display: flex;
      justify-content: space-between;
      margin-top: 16rpx;
      
      .goods-price {
        font-size: 28rpx;
        color: #E4393C;
        font-weight: 600;
      }
      
      .goods-qty {
        font-size: 26rpx;
        color: #999;
      }
    }
  }
}

.fee-card {
  .fee-row {
    display: flex;
    justify-content: space-between;
    font-size: 26rpx;
    color: #666;
    padding: 12rpx 0;
  }
  
  .fee-total {
    display: flex;
    justify-content: space-between;
    font-size: 28rpx;
    font-weight: 600;
    padding-top: 20rpx;
    border-top: 1px solid #EEE;
    margin-top: 12rpx;
    
    .total-price {
      font-size: 36rpx;
      color: #E4393C;
    }
  }
}

.submit-bar {
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
  
  .total-info {
    font-size: 26rpx;
    
    .total-price {
      font-size: 36rpx;
      color: #E4393C;
      font-weight: 700;
    }
  }
  
  .btn-submit {
    margin-left: auto;
    background: #E4393C;
    color: #fff;
    padding: 0 50rpx;
    height: 72rpx;
    border-radius: 36rpx;
    font-size: 28rpx;
    line-height: 72rpx;
    
    &:disabled {
      background: #CCC;
    }
  }
}
</style>
