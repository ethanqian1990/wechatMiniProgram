// 购物车 Store
import { defineStore } from 'pinia';
import { getCartList, addToCart, updateCartItem, removeCartItem, clearCart } from '@/api/cart';

export const useCartStore = defineStore('cart', {
  state: () => ({
    items: [], // 购物车商品列表
    loading: false
  }),
  
  getters: {
    totalCount: (state) => state.items.reduce((sum, item) => sum + item.quantity, 0),
    
    checkedItems: (state) => state.items.filter(item => item.checked),
    
    totalPrice: (state) => {
      return state.items
        .filter(item => item.checked)
        .reduce((sum, item) => {
          const price = item.sku?.price || item.product?.price || 0;
          return sum + price * item.quantity;
        }, 0);
    },
    
    hasLoginItems: (state) => state.items.some(item => !item.valid),
    
    allChecked: (state) => state.items.length > 0 && state.items.every(item => item.checked)
  },
  
  actions: {
    async fetchCart() {
      this.loading = true;
      try {
        const data = await getCartList();
        this.items = (data.list || data).map(item => ({
          ...item,
          checked: true
        }));
      } catch (e) {
        console.error('获取购物车失败', e);
      } finally {
        this.loading = false;
      }
    },
    
    async addItem(productId, skuId, quantity = 1) {
      try {
        await addToCart({ product_id: productId, sku_id: skuId, quantity });
        uni.showToast({ title: '已加入购物车', icon: 'success' });
        // 刷新购物车
        await this.fetchCart();
      } catch (e) {
        uni.showToast({ title: e.message || '添加失败', icon: 'none' });
        throw e;
      }
    },
    
    async updateItemQuantity(id, quantity) {
      try {
        await updateCartItem(id, { quantity });
        const item = this.items.find(i => i.id === id);
        if (item) item.quantity = quantity;
      } catch (e) {
        uni.showToast({ title: '更新失败', icon: 'none' });
        throw e;
      }
    },
    
    async removeItem(id) {
      try {
        await removeCartItem(id);
        this.items = this.items.filter(i => i.id !== id);
      } catch (e) {
        uni.showToast({ title: '删除失败', icon: 'none' });
        throw e;
      }
    },
    
    async clearAll() {
      try {
        await clearCart();
        this.items = [];
      } catch (e) {
        uni.showToast({ title: '清空失败', icon: 'none' });
        throw e;
      }
    },
    
    toggleItem(id) {
      const item = this.items.find(i => i.id === id);
      if (item) item.checked = !item.checked;
    },
    
    toggleAll(checked) {
      this.items.forEach(item => item.checked = checked);
    }
  }
});
