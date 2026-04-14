<template>
  <view class="page">
    <!-- 购物车列表 -->
    <view class="cart-body">
      <view v-if="!cartStore.items.length" class="empty">
        <text>购物车为空</text>
        <button class="btn" @click="goShopping">去逛逛</button>
      </view>
      
      <view v-else>
        <view v-for="item in cartStore.items" :key="item.id" class="cart-item">
          <checkbox :checked="item.checked" @change="() => cartStore.toggleItem(item.id)" />
          <image class="item-img" :src="item.product?.main_image || item.sku?.image" mode="aspectFill" />
          <view class="item-info">
            <text class="item-name">{{ item.product?.name || item.sku?.name }}</text>
            <text class="item-sku">{{ item.sku?.sku_name }}</text>
            <view class="item-bottom">
              <text class="item-price">¥{{ item.sku?.price || item.product?.price }}</text>
              <view class="qty-editor">
                <view class="qty-btn" @click="decrease(item)">-</view>
                <text class="qty-num">{{ item.quantity }}</text>
                <view class="qty-btn" @click="increase(item)">+</view>
              </view>
            </view>
          </view>
          <view class="item-delete" @click="removeItem(item)">删除</view>
        </view>
      </view>
    </view>
    
    <!-- 底部结算栏 -->
    <view v-if="cartStore.items.length" class="cart-footer">
      <view class="footer-left">
        <checkbox :checked="cartStore.allChecked" @change="() => cartStore.toggleAll(!cartStore.allChecked)" />
        <text>全选</text>
      </view>
      <view class="footer-total">
        合计: <text class="total-price">¥{{ cartStore.totalPrice }}</text>
      </view>
      <button class="btn-checkout" :disabled="!cartStore.checkedItems.length" @click="checkout">结算</button>
    </view>
  </view>
</template>

<script setup>
import { onShow } from '@dcloudio/uni-app';
import { useCartStore } from '@/store/cart';
import { useUserStore } from '@/store/user';

const cartStore = useCartStore();
const userStore = useUserStore();

onShow(async () => {
  if (userStore.isLoggedIn) {
    await cartStore.fetchCart();
  }
});

async function decrease(item) {
  if (item.quantity <= 1) return;
  await cartStore.updateItemQuantity(item.id, item.quantity - 1);
}

async function increase(item) {
  await cartStore.updateItemQuantity(item.id, item.quantity + 1);
}

async function removeItem(item) {
  await cartStore.removeItem(item.id);
}

function checkout() {
  if (!userStore.isLoggedIn) {
    uni.showToast({ title: '请先登录', icon: 'none' });
    return;
  }
  if (!cartStore.checkedItems.length) {
    uni.showToast({ title: '请选择商品', icon: 'none' });
    return;
  }
  uni.navigateTo({ url: '/pages/confirm/index' });
}

function goShopping() {
  uni.switchTab({ url: '/pages/index/index' });
}
</script>

<style lang="scss" scoped>
.page {
  min-height: 100vh;
  background: #F5F5F5;
  padding-bottom: 120rpx;
}

.cart-body {
  padding: 20rpx;
}

.cart-item {
  display: flex;
  align-items: center;
  background: #fff;
  border-radius: 16rpx;
  padding: 20rpx;
  margin-bottom: 20rpx;
  
  .item-img {
    width: 160rpx;
    height: 160rpx;
    border-radius: 8rpx;
    margin: 0 20rpx;
    background: #EEE;
  }
  
  .item-info {
    flex: 1;
    
    .item-name {
      font-size: 28rpx;
      display: block;
      margin-bottom: 8rpx;
    }
    
    .item-sku {
      font-size: 24rpx;
      color: #999;
      display: block;
      margin-bottom: 16rpx;
    }
    
    .item-bottom {
      display: flex;
      justify-content: space-between;
      align-items: center;
      
      .item-price {
        font-size: 32rpx;
        color: #E4393C;
        font-weight: 700;
      }
      
      .qty-editor {
        display: flex;
        align-items: center;
        
        .qty-btn {
          width: 48rpx;
          height: 48rpx;
          border: 1px solid #DDD;
          border-radius: 8rpx;
          display: flex;
          align-items: center;
          justify-content: center;
          font-size: 28rpx;
        }
        
        .qty-num {
          padding: 0 20rpx;
          font-weight: 700;
        }
      }
    }
  }
  
  .item-delete {
    font-size: 24rpx;
    color: #E4393C;
    margin-left: 20rpx;
  }
}

.cart-footer {
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
  
  .footer-left {
    display: flex;
    align-items: center;
    gap: 8rpx;
    font-size: 26rpx;
  }
  
  .footer-total {
    flex: 1;
    text-align: right;
    margin-right: 20rpx;
    font-size: 26rpx;
    
    .total-price {
      font-size: 36rpx;
      color: #E4393C;
      font-weight: 700;
    }
  }
  
  .btn-checkout {
    background: #E4393C;
    color: #fff;
    padding: 0 40rpx;
    height: 72rpx;
    border-radius: 36rpx;
    font-size: 28rpx;
    line-height: 72rpx;
    
    &:disabled {
      background: #CCC;
    }
  }
}

.empty {
  text-align: center;
  padding: 200rpx 0;
  color: #999;
  
  .btn {
    margin-top: 40rpx;
    background: #E4393C;
    color: #fff;
    padding: 20rpx 60rpx;
    border-radius: 999rpx;
  }
}
</style>
