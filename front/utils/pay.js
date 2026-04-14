import { getPayStatus } from '@/api/pay';

function sleep(ms) {
  return new Promise((r) => setTimeout(r, ms));
}

// 支付后兜底：轮询订单支付状态，解决回调延迟/网络异常导致的“支付成功但页面未刷新”
export async function pollPayStatus(orderId, opts = {}) {
  const {
    maxMs = 30000,
    steps = [1000, 2000, 4000, 8000, 15000], // 递增间隔，最多 30s
  } = opts;

  const start = Date.now();
  let i = 0;
  while (Date.now() - start < maxMs) {
    try {
      const data = await getPayStatus(orderId);
      const status = data?.status;
      if (status && status !== 'PENDING_PAY') return { ok: true, status };
    } catch (e) {
      // 忽略中间错误，继续轮询
    }
    const wait = steps[Math.min(i, steps.length - 1)];
    i += 1;
    await sleep(wait);
  }
  return { ok: false, status: 'UNKNOWN' };
}

