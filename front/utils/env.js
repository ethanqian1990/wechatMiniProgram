// 运行时环境配置（兼容 HBuilderX/小程序运行时，无需 .env 构建注入）
// 约定：优先读取本地存储 `api_base_url`，便于测试/灰度时动态切换。

export function getApiBaseURL() {
  const fromStorage = uni.getStorageSync('api_base_url');
  if (fromStorage) return String(fromStorage).replace(/\/+$/, '');

  // #ifdef H5
  // H5 可用同源代理；这里给默认值，生产建议配 Nginx 反代到 /api/v1
  return '/api/v1';
  // #endif

  // #ifdef MP-WEIXIN
  // 小程序端默认给一个占位域名，避免误用 localhost
  // 请在上线前通过存储或发布配置写入正确域名，例如：https://api.example.com/api/v1
  return 'https://api.example.com/api/v1';
  // #endif

  return '/api/v1';
}

