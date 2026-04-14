// 首页 API
import { get } from '@/utils/request';

export const getHomeBanners = (regionCode) => get('/home/banners', { region_code: regionCode });
export const getHomeFeatured = (regionCode, categoryId) => {
  const data = { region_code: regionCode };
  if (categoryId) data.category_id = categoryId;
  return get('/home/featured', data);
};
export const getHomeProducts = (regionCode, categoryId) => {
  const data = { region_code: regionCode };
  if (categoryId) data.category_id = categoryId;
  return get('/home/products', data);
};
