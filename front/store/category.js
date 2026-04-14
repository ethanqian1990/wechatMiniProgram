// 分类 Store
import { defineStore } from 'pinia';
import { getCategories } from '@/api/category';

export const useCategoryStore = defineStore('category', {
  state: () => ({
    categories: [], // 分类列表
    currentCategory: '', // 当前选中的分类 id
    categoryMap: {} // id -> name 映射
  }),
  
  actions: {
    async fetchCategories() {
      try {
        const data = await getCategories();
        this.categories = data.list || data;
        // 建立 id -> name 映射
        this.categoryMap = {};
        this.categories.forEach(cat => {
          this.categoryMap[cat.id] = cat.name;
        });
        return this.categories;
      } catch (e) {
        console.error('获取分类失败', e);
        // 使用默认分类
        this.categories = [
          { id: 'c1', name: '农产品' },
          { id: 'c2', name: '生鲜水果' },
          { id: 'c3', name: '粮油副食' },
          { id: 'c4', name: '土特产' }
        ];
        return this.categories;
      }
    },
    
    setCategory(id) {
      this.currentCategory = id;
    }
  }
});
