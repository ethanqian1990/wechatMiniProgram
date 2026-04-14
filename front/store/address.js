// 地址 Store
import { defineStore } from 'pinia';
import { getAddressList, addAddress, updateAddress, deleteAddress, setDefaultAddress } from '@/api/address';

export const useAddressStore = defineStore('address', {
  state: () => ({
    addresses: [],
    selectedAddressId: uni.getStorageSync('selected_address_id') || ''
  }),
  
  getters: {
    defaultAddress: (state) => state.addresses.find(a => a.is_default) || state.addresses[0],
    
    selectedAddress: (state) => state.addresses.find(a => a.id === state.selectedAddressId) || state.addresses[0]
  },
  
  actions: {
    async fetchAddresses() {
      try {
        const data = await getAddressList();
        this.addresses = data.list || data;
        return this.addresses;
      } catch (e) {
        console.error('获取地址列表失败', e);
        return [];
      }
    },
    
    async add(addrData) {
      try {
        const data = await addAddress(addrData);
        await this.fetchAddresses();
        return data;
      } catch (e) {
        uni.showToast({ title: e.message || '添加失败', icon: 'none' });
        throw e;
      }
    },
    
    async update(id, addrData) {
      try {
        await updateAddress(id, addrData);
        await this.fetchAddresses();
      } catch (e) {
        uni.showToast({ title: e.message || '更新失败', icon: 'none' });
        throw e;
      }
    },
    
    async remove(id) {
      try {
        await deleteAddress(id);
        await this.fetchAddresses();
      } catch (e) {
        uni.showToast({ title: e.message || '删除失败', icon: 'none' });
        throw e;
      }
    },
    
    async setDefault(id) {
      try {
        await setDefaultAddress(id);
        await this.fetchAddresses();
      } catch (e) {
        uni.showToast({ title: e.message || '设置失败', icon: 'none' });
        throw e;
      }
    },
    
    selectAddress(id) {
      this.selectedAddressId = id;
      uni.setStorageSync('selected_address_id', id);
    }
  }
});
