-- 增量迁移：为已存在 orders 表补齐物流字段与支付幂等交易号
-- 注意：如果你是全新初始化数据库，001_init.sql 已包含这些字段，可跳过此迁移

ALTER TABLE orders
  ADD COLUMN IF NOT EXISTS express_company VARCHAR(50) COMMENT '快递公司' AFTER region_code,
  ADD COLUMN IF NOT EXISTS express_no VARCHAR(64) COMMENT '快递单号' AFTER express_company,
  ADD COLUMN IF NOT EXISTS paid_transaction_id VARCHAR(64) COMMENT '支付交易号(幂等)' AFTER express_no;

-- paid_transaction_id 唯一索引（MySQL 允许多条 NULL，不影响未支付订单）
CREATE UNIQUE INDEX IF NOT EXISTS uk_paid_tx ON orders(paid_transaction_id);

