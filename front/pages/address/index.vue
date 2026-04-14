<template>
  <view class="page">
    <!-- 地址列表 -->
    <view v-if="addressStore.addresses.length" class="address-list">
      <view 
        v-for="addr in addressStore.addresses" 
        :key="addr.id" 
        class="address-card"
        :class="{ selected: addr.id === addressStore.selectedAddressId }"
        @click="onSelect(addr)"
      >
        <view class="address-main">
          <view class="address-top">
            <text class="receiver">{{ addr.receiver }}</text>
            <text class="phone">{{ addr.phone }}</text>
            <view v-if="addr.is_default" class="default-tag">默认</view>
          </view>
          <text class="address-detail">{{ addr.province }}{{ addr.city }}{{ addr.district }}{{ addr.detail }}</text>
        </view>
        
        <view class="address-actions">
          <view v-if="mode === 'select'" class="action-btn select-btn" @click.stop="onSelect(addr)">选择</view>
          <view v-if="!addr.is_default" class="action-btn" @click.stop="onSetDefault(addr.id)">设为默认</view>
          <view class="action-btn" @click.stop="onEdit(addr)">编辑</view>
          <view class="action-btn delete" @click.stop="onDelete(addr.id)">删除</view>
        </view>
      </view>
    </view>
    
    <!-- 空状态 -->
    <view v-else class="empty">
      <text>暂无收货地址</text>
    </view>
    
    <!-- 新增按钮 -->
    <view class="add-btn-wrap">
      <button class="btn-add" @click="onAdd">+ 新增收货地址</button>
    </view>
  </view>
</template>

<script setup>
import { onMounted } from 'vue';
import { useAddressStore } from '@/store/address';

const addressStore = useAddressStore();

let mode = 'list'; // list | select

onMounted(async () => {
  const pages = getCurrentPages();
  const options = pages[pages.length - 1].options;
  mode = options.mode || 'list';
  
  await addressStore.fetchAddresses();
});

function onSelect(addr) {
  addressStore.selectAddress(addr.id);
  if (mode === 'select') {
    uni.navigateBack();
  }
}

async function onSetDefault(id) {
  await addressStore.setDefault(id);
  uni.showToast({ title: '设置成功', icon: 'success' });
}

function onEdit(addr) {
  uni.navigateTo({ url: `/pages/address-edit/index?address=${encodeURIComponent(JSON.stringify(addr))}` });
}

async function onDelete(id) {
  uni.showModal({
    title: '提示',
    content: '确定删除该收货地址？',
    success: async (res) => {
      if (res.confirm) {
        try {
          await addressStore.remove(id);
          uni.showToast({ title: '已删除', icon: 'success' });
        } catch (e) {
          uni.showToast({ title: '删除失败', icon: 'none' });
        }
      }
    }
  });
}

function onAdd() {
  uni.navigateTo({ url: '/pages/address-edit/index' });
}
</script>

<style lang="scss" scoped>
.page {
  min-height: 100vh;
  background: #F5F5F5;
  padding: 24rpx;
  padding-bottom: 120rpx;
}

.address-list {
  .address-card {
    background: #fff;
    border-radius: 16rpx;
    padding: 24rpx;
    margin-bottom: 20rpx;
    
    &.selected {
      border: 2rpx solid #E4393C;
    }
    
    .address-main {
      margin-bottom: 20rpx;
      
      .address-top {
        display: flex;
        align-items: center;
        gap: 16rpx;
        margin-bottom: 12rpx;
        
        .receiver {
          font-size: 28rpx;
          font-weight: 600;
        }
        
        .phone {
          font-size: 26rpx;
          color: #666;
        }
        
        .default-tag {
          font-size: 20rpx;
          color: #E4393C;
          background: #FFECEB;
          padding: 4rpx 12rpx;
          border-radius: 4rpx;
        }
      }
      
      .address-detail {
        font-size: 26rpx;
        color: #666;
      }
    }
    
    .address-actions {
      display: flex;
      gap: 16rpx;
      border-top: 1px solid #F5F5F5;
      padding-top: 20rpx;
      
      .action-btn {
        font-size: 24rpx;
        color: #666;
        
        &.select-btn {
          color: #E4393C;
          font-weight: 600;
        }
        
        &.delete {
          color: #E4393C;
        }
      }
    }
  }
}

.empty {
  text-align: center;
  padding: 200rpx;
  color: #999;
}

.add-btn-wrap {
  position: fixed;
  bottom: 0;
  left: 0;
  right: 0;
  padding: 20rpx 24rpx;
  background: #fff;
  border-top: 1px solid #EEE;
  
  .btn-add {
    width: 100%;
    background: #E4393C;
    color: #fff;
    border-radius: 999rpx;
    padding: 24rpx;
    font-size: 28rpx;
  }
}
</style>
