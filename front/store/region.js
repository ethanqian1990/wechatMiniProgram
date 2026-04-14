// 区域 Store
import { defineStore } from 'pinia';
import { getRegions } from '@/api/region';

export const useRegionStore = defineStore('region', {
  state: () => ({
    currentRegion: uni.getStorageSync('region_code') || 'EC',
    regionList: [],
    regionMap: {} // id -> name 映射
  }),
  
  getters: {
    regionName: (state) => {
      const region = state.regionList.find(r => r.code === state.currentRegion);
      return region?.name || '华东';
    }
  },
  
  actions: {
    async fetchRegions() {
      try {
        const data = await getRegions();
        this.regionList = data.list || data;
        // 建立 code -> name 映射
        this.regionMap = {};
        this.regionList.forEach(r => {
          this.regionMap[r.code] = r.name;
        });
        return this.regionList;
      } catch (e) {
        console.error('获取区域列表失败', e);
        // 使用默认区域
        this.regionList = [
          { code: 'EC', name: '华东' },
          { code: 'NC', name: '华北' },
          { code: 'SC', name: '华南' },
          { code: 'NE', name: '东北' }
        ];
        return this.regionList;
      }
    },
    
    setRegion(code) {
      this.currentRegion = code;
      uni.setStorageSync('region_code', code);
    }
  }
});
