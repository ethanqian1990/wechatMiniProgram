<template>
  <view class="page">
    <!-- 顶部区域选择栏 -->
    <view class="topbar" @click="showRegionPicker = true">
      <text class="region-name"><b>{{ regionStore.regionName }}</b> · 配送区域</text>
      <text class="switch-btn">切换 ▼</text>
    </view>
    
    <!-- 搜索栏 -->
    <view class="search-wrap" @click="goSearch">
      <view class="search-inner">
        <text class="search-icon">🔍</text>
        <text class="placeholder">搜索好物、品牌、品类</text>
      </view>
    </view>
    
    <!-- Banner 轮播 -->
    <swiper class="banner" :indicator-dots="true" :autoplay="true" :interval="3000" :circular="true" indicator-active-color="#E4393C">
      <swiper-item v-for="banner in banners" :key="banner.id" @click="onBannerClick(banner)">
        <view class="banner-slide" :style="{ backgroundImage: banner.image ? `url(${banner.image})` : banner.bg }">
          <view class="banner-overlay"></view>
          <view class="banner-content">
            <text v-if="banner.kicker" class="kicker">{{ banner.kicker }}</text>
            <text class="title">{{ banner.title }}</text>
            <text v-if="banner.subtitle" class="subtitle">{{ banner.subtitle }}</text>
            <text v-if="banner.promo" class="promo">{{ banner.promo }}</text>
            <view class="cta-btn">{{ banner.cta || '立即选购' }} ›</view>
          </view>
        </view>
      </swiper-item>
    </swiper>
    
    <!-- 分类横滑 -->
    <scroll-view class="category-scroll" scroll-x="true">
      <view 
        v-for="cat in categoryList" 
        :key="cat.id || 'all'"
        class="category-pill"
        :class="{ active: currentCategory === (cat.id || '') }"
        @click="onCategoryChange(cat.id || '')"
      >
        {{ cat.name }}
      </view>
    </scroll-view>
    
    <!-- 主推区 -->
    <view class="section-header">
      <view class="section-title">
        <view class="section-tag">主推</view>
        <text>本区销量第一</text>
      </view>
      <text class="section-more" @click="goCategory">进分类 ›</text>
    </view>
    
    <view v-if="featured" class="featured-card" @click="goDetail(featured.id)">
      <image class="featured-img" :src="featured.main_image || featured.image" mode="aspectFill" />
      <view class="featured-body">
        <view v-if="featured.tags?.length" class="featured-tag">{{ featured.tags[0] }}</view>
        <text class="featured-name">{{ featured.name }}</text>
        <text class="featured-desc">{{ featured.description }}</text>
        <view class="featured-price-row">
          <text class="featured-price">¥{{ featured.price }}</text>
          <text v-if="featured.original_price" class="featured-original">¥{{ featured.original_price }}</text>
          <view class="featured-buy">立即购买 ›</view>
        </view>
      </view>
    </view>
    
    <!-- 列表区 -->
    <view class="section-header">
      <view class="section-title">
        <view class="section-tag" style="background: #111;">列表</view>
        <text>其余在售</text>
      </view>
      <text class="section-more" @click="goSearchPage">搜索 ›</text>
    </view>
    
    <view class="goods-grid">
      <view v-for="goods in products" :key="goods.id" class="goods-card" @click="goDetail(goods.id)">
        <image class="goods-img" :src="goods.main_image" mode="aspectFill" />
        <view class="goods-body">
          <text v-if="goods.tags?.length" class="goods-tag">{{ goods.tags[0] }}</text>
          <text class="goods-name">{{ goods.name }}</text>
          <text v-if="goods.promo_text" class="goods-promo">{{ goods.promo_text }}</text>
          <view class="goods-price-row">
            <view>
              <text class="goods-price">¥{{ goods.price }}</text>
              <text v-if="goods.original_price" class="goods-original">¥{{ goods.original_price }}</text>
            </view>
            <view class="cart-btn" @click.stop="addToCart(goods)">🛒</view>
          </view>
        </view>
      </view>
    </view>
    
    <!-- 加载更多 -->
    <view v-if="hasMore" class="load-more" @click="loadMore">
      <text>加载更多</text>
    </view>
    <view v-else class="no-more">
      <text>没有更多了</text>
    </view>
    
    <!-- 区域选择弹窗 -->
    <view v-if="showRegionPicker" class="region-overlay" @click="showRegionPicker = false">
      <view class="region-panel" @click.stop>
        <view class="region-title">选择区域</view>
        <view class="region-tags">
          <view 
            v-for="region in regionStore.regionList" 
            :key="region.code"
            class="region-tag"
            :class="{ active: region.code === regionStore.currentRegion }"
            @click="selectRegion(region.code)"
          >
            {{ region.name }}
          </view>
        </view>
        <view class="region-actions">
          <button class="btn-cancel" @click="showRegionPicker = false">取消</button>
          <button class="btn-confirm" @click="confirmRegion">确定</button>
        </view>
      </view>
    </view>
  </view>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue';
import { useRegionStore } from '@/store/region';
import { useCategoryStore } from '@/store/category';
import { useCartStore } from '@/store/cart';
import { getHomeBanners, getHomeFeatured, getHomeProducts } from '@/api/home';

const regionStore = useRegionStore();
const categoryStore = useCategoryStore();
const cartStore = useCartStore();

const banners = ref([]);
const featured = ref(null);
const products = ref([]);
const currentCategory = ref('');
const page = ref(1);
const hasMore = ref(true);
const showRegionPicker = ref(false);
const pendingRegion = ref('');

const categoryList = computed(() => {
  return [{ id: '', name: '全部' }, ...categoryStore.categories];
});

onMounted(async () => {
  pendingRegion.value = regionStore.currentRegion;
  await loadHomeData();
});

async function loadHomeData() {
  try {
    const regionCode = regionStore.currentRegion;
    const categoryId = currentCategory.value;
    
    const [bannerData, featuredData, productsData] = await Promise.all([
      getHomeBanners(regionCode),
      getHomeFeatured(regionCode, categoryId),
      getHomeProducts(regionCode, categoryId)
    ]);
    
    banners.value = bannerData.banners || [];
    featured.value = featuredData || null;
    products.value = productsData.list || productsData || [];
    hasMore.value = (productsData.list || productsData || []).length >= 24;
  } catch (e) {
    console.error('加载首页数据失败', e);
    uni.showToast({ title: '加载失败', icon: 'none' });
  }
}

async function onCategoryChange(catId) {
  currentCategory.value = catId;
  page.value = 1;
  await loadHomeData();
}

function selectRegion(code) {
  pendingRegion.value = code;
}

async function confirmRegion() {
  regionStore.setRegion(pendingRegion.value);
  showRegionPicker = false;
  page.value = 1;
  await loadHomeData();
}

function onBannerClick(banner) {
  if (!banner.goto) return;
  const { type, product_id, keyword, category_id } = banner.goto;
  if (type === 'product' && product_id) {
    goDetail(product_id);
  } else if (type === 'search' && keyword) {
    uni.navigateTo({ url: `/pages/search/result?keyword=${keyword}` });
  } else if (type === 'category' && category_id) {
    uni.switchTab({ url: '/pages/category/index' });
  }
}

function goDetail(id) {
  uni.navigateTo({ url: `/pages/detail/index?id=${id}` });
}

function goSearch() {
  uni.navigateTo({ url: '/pages/search/index' });
}

function goSearchPage() {
  uni.navigateTo({ url: '/pages/search/result' });
}

function goCategory() {
  uni.switchTab({ url: '/pages/category/index' });
}

async function addToCart(goods) {
  if (!goods.skus?.length) {
    uni.showToast({ title: '暂无可售规格', icon: 'none' });
    return;
  }
  const sku = goods.skus[0];
  try {
    await cartStore.addItem(goods.id, sku.id, 1);
  } catch (e) {
    console.error('加入购物车失败', e);
  }
}

async function loadMore() {
  page.value++;
  try {
    const regionCode = regionStore.currentRegion;
    const categoryId = currentCategory.value;
    const data = await getHomeProducts(regionCode, categoryId);
    const newProducts = data.list || data || [];
    products.value = [...products.value, ...newProducts];
    hasMore.value = newProducts.length >= 24;
  } catch (e) {
    page.value--;
    uni.showToast({ title: '加载失败', icon: 'none' });
  }
}
</script>

<style lang="scss" scoped>
.page {
  min-height: 100vh;
  background: #F5F5F5;
  padding-bottom: 20rpx;
}

.topbar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16rpx 24rpx 8rpx;
  font-size: 24rpx;
  color: #666;
  
  .region-name b {
    color: #111;
    font-size: 26rpx;
  }
  
  .switch-btn {
    color: #E4393C;
  }
}

.search-wrap {
  padding: 0 24rpx 20rpx;
}

.search-inner {
  display: flex;
  align-items: center;
  background: #F0F0F0;
  border-radius: 44rpx;
  padding: 20rpx 28rpx;
  
  .search-icon {
    margin-right: 16rpx;
  }
  
  .placeholder {
    color: #999;
    font-size: 26rpx;
  }
}

.banner {
  height: 400rpx;
  margin: 0 24rpx;
  border-radius: 28rpx;
  overflow: hidden;
  
  .banner-slide {
    width: 100%;
    height: 100%;
    background-size: cover;
    background-position: center;
    display: flex;
    align-items: flex-end;
    position: relative;
  }
  
  .banner-overlay {
    position: absolute;
    inset: 0;
    background: linear-gradient(to top, rgba(0,0,0,0.78) 0%, rgba(0,0,0,0.35) 42%, transparent 72%);
  }
  
  .banner-content {
    position: relative;
    z-index: 1;
    padding: 28rpx 32rpx;
    color: #fff;
    width: 100%;
    
    .kicker {
      font-size: 22rpx;
      opacity: 0.92;
      display: block;
      margin-bottom: 8rpx;
    }
    
    .title {
      font-size: 34rpx;
      font-weight: 800;
      display: block;
      margin-bottom: 8rpx;
      text-shadow: 0 2rpx 16rpx rgba(0,0,0,0.35);
    }
    
    .subtitle {
      font-size: 24rpx;
      opacity: 0.88;
      display: block;
      margin-bottom: 12rpx;
    }
    
    .promo {
      font-size: 24rpx;
      font-weight: 700;
      color: #FFB74D;
      display: block;
      margin-bottom: 16rpx;
    }
    
    .cta-btn {
      display: inline-block;
      padding: 12rpx 28rpx;
      border-radius: 999rpx;
      border: 1px solid rgba(255,255,255,0.85);
      font-size: 24rpx;
      font-weight: 600;
      background: rgba(255,255,255,0.12);
      backdrop-filter: blur(8rpx);
    }
  }
}

.category-scroll {
  white-space: nowrap;
  padding: 0 24rpx 20rpx;
}

.category-pill {
  display: inline-block;
  padding: 16rpx 28rpx;
  border-radius: 999rpx;
  background: #fff;
  border: 1px solid #EEE;
  font-size: 24rpx;
  color: #555;
  margin-right: 16rpx;
  
  &.active {
    background: #111;
    color: #fff;
    border-color: #111;
  }
}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12rpx 24rpx 20rpx;
  
  .section-title {
    display: flex;
    align-items: center;
    gap: 16rpx;
    
    .section-tag {
      font-size: 20rpx;
      font-weight: 700;
      color: #fff;
      background: #E4393C;
      padding: 4rpx 12rpx;
      border-radius: 8rpx;
    }
    
    text {
      font-size: 30rpx;
      font-weight: 700;
      color: #111;
    }
  }
  
  .section-more {
    font-size: 24rpx;
    color: #999;
  }
}

.featured-card {
  display: flex;
  margin: 0 24rpx 24rpx;
  background: #fff;
  border-radius: 28rpx;
  overflow: hidden;
  box-shadow: 0 4rpx 24rpx rgba(0,0,0,0.06);
  
  .featured-img {
    width: 236rpx;
    height: 236rpx;
    flex-shrink: 0;
  }
  
  .featured-body {
    flex: 1;
    padding: 24rpx;
    display: flex;
    flex-direction: column;
    justify-content: center;
    
    .featured-tag {
      font-size: 20rpx;
      color: #E4393C;
      border: 1px solid rgba(228,57,60,0.45);
      padding: 4rpx 12rpx;
      border-radius: 8rpx;
      align-self: flex-start;
      margin-bottom: 12rpx;
    }
    
    .featured-name {
      font-size: 28rpx;
      font-weight: 700;
      color: #222;
      display: -webkit-box;
      -webkit-line-clamp: 2;
      -webkit-box-orient: vertical;
      overflow: hidden;
    }
    
    .featured-desc {
      font-size: 22rpx;
      color: #999;
      margin-top: 8rpx;
      display: -webkit-box;
      -webkit-line-clamp: 1;
      -webkit-box-orient: vertical;
      overflow: hidden;
    }
    
    .featured-price-row {
      display: flex;
      align-items: baseline;
      gap: 16rpx;
      margin-top: 16rpx;
      
      .featured-price {
        font-size: 36rpx;
        font-weight: 800;
        color: #E4393C;
      }
      
      .featured-original {
        font-size: 24rpx;
        color: #bbb;
        text-decoration: line-through;
      }
      
      .featured-buy {
        margin-left: auto;
        padding: 10rpx 24rpx;
        border-radius: 999rpx;
        background: #F5F5F5;
        color: #333;
        font-size: 24rpx;
        font-weight: 600;
      }
    }
  }
}

.goods-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 16rpx;
  padding: 0 24rpx;
}

.goods-card {
  background: #fff;
  border-radius: 24rpx;
  overflow: hidden;
  box-shadow: 0 2rpx 12rpx rgba(0,0,0,0.06);
  
  .goods-img {
    width: 100%;
    aspect-ratio: 1;
    background: #EEE;
  }
  
  .goods-body {
    padding: 16rpx;
    
    .goods-tag {
      font-size: 18rpx;
      color: #E4393C;
      border: 1px solid rgba(228,57,60,0.4);
      padding: 2rpx 8rpx;
      border-radius: 6rpx;
      display: inline-block;
      margin-bottom: 8rpx;
    }
    
    .goods-name {
      font-size: 22rpx;
      font-weight: 600;
      color: #222;
      display: -webkit-box;
      -webkit-line-clamp: 2;
      -webkit-box-orient: vertical;
      overflow: hidden;
      line-height: 1.3;
    }
    
    .goods-promo {
      font-size: 18rpx;
      color: #E4393C;
      margin-top: 4rpx;
      display: block;
    }
    
    .goods-price-row {
      display: flex;
      align-items: center;
      justify-content: space-between;
      margin-top: 12rpx;
      
      .goods-price {
        font-size: 26rpx;
        font-weight: 800;
        color: #E4393C;
      }
      
      .goods-original {
        font-size: 18rpx;
        color: #bbb;
        text-decoration: line-through;
        margin-left: 8rpx;
      }
      
      .cart-btn {
        width: 44rpx;
        height: 44rpx;
        border-radius: 50%;
        background: #FFECEB;
        color: #E4393C;
        font-size: 24rpx;
        display: flex;
        align-items: center;
        justify-content: center;
      }
    }
  }
}

.load-more, .no-more {
  text-align: center;
  padding: 32rpx;
  color: #999;
  font-size: 24rpx;
}

.region-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0,0,0,0.45);
  z-index: 999;
  display: flex;
  align-items: center;
  justify-content: center;
}

.region-panel {
  width: 640rpx;
  background: #fff;
  border-radius: 24rpx;
  padding: 32rpx;
  
  .region-title {
    font-size: 32rpx;
    font-weight: 700;
    margin-bottom: 24rpx;
  }
  
  .region-tags {
    display: grid;
    grid-template-columns: repeat(2, 1fr);
    gap: 16rpx;
    
    .region-tag {
      padding: 24rpx;
      border-radius: 12rpx;
      background: #F5F5F5;
      text-align: center;
      font-size: 28rpx;
      
      &.active {
        background: #FFECEB;
        color: #E4393C;
        border: 1px solid #E4393C;
      }
    }
  }
  
  .region-actions {
    display: flex;
    gap: 16rpx;
    margin-top: 24rpx;
    
    button {
      flex: 1;
      padding: 24rpx;
      border-radius: 12rpx;
      font-size: 28rpx;
    }
    
    .btn-cancel {
      background: #F5F5F5;
      color: #333;
    }
    
    .btn-confirm {
      background: #E4393C;
      color: #fff;
    }
  }
}
</style>
