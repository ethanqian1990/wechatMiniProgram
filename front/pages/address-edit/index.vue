<template>
  <view class="page">
    <view class="form-card">
      <!-- 收货人 -->
      <view class="form-group" :class="{ error: errors.receiver }">
        <text class="label">收货人</text>
        <input v-model="form.receiver" class="input" placeholder="请输入收货人姓名" />
        <text v-if="errors.receiver" class="error-msg">{{ errors.receiver }}</text>
      </view>
      
      <!-- 手机号 -->
      <view class="form-group" :class="{ error: errors.phone }">
        <text class="label">手机号</text>
        <input v-model="form.phone" class="input" type="number" placeholder="请输入手机号" maxlength="11" />
        <text v-if="errors.phone" class="error-msg">{{ errors.phone }}</text>
      </view>
      
      <!-- 省份 -->
      <view class="form-group" :class="{ error: errors.province }">
        <text class="label">省份</text>
        <picker :value="provinceIndex" :range="provinces" @change="onProvinceChange">
          <view class="picker-value">{{ form.province || '请选择省份' }}</view>
        </picker>
        <text v-if="errors.province" class="error-msg">{{ errors.province }}</text>
      </view>
      
      <!-- 城市 -->
      <view class="form-group" :class="{ error: errors.city }">
        <text class="label">城市</text>
        <picker :value="cityIndex" :range="cities" @change="onCityChange">
          <view class="picker-value">{{ form.city || '请选择城市' }}</view>
        </picker>
        <text v-if="errors.city" class="error-msg">{{ errors.city }}</text>
      </view>
      
      <!-- 区/县 -->
      <view class="form-group" :class="{ error: errors.district }">
        <text class="label">区/县</text>
        <picker :value="districtIndex" :range="districts" @change="onDistrictChange">
          <view class="picker-value">{{ form.district || '请选择区/县' }}</view>
        </picker>
        <text v-if="errors.district" class="error-msg">{{ errors.district }}</text>
      </view>
      
      <!-- 详细地址 -->
      <view class="form-group" :class="{ error: errors.detail }">
        <text class="label">详细地址</text>
        <textarea v-model="form.detail" class="textarea" placeholder="请输入详细地址" />
        <text v-if="errors.detail" class="error-msg">{{ errors.detail }}</text>
      </view>
      
      <!-- 设为默认 -->
      <view class="form-group checkbox-group">
        <checkbox v-model="form.is_default" />
        <text>设为默认地址</text>
      </view>
    </view>
    
    <!-- 保存按钮 -->
    <view class="submit-bar">
      <button class="btn-submit" :disabled="submitting" @click="onSubmit">保存</button>
    </view>
  </view>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue';
import { useAddressStore } from '@/store/address';

const addressStore = useAddressStore();

const form = ref({
  receiver: '',
  phone: '',
  province: '',
  city: '',
  district: '',
  detail: '',
  is_default: false
});

const errors = ref({});
const submitting = ref(false);
const isEdit = ref(false);
const editId = ref('');

const provinces = ['北京市', '上海市', '广东省', '江苏省', '浙江省', '四川省', '湖北省', '湖南省'];
const cities = ['市辖区', '市辖区', '广州市', '深圳市', '南京市', '杭州市', '成都市', '武汉市', '长沙市'];
const districts = ['东城区', '西城区', '朝阳区', '海淀区', '浦东新区', '黄浦区'];

const provinceIndex = ref(0);
const cityIndex = ref(0);
const districtIndex = ref(0);

onMounted(() => {
  const pages = getCurrentPages();
  const options = pages[pages.length - 1].options;
  
  if (options.address) {
    // 编辑模式
    const addr = JSON.parse(decodeURIComponent(options.address));
    isEdit.value = true;
    editId.value = addr.id;
    form.value = {
      receiver: addr.receiver,
      phone: addr.phone,
      province: addr.province,
      city: addr.city,
      district: addr.district,
      detail: addr.detail,
      is_default: addr.is_default
    };
  }
});

function onProvinceChange(e) {
  provinceIndex.value = e.detail.value;
  form.value.province = provinces[e.detail.value];
}

function onCityChange(e) {
  cityIndex.value = e.detail.value;
  form.value.city = cities[e.detail.value];
}

function onDistrictChange(e) {
  districtIndex.value = e.detail.value;
  form.value.district = districts[e.detail.value];
}

function validate() {
  errors.value = {};
  
  if (!form.value.receiver.trim()) {
    errors.value.receiver = '请填写收货人';
  }
  
  if (!form.value.phone.trim()) {
    errors.value.phone = '请填写手机号';
  } else if (!/^1\d{10}$/.test(form.value.phone)) {
    errors.value.phone = '请输入11位有效手机号';
  }
  
  if (!form.value.province) {
    errors.value.province = '请选择省份';
  }
  
  if (!form.value.city) {
    errors.value.city = '请选择城市';
  }
  
  if (!form.value.district) {
    errors.value.district = '请选择区/县';
  }
  
  if (!form.value.detail.trim()) {
    errors.value.detail = '请填写详细地址';
  }
  
  return Object.keys(errors.value).length === 0;
}

async function onSubmit() {
  if (!validate()) return;
  
  submitting.value = true;
  
  try {
    if (isEdit.value) {
      await addressStore.update(editId.value, form.value);
    } else {
      await addressStore.add(form.value);
    }
    
    uni.showToast({ title: '保存成功', icon: 'success' });
    setTimeout(() => {
      uni.navigateBack();
    }, 1500);
  } catch (e) {
    uni.showToast({ title: e.message || '保存失败', icon: 'none' });
  } finally {
    submitting.value = false;
  }
}
</script>

<style lang="scss" scoped>
.page {
  min-height: 100vh;
  background: #F5F5F5;
  padding-bottom: 120rpx;
}

.form-card {
  margin: 24rpx;
  background: #fff;
  border-radius: 16rpx;
  padding: 24rpx;
}

.form-group {
  margin-bottom: 24rpx;
  
  .label {
    font-size: 26rpx;
    color: #333;
    display: block;
    margin-bottom: 12rpx;
  }
  
  .input {
    width: 100%;
    padding: 20rpx;
    border: 1px solid #DDD;
    border-radius: 8rpx;
    font-size: 28rpx;
  }
  
  .textarea {
    width: 100%;
    padding: 20rpx;
    border: 1px solid #DDD;
    border-radius: 8rpx;
    font-size: 28rpx;
    min-height: 120rpx;
  }
  
  .picker-value {
    padding: 20rpx;
    border: 1px solid #DDD;
    border-radius: 8rpx;
    font-size: 28rpx;
    color: #333;
  }
  
  &.error {
    .input, .textarea, .picker-value {
      border-color: #E4393C;
    }
    
    .error-msg {
      font-size: 22rpx;
      color: #E4393C;
      margin-top: 8rpx;
    }
  }
}

.checkbox-group {
  display: flex;
  align-items: center;
  gap: 12rpx;
  
  text {
    font-size: 26rpx;
  }
}

.submit-bar {
  position: fixed;
  bottom: 0;
  left: 0;
  right: 0;
  padding: 20rpx 24rpx;
  background: #fff;
  border-top: 1px solid #EEE;
  
  .btn-submit {
    width: 100%;
    background: #E4393C;
    color: #fff;
    border-radius: 999rpx;
    padding: 24rpx;
    font-size: 28rpx;
    
    &:disabled {
      background: #CCC;
    }
  }
}
</style>
