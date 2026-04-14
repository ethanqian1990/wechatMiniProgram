<template>
  <view class="page">
    <!-- 商品轮播图 -->
    <swiper class="gallery" :indicator-dots="true" :autoplay="true">
      <swiper-item v-for="(img, idx) in images" :key="idx">
        <image :src="img" mode="aspectFill" class="gallery-img" />
      </swiper-item>
    </swiper>
    
    <!-- 商品信息 -->
    <view class="goods-info card">
      <view class="price-row">
        <text class="price">¥{{ currentSku?.price || goods.price }}</text>
        <text v-if="currentSku?.original_price || goods.original_price" class="original-price">
          ¥{{ currentSku?.original_price || goods.original_price }}
        </text>
      </view>
      <text class="name">{{ goods.name }}</text>
      <text v-if="goods.description" class="desc">{{ goods.description }}</text>
    </view>
    
    <!-- SKU 选择 -->
    <view v-if="goods.skus?.length" class="card">
      <view class="section-title">规格</view>
      <view class="sku-list">
        <view 
          v-for="sku in goods.skus" 
          :key="sku.id"
          class="sku-item"
          :class="{ active: currentSkuId === sku.id, disabled: sku.stock <= 0 }"
          @click="selectSku(sku)"
        >
          {{ sku.sku_name }}
          <text v-if="sku.stock <= 0" class="stock-tip">售罄</text>
        </view>
      </view>
    </view>
    
    <!-- 数量选择 -->
    <view class="card">
      <view class="section-title">数量</view>
      <view class="qty-row">
        <view class="qty-btn" @click="decrease">-</view>
        <text class="qty-num">{{ quantity }}</text>
        <view class="qty-btn" @click="increase">+</view>
        <text class="stock-info">库存：{{ currentSku?.stock || goods.stock || 0 }}</text>
      </view>
    </view>
    
    <!-- 底部操作栏 -->
    <view class="action-bar">
      <view class="action-left">
        <view class="action-icon" @click="goCart">🛒</view>
      </view>
      <view class="action-right">
        <button class="btn-cart" @click="onAddCart">加入购物车</button>
        <button class="btn-buy" @click="onBuyNow">立即购买</button>
      </view>
    </view>
  </view>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue';
import { useCartStore } from '@/store/cart';
import { useUserStore } from '@/store/user';
import { getProduct } from '@/api/goods';

const cartStore = useCartStore();
const userStore = useUserStore();

const goods = ref({});
const currentSkuId = ref('');
const quantity = ref(1);

const currentSku = computed(() => {
  return goods.value.skus?.find(s => s.id === currentSkuId.value);
});

const images = computed(() => {
  const imgs = [];
  if (goods.value.main_image) imgs.push(goods.value.main_image);
  if (goods.value.images?.length) imgs.push(...goods.value.images);
  return imgs.length ? imgs : ['/static/images/default.png'];
});

onMounted(async () => {
  const pages = getCurrentPages();
  const id = pages[pages.length - 1].options.id;
  if (id) {
    await loadGoods(id);
  }
});

async function loadGoods(id) {
  try {
    const data = await getProduct(id);
    goods.value = data;
    if (data.skus?.length) {
      // 默认选第一个有库存的 SKU
      const availableSku = data.skus.find(s => s.stock > 0) || data.skus[0];
      currentSkuId.value = availableSku.id;
    }
  } catch (e) {
    console.error('加载商品失败', e);
    uni.showToast({ title: '加载失败', icon: 'none' });
  }
}

function selectSku(sku) {
  if (sku.stock <= 0) return;
  currentSkuId.value = sku.id;
}

function decrease() {
  if (quantity.value <= 1) return;
  quantity.value--;
}

function increase() {
  const maxStock = currentSku.value?.stock || goods.value.stock || 0;
  if (quantity.value >= maxStock) {
    uni.showToast({ title: '库存不足', icon: 'none' });
    return;
  }
  quantity.value++;
}

async function onAddCart() {
  if (!userStore.isLoggedIn) {
    uni.showToast({ title: '请先登录', icon: 'none' });
    return;
  }
  if (!currentSkuId.value) {
    uni.showToast({ title: '请选择规格', icon: 'none' });
    return;
  }
  await cartStore.addItem(goods.value.id, currentSkuId.value, quantity.value);
}

function onBuyNow() {
  if (!userStore.isLoggedIn) {
    uni.showToast({ title: '请先登录', icon: 'none' });
    return;
  }
  if (!currentSkuId.value) {
    uni.showToast({ title: '请选择规格', icon: 'none' });
    return;
  }
  
  // 跳转到确认订单页，带上商品信息
  uni.navigateTo({
    url: `/pages/confirm/index?goodsId=${goods.value.id}&skuId=${currentSkuId.value}&quantity=${quantity.value}`
  });
}

function goCart() {
  uni.switchTab({ url: '/pages/cart/index' });
}
</script>

<style lang="scss" scoped>
.page {
  min-height: 100vh;
  background: #F5F5F5;
  padding-bottom: 120rpx;
}

.gallery {
  height: 600rpx;
  
  .gallery-img {
    width: 100%;
    height: 100%;
  }
}

.card {
  margin: 24rpx;
  background: #fff;
  border-radius: 16rpx;
  padding: 24rpx;
}

.goods-info {
  .price-row {
    display: flex;
    align-items: baseline;
    gap: 16rpx;
    margin-bottom: 16rpx;
    
    .price {
      font-size: 44rpx;
      font-weight: 800;
      color: #E4393C;
    }
    
    .original-price {
      font-size: 28rpx;
      color: #999;
      text-decoration: line-through;
    }
  }
  
  .name {
    font-size: 32rpx;
    font-weight: 600;
    display: block;
    margin-bottom: 12rpx;
  }
  
  .desc {
    font-size: 26rpx;
    color: #666;
    display: block;
  }
}

.section-title {
  font-size: 28rpx;
  font-weight: 600;
  margin-bottom: 20rpx;
}

.sku-list {
  display: flex;
  flex-wrap: wrap;
  gap: 16rpx;
  
  .sku-item {
    padding: 16rpx 24rpx;
    border-radius: 8rpx;
    background: #F5F5F5;
    font-size: 26rpx;
    position: relative;
    
    &.active {
      background: #FFECEB;
      color: #E4393C;
      border: 1px solid #E4393C;
    }
    
    &.disabled {
      opacity: 0.5;
    }
    
    .stock-tip {
      font-size: 20rpx;
      color: #999;
    }
  }
}

.qty-row {
  display: flex;
  align-items: center;
  gap: 24rpx;
  
  .qty-btn {
    width: 72rpx;
    height: 72rpx;
    border: 1px solid #DDD;
    border-radius: 20rpx;
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 32rpx;
  }
  
  .qty-num {
    font-weight: 700;
    font-size: 32rpx;
    min-width: 60rpx;
    text-align: center;
  }
  
  .stock-info {
    font-size: 24rpx;
    color: #999;
    margin-left: auto;
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
    .action-icon {
      font-size: 44rpx;
      padding: 10rpx;
    }
  }
  
  .action-right {
    flex: 1;
    display: flex;
    gap: 16rpx;
    margin-left: 24rpx;
    
    button {
      flex: 1;
      height: 72rpx;
      border-radius: 36rpx;
      font-size: 28rpx;
      line-height: 72rpx;
    }
    
    .btn-cart {
      background: #FFECEB;
      color: #E4393C;
    }
    
    .btn-buy {
      background: #E4393C;
      color: #fff;
    }
  }
}
</style>
