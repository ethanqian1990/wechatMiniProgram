// 区域 API
import { get } from '@/utils/request';

export const getRegions = () => get('/regions');
export const getRegion = (code) => get(`/regions/${code}`);
