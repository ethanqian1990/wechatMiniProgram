// 订单 Store
import { defineStore } from 'pinia';
import { getOrders, getOrder, createOrder, cancelOrder, confirmOrder } from '@/api/order';

export const useOrderStore = defineStore('order', {
  state: () => ({
    orders: [],
    currentOrder: null,
    loading: false
  }),
  
  getters: {
    orderMap: (state) => {
      const map = {};
      state.orders.forEach(o => { map[o.id] = o; });
      return map;
    }
  },
  
  actions: {
    async fetchOrders(params = {}) {
      this.loading = true;
      try {
        const data = await getOrders(params);
        this.orders = data.list || data || [];
        return this.orders;
      } catch (e) {
        console.error('获取订单列表失败', e);
        return [];
      } finally {
        this.loading = false;
      }
    },
    
    async fetchOrder(id) {
      try {
        const data = await getOrder(id);
        this.currentOrder = data;
        return data;
      } catch (e) {
        console.error('获取订单详情失败', e);
        throw e;
      }
    },
    
    async create(orderData) {
      try {
        const data = await createOrder(orderData);
        return data;
      } catch (e) {
        if (e.code === 400 && e.message === '库存不足') {
          uni.showToast({ title: '库存不足，请调整数量', icon: 'none' });
        }
        throw e;
      }
    },
    
    async cancel(id) {
      try {
        await cancelOrder(id);
        await this.fetchOrders();
      } catch (e) {
        uni.showToast({ title: e.message || '取消失败', icon: 'none' });
        throw e;
      }
    },
    
    async confirm(id) {
      try {
        await confirmOrder(id);
        await this.fetchOrders();
      } catch (e) {
        uni.showToast({ title: e.message || '确认失败', icon: 'none' });
        throw e;
      }
    }
  }
});
