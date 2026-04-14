# 前端工程说明（UniApp / 微信小程序）

## 运行方式（推荐：HBuilderX）

本目录为 UniApp 工程（Vue3 + Pinia）。目前以 **HBuilderX** 作为主要开发/打包方式。

- **打开工程**：用 HBuilderX 打开 `front/`
- **运行到微信开发者工具**：HBuilderX 顶部菜单选择「运行」→「运行到小程序模拟器」→「微信开发者工具」
- **发行**：HBuilderX 顶部菜单选择「发行」→「小程序-微信」

## 后端地址配置（上线/联调必读）

请求基址不再写死本地 `localhost`，优先读取本地存储：

- key：`api_base_url`
- 值示例：`https://api.example.com/api/v1`

在微信开发者工具控制台或页面代码中设置一次即可：

```js
uni.setStorageSync('api_base_url', 'https://api.example.com/api/v1')
```

若未设置，将使用 `front/utils/env.js` 中的默认占位地址（上线前请务必配置）。

## 目录结构（节选）

- `pages/`：页面（首页/分类/搜索/购物车/订单/地址等）
- `store/`：Pinia 状态
- `api/`：接口封装（统一走 `utils/request.js`）
- `utils/`：请求封装与支付轮询兜底等工具

## 支付链路说明（符合 PRD 3.4）

- 支付调起后会进入“支付结果确认中”状态，并轮询：
  - `GET /api/v1/pay/status/{order_id}`
- 超时（默认 30s）会提示用户到订单详情查看，避免回调延迟导致的“支付成功但页面不刷新”。

