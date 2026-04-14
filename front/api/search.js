// 搜索 API
import { get, del } from '@/utils/request';

export const getSearchSuggest = (keyword) => get('/search/suggest', { keyword });
export const getSearchHot = () => get('/search/hot');
export const getSearchHistory = () => get('/search/history');
export const clearSearchHistory = () => del('/search/history');
