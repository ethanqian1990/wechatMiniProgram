<template>
  <view class="page">
    <!-- 搜索栏 -->
    <view class="search-bar">
      <input v-model="keyword" class="search-input" placeholder="搜索商品（回车搜索）" confirm-type="search" @confirm="onSearch" @input="onInput" />
    </view>
    
    <!-- 搜索建议 -->
    <view v-if="suggestions.length" class="suggest">
      <view v-for="s in suggestions" :key="s" class="suggest-item" @click="onSuggestClick(s)">
        {{ s }}
      </view>
    </view>
    
    <!-- 热门搜索 -->
    <view v-if="!keyword && !suggestions.length" class="section">
      <view class="section-title">
        <text>热门搜索</text>
      </view>
      <view class="tags">
        <view v-for="tag in hotTags" :key="tag" class="tag" @click="onTagClick(tag)">
          {{ tag }}
        </view>
      </view>
    </view>
    
    <!-- 搜索历史 -->
    <view v-if="!keyword && !suggestions.length && history.length" class="section">
      <view class="section-title">
        <text>搜索历史</text>
        <text class="clear" @click="clearHistory">清空</text>
      </view>
      <view class="tags">
        <view v-for="h in history" :key="h" class="tag" @click="onTagClick(h)">
          {{ h }}
        </view>
      </view>
    </view>
  </view>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import { getSearchSuggest, getSearchHot, getSearchHistory, clearSearchHistory } from '@/api/search';

const keyword = ref('');
const suggestions = ref([]);
const hotTags = ref([]);
const history = ref([]);

let suggestTimer = null;

onMounted(async () => {
  await Promise.all([loadHot(), loadHistory()]);
});

async function loadHot() {
  try {
    const data = await getSearchHot();
    hotTags.value = data.list || data || [];
  } catch (e) {
    hotTags.value = ['大米', '苹果', '草莓', '鸡蛋', '蜂蜜', '蔬菜'];
  }
}

async function loadHistory() {
  try {
    const data = await getSearchHistory();
    history.value = data.list || data || [];
  } catch (e) {
    history.value = [];
  }
}

function onInput() {
  if (suggestTimer) clearTimeout(suggestTimer);
  suggestTimer = setTimeout(async () => {
    if (!keyword.value.trim()) {
      suggestions.value = [];
      return;
    }
    try {
      const data = await getSearchSuggest(keyword.value);
      suggestions.value = data.list || data || [];
    } catch (e) {
      suggestions.value = [];
    }
  }, 300);
}

function onSearch() {
  if (!keyword.value.trim()) return;
  uni.navigateTo({ url: `/pages/search/result?keyword=${encodeURIComponent(keyword.value)}` });
}

function onTagClick(tag) {
  keyword.value = tag;
  onSearch();
}

function onSuggestClick(s) {
  keyword.value = s;
  onSearch();
}

async function clearHistory() {
  try {
    await clearSearchHistory();
    history.value = [];
  } catch (e) {
    console.error('清空失败', e);
  }
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

.suggest {
  background: #fff;
  padding: 0 24rpx;
  
  .suggest-item {
    padding: 24rpx 0;
    border-bottom: 1px solid #F5F5F5;
    font-size: 26rpx;
  }
}

.section {
  margin: 24rpx;
  background: #fff;
  border-radius: 16rpx;
  padding: 24rpx;
  
  .section-title {
    font-size: 28rpx;
    font-weight: 600;
    margin-bottom: 20rpx;
    display: flex;
    justify-content: space-between;
    
    .clear {
      font-weight: 400;
      color: #E4393C;
    }
  }
  
  .tags {
    display: flex;
    flex-wrap: wrap;
    gap: 16rpx;
    
    .tag {
      padding: 12rpx 24rpx;
      background: #F5F5F5;
      border-radius: 999rpx;
      font-size: 24rpx;
    }
  }
}
</style>
