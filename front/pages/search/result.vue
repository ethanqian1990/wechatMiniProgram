<template>
  <view class="page">
    <!-- 搜索栏 -->
    <view class="search-bar">
      <input v-model="keyword" class="search-input" placeholder="搜索商品（回车搜索）" confirm-type="search" @confirm="onSearch" />
    </view>
    
    <!-- 筛选 -->
    <view class="filter-bar">
      <view 
        v-for="f in filters" 
        :key="f.value"
        class="filter-chip"
        :class="{ active: sort === f.value }"
        @click="onFilter(f.value)"
      >
        {{ f.label }}
      </view>
    </view>
    
    <!-- 商品列表 -->
    <view v-if="loading" class="loading">加载中...</view>
    <view v-else-if="!products.length" class="empty">未找到相关商品</view>
    <view v-else class="goods-list">
      <view v-for="goods in products" :key="goods.id" class="goods-item" @click="goDetail(goods.id)">
        <image class="goods-img" :src="goods.main_image" mode="aspectFill" />
        <view class="goods-info">
          <text class="goods-name">{{ goods.name }}</text>
          <text class="goods-meta">销量 {{ goods.sales_count }} · 库存 {{ goods.stock }}</text>
          <view class="goods-bottom">
            <text class="goods-price">¥{{ goods.price }}</text>
            <view class="goods-actions">
              <button class="btn-cart" @click.stop="addToCart(goods)">加入购物车</button>
              <button class="btn-buy" @click.stop="buyNow(goods)">购买</button>
            </view>
          </view>
        </view>
      </view>
    </view>
  </view>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import { useRegionStore } from '@/store/region';
import { useCartStore } from '@/store/cart';
import { getProducts } from '@/api/goods';

const regionStore = useRegionStore();
const cartStore = useCartStore();

const keyword = ref('');
const sort = ref('default');
const products = ref([]);
const loading = ref(false);

const filters = [
  { label: '综合', value: 'default' },
  { label: '价格↑', value: 'price_asc' },
  { label: '价格↓', value: 'price_desc' },
  { label: '销量', value: 'sales' }
];

onMounted(() => {
  const pages = getCurrentPages();
  const currentPage = pages[pages.length - 1];
  if (currentPage.options.keyword) {
    keyword.value = decodeURIComponent(currentPage.options.keyword);
  }
  loadProducts();
});

async function loadProducts() {
  loading.value = true;
  try {
    const data = await getProducts({
      region_code: regionStore.currentRegion,
      keyword: keyword.value,
      sort: sort.value === 'default' ? '' : sort.value,
      on_shelf: 1
    });
    products.value = data.list || data || [];
  } catch (e) {
    console.error('搜索失败', e);
  } finally {
    loading.value = false;
  }
}

function onSearch() {
  loadProducts();
}

function onFilter(value) {
  sort.value = value;
  loadProducts();
}

function goDetail(id) {
  uni.navigateTo({ url: `/pages/detail/index?id=${id}` });
}

async function addToCart(goods) {
  if (!goods.skus?.length) {
    uni.showToast({ title: '暂无可售规格', icon: 'none' });
    return;
  }
  await cartStore.addItem(goods.id, goods.skus[0].id, 1);
}

function buyNow(goods) {
  goDetail(goods.id);
}
</script>

<style lang="scss" scoped>
.page {
  min-height: 100vh;
  background: #F5F5F5;
}

.search-bar {
  padding: 20rpx 24rpx;
  background: #fff;
  
  .search-input {
    background: #F5F5F5;
    border-radius: 44rpx;
    padding: 20rpx 30rpx;
    font-size: 26rpx;
  }
}

.filter-bar {
  display: flex;
  gap: 16rpx;
  padding: 20rpx 24rpx;
  background: #fff;
  
  .filter-chip {
    padding: 12rpx 24rpx;
    border-radius: 999rpx;
    background: #F5F5F5;
    font-size: 24rpx;
    color: #666;
    
    &.active {
      border: 1px solid #E4393C;
      color: #E4393C;
      background: #FFECEB;
    }
  }
}

.goods-list {
  padding: 20rpx;
}

.goods-item {
  display: flex;
  background: #fff;
  border-radius: 16rpx;
  padding: 20rpx;
  margin-bottom: 20rpx;
  
  .goods-img {
    width: 200rpx;
    height: 200rpx;
    border-radius: 8rpx;
    background: #EEE;
    flex-shrink: 0;
  }
  
  .goods-info {
    flex: 1;
    margin-left: 20rpx;
    display: flex;
    flex-direction: column;
    
    .goods-name {
      font-size: 28rpx;
      display: -webkit-box;
      -webkit-line-clamp: 2;
      -webkit-box-orient: vertical;
      overflow: hidden;
    }
    
    .goods-meta {
      font-size: 24rpx;
      color: #999;
      margin-top: 8rpx;
    }
    
    .goods-bottom {
      margin-top: auto;
      display: flex;
      justify-content: space-between;
      align-items: center;
      
      .goods-price {
        font-size: 32rpx;
        color: #E4393C;
        font-weight: 700;
      }
      
      .goods-actions {
        display: flex;
        gap: 12rpx;
        
        button {
          padding: 12rpx 20rpx;
          border-radius: 8rpx;
          font-size: 22rpx;
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
  }
}

.loading, .empty {
  text-align: center;
  padding: 200rpx;
  color: #999;
}
</style>
