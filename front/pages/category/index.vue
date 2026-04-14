<template>
  <view class="page">
    <!-- 搜索栏 -->
    <view class="search-wrap">
      <input v-model="keyword" class="search-input" placeholder="搜索商品（回车搜索）" confirm-type="search" @confirm="onSearch" />
    </view>
    
    <!-- 分类侧边栏 -->
    <view class="category-layout">
      <scroll-view class="category-left" scroll-y="true">
        <view 
          v-for="cat in categoryList" 
          :key="cat.id"
          class="category-item"
          :class="{ active: currentCategory === cat.id }"
          @click="onCategoryChange(cat.id)"
        >
          {{ cat.name }}
        </view>
      </scroll-view>
      
      <!-- 商品列表 -->
      <scroll-view class="category-right" scroll-y="true">
        <view v-if="loading" class="loading">加载中...</view>
        <view v-else-if="!products.length" class="empty">该分类暂无商品</view>
        <view v-else class="goods-grid">
          <view v-for="goods in products" :key="goods.id" class="goods-card" @click="goDetail(goods.id)">
            <image class="goods-img" :src="goods.main_image" mode="aspectFill" />
            <view class="goods-info">
              <text class="goods-name">{{ goods.name }}</text>
              <text class="goods-price">¥{{ goods.price }}</text>
            </view>
          </view>
        </view>
      </scroll-view>
    </view>
  </view>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import { useCategoryStore } from '@/store/category';
import { useRegionStore } from '@/store/region';
import { getProducts } from '@/api/goods';

const categoryStore = useCategoryStore();
const regionStore = useRegionStore();

const keyword = ref('');
const currentCategory = ref('');
const products = ref([]);
const loading = ref(false);

const categoryList = computed(() => categoryStore.categories);

import { computed } from 'vue';

onMounted(async () => {
  await loadProducts();
});

async function loadProducts() {
  loading.value = true;
  try {
    const data = await getProducts({
      region_code: regionStore.currentRegion,
      category_id: currentCategory.value,
      keyword: keyword.value,
      on_shelf: 1
    });
    products.value = data.list || data || [];
  } catch (e) {
    console.error('获取商品失败', e);
  } finally {
    loading.value = false;
  }
}

function onCategoryChange(catId) {
  currentCategory.value = catId;
  loadProducts();
}

function onSearch() {
  loadProducts();
}

function goDetail(id) {
  uni.navigateTo({ url: `/pages/detail/index?id=${id}` });
}
</script>

<style lang="scss" scoped>
.page {
  height: 100vh;
  display: flex;
  flex-direction: column;
}

.search-wrap {
  padding: 16rpx 24rpx;
  background: #fff;
  
  .search-input {
    background: #F5F5F5;
    border-radius: 44rpx;
    padding: 20rpx 30rpx;
    font-size: 26rpx;
  }
}

.category-layout {
  flex: 1;
  display: flex;
  overflow: hidden;
}

.category-left {
  width: 160rpx;
  height: 100%;
  background: #F5F5F5;
  
  .category-item {
    padding: 30rpx 20rpx;
    text-align: center;
    font-size: 24rpx;
    color: #666;
    
    &.active {
      background: #fff;
      color: #E4393C;
      border-left: 6rpx solid #E4393C;
    }
  }
}

.category-right {
  flex: 1;
  height: 100%;
  padding: 20rpx;
}

.goods-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 20rpx;
}

.goods-card {
  background: #fff;
  border-radius: 16rpx;
  overflow: hidden;
  
  .goods-img {
    width: 100%;
    aspect-ratio: 1;
    background: #EEE;
  }
  
  .goods-info {
    padding: 16rpx;
    
    .goods-name {
      font-size: 26rpx;
      display: block;
      margin-bottom: 8rpx;
    }
    
    .goods-price {
      font-size: 28rpx;
      color: #E4393C;
      font-weight: 700;
    }
  }
}

.loading, .empty {
  text-align: center;
  padding: 100rpx;
  color: #999;
}
</style>
