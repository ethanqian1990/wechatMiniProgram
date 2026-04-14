// 请求封装
const BASE_URL = 'http://localhost:8080/api/v1';

const request = (options) => {
  return new Promise((resolve, reject) => {
    const token = uni.getStorageSync('token');
    
    uni.request({
      url: BASE_URL + options.url,
      method: options.method || 'GET',
      header: {
        'Authorization': token ? `Bearer ${token}` : '',
        'Content-Type': 'application/json'
      },
      data: options.data,
      success: (res) => {
        if (res.statusCode === 200) {
          if (res.data.code === 200) {
            resolve(res.data.data);
          } else if (res.data.code === 401) {
            // Token 过期
            uni.removeStorageSync('token');
            uni.removeStorageSync('userInfo');
            uni.showToast({ title: '登录已过期', icon: 'none' });
            reject({ code: 401, message: '未登录' });
          } else if (res.data.code === 400 && res.data.message === '库存不足') {
            reject({ code: 400, message: '库存不足' });
          } else {
            reject(res.data);
          }
        } else if (res.statusCode === 401) {
          uni.removeStorageSync('token');
          reject({ code: 401, message: '未登录' });
        } else {
          reject(res.data);
        }
      },
      fail: (err) => {
        uni.showToast({ title: '网络请求失败', icon: 'none' });
        reject(err);
      }
    });
  });
};

// GET 请求
export const get = (url, data) => request({ url, method: 'GET', data });

// POST 请求
export const post = (url, data) => request({ url, method: 'POST', data });

// PUT 请求
export const put = (url, data) => request({ url, method: 'PUT', data });

// DELETE 请求
export const del = (url, data) => request({ url, method: 'DELETE', data });

export default { get, post, put, del };
