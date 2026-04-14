// 分类 API
import { get } from '@/utils/request';

export const getCategories = () => get('/categories');
export const getCategory = (id) => get(`/categories/${id}`);
