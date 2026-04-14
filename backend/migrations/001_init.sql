-- 购物商城小程序数据库初始化脚本
-- 版本: 1.0
-- 日期: 2026-04-14

-- 创建数据库
CREATE DATABASE IF NOT EXISTS wechat_mall DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
USE wechat_mall;

-- 用户表
CREATE TABLE IF NOT EXISTS users (
    id              VARCHAR(36) PRIMARY KEY COMMENT '用户ID',
    openid          VARCHAR(128) UNIQUE COMMENT '微信openid',
    nickname        VARCHAR(100) COMMENT '用户昵称',
    avatar          VARCHAR(500) COMMENT '头像URL',
    phone           VARCHAR(20) COMMENT '手机号',
    status          TINYINT DEFAULT 1 COMMENT '状态: 1正常 0禁用',
    last_login_at   DATETIME COMMENT '最后登录时间',
    created_at      DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at      DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at      DATETIME DEFAULT NULL,
    INDEX idx_openid (openid),
    INDEX idx_phone (phone)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户表';

-- 区域表
CREATE TABLE IF NOT EXISTS regions (
    code            VARCHAR(20) PRIMARY KEY COMMENT '区域编码如EC、NE',
    id              VARCHAR(36) COMMENT '内部ID(可选)',
    region_name     VARCHAR(50) NOT NULL COMMENT '区域名称',
    is_enabled      TINYINT DEFAULT 1 COMMENT '是否启用',
    sort_order      INT DEFAULT 0 COMMENT '排序',
    created_at      DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at      DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_enabled (is_enabled),
    INDEX idx_sort (sort_order)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='区域表';

-- 初始化区域数据
INSERT INTO regions (code, region_name, sort_order) VALUES
('EC', '华东', 1),
('NC', '华北', 2),
('CC', '华中', 3),
('SC', '华南', 4),
('SW', '西南', 5),
('NW', '西北', 6),
('NE', '东北', 7);

-- 商品分类表
CREATE TABLE IF NOT EXISTS categories (
    id              VARCHAR(36) PRIMARY KEY,
    name            VARCHAR(50) NOT NULL COMMENT '分类名称',
    parent_id       VARCHAR(36) DEFAULT NULL COMMENT '父分类ID',
    level           TINYINT DEFAULT 1 COMMENT '层级:1一级 2二级',
    icon            VARCHAR(500) COMMENT '图标',
    sort_order      INT DEFAULT 0,
    is_enabled      TINYINT DEFAULT 1,
    created_at      DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at      DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_parent (parent_id),
    INDEX idx_level (level),
    INDEX idx_enabled (is_enabled)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='商品分类表';

-- 商品表
CREATE TABLE IF NOT EXISTS products (
    id              VARCHAR(36) PRIMARY KEY,
    name            VARCHAR(200) NOT NULL COMMENT '商品名称',
    description     TEXT COMMENT '商品描述',
    main_image      VARCHAR(500) COMMENT '主图URL',
    images          JSON COMMENT '图片列表',
    category_id     VARCHAR(36) COMMENT '分类ID',
    price           DECIMAL(10,2) NOT NULL DEFAULT 0 COMMENT '售价',
    original_price  DECIMAL(10,2) DEFAULT NULL COMMENT '划线价',
    stock           INT DEFAULT 0 COMMENT '总库存',
    sales_count     INT DEFAULT 0 COMMENT '销量',
    is_on_shelf     TINYINT DEFAULT 1 COMMENT '是否上架:1上架 0下架',
    show_on_home    TINYINT DEFAULT 1 COMMENT '是否在首页展示',
    available_regions JSON COMMENT '可售区域列表',
    tags            JSON COMMENT '标签',
    promo_text      VARCHAR(200) COMMENT '促销文案',
    low_stock_threshold INT DEFAULT 10 COMMENT '低库存预警阈值',
    created_at      DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at      DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at      DATETIME DEFAULT NULL,
    INDEX idx_category (category_id),
    INDEX idx_on_shelf (is_on_shelf),
    INDEX idx_show_home (show_on_home),
    INDEX idx_sales (sales_count DESC)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='商品表';

-- 商品SKU表
CREATE TABLE IF NOT EXISTS product_skus (
    id              VARCHAR(36) PRIMARY KEY,
    product_id      VARCHAR(36) NOT NULL,
    sku_name        VARCHAR(200) NOT NULL COMMENT 'SKU名称',
    specs           JSON NOT NULL COMMENT '规格详情',
    price           DECIMAL(10,2) NOT NULL COMMENT '售价',
    original_price  DECIMAL(10,2) DEFAULT NULL COMMENT '划线价',
    stock           INT DEFAULT 0 COMMENT '库存',
    image           VARCHAR(500) COMMENT 'SKU图片',
    is_enabled      TINYINT DEFAULT 1,
    created_at      DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at      DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_product (product_id),
    INDEX idx_enabled (is_enabled),
    FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='商品SKU表';

-- 购物车表
CREATE TABLE IF NOT EXISTS carts (
    id              VARCHAR(36) PRIMARY KEY,
    user_id         VARCHAR(36) NOT NULL,
    product_id      VARCHAR(36) NOT NULL,
    sku_id          VARCHAR(36) DEFAULT NULL,
    quantity        INT NOT NULL DEFAULT 1,
    created_at      DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at      DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_user (user_id),
    INDEX idx_user_product (user_id, product_id, sku_id),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='购物车表';

-- 收货地址表
CREATE TABLE IF NOT EXISTS addresses (
    id              VARCHAR(36) PRIMARY KEY,
    user_id         VARCHAR(36) NOT NULL,
    receiver        VARCHAR(50) NOT NULL COMMENT '收货人',
    phone           VARCHAR(20) NOT NULL COMMENT '手机号',
    province        VARCHAR(50) NOT NULL COMMENT '省份',
    city            VARCHAR(50) NOT NULL COMMENT '城市',
    district        VARCHAR(50) NOT NULL COMMENT '区县',
    detail          VARCHAR(200) NOT NULL COMMENT '详细地址',
    is_default      TINYINT DEFAULT 0 COMMENT '是否默认:1默认 0非默认',
    created_at      DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at      DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_user (user_id),
    INDEX idx_default (user_id, is_default),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='收货地址表';

-- 订单表
CREATE TABLE IF NOT EXISTS orders (
    id                  VARCHAR(36) PRIMARY KEY,
    order_no            VARCHAR(64) UNIQUE NOT NULL COMMENT '订单号',
    user_id             VARCHAR(36) NOT NULL,
    total_amount        DECIMAL(10,2) NOT NULL COMMENT '商品总额',
    discount_amount     DECIMAL(10,2) DEFAULT 0 COMMENT '优惠金额',
    freight_amount      DECIMAL(10,2) DEFAULT 0 COMMENT '运费',
    final_amount        DECIMAL(10,2) NOT NULL COMMENT '实付金额',
    status              VARCHAR(20) NOT NULL DEFAULT 'PENDING_PAY' COMMENT '订单状态',
    receiver            VARCHAR(50) COMMENT '收货人',
    phone               VARCHAR(20) COMMENT '收货手机',
    province            VARCHAR(50) COMMENT '省份',
    city                VARCHAR(50) COMMENT '城市',
    district            VARCHAR(50) COMMENT '区县',
    detail              VARCHAR(200) COMMENT '详细地址',
    pay_at              DATETIME COMMENT '支付时间',
    ship_at             DATETIME COMMENT '发货时间',
    receive_at          DATETIME COMMENT '收货时间',
    cancel_at           DATETIME COMMENT '取消时间',
    cancel_reason       VARCHAR(200) COMMENT '取消原因',
    client_order_token  VARCHAR(64) COMMENT '客户端订单token',
    region_code         VARCHAR(20) COMMENT '下单区域',
    created_at          DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at          DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_user (user_id),
    INDEX idx_status (status),
    INDEX idx_order_no (order_no),
    INDEX idx_token (client_order_token),
    INDEX idx_created (created_at),
    INDEX idx_region (region_code),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='订单表';

-- 订单商品表
CREATE TABLE IF NOT EXISTS order_items (
    id              VARCHAR(36) PRIMARY KEY,
    order_id        VARCHAR(36) NOT NULL,
    product_id      VARCHAR(36) NOT NULL,
    sku_id          VARCHAR(36) DEFAULT NULL,
    product_name    VARCHAR(200) NOT NULL,
    product_image   VARCHAR(500),
    sku_name        VARCHAR(200) COMMENT 'SKU规格',
    price           DECIMAL(10,2) NOT NULL COMMENT '购买时单价',
    quantity        INT NOT NULL,
    subtotal        DECIMAL(10,2) NOT NULL COMMENT '小计',
    created_at      DATETIME DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_order (order_id),
    INDEX idx_product (product_id),
    FOREIGN KEY (order_id) REFERENCES orders(id) ON DELETE CASCADE,
    FOREIGN KEY (product_id) REFERENCES products(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='订单商品表';

-- 库存冻结/预占表
CREATE TABLE IF NOT EXISTS stock_reservations (
    id              VARCHAR(36) PRIMARY KEY,
    order_id        VARCHAR(36) NOT NULL COMMENT '订单ID',
    user_id         VARCHAR(36) NOT NULL COMMENT '用户ID',
    product_id      VARCHAR(36) NOT NULL,
    sku_id          VARCHAR(36) NOT NULL,
    quantity        INT NOT NULL COMMENT '冻结数量',
    status          VARCHAR(20) NOT NULL DEFAULT 'FROZEN' COMMENT 'FROZEN|CONSUMED|RELEASED',
    expire_at       DATETIME NOT NULL COMMENT '冻结过期时间',
    created_at      DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at      DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY uk_order_sku (order_id, sku_id),
    INDEX idx_order (order_id),
    INDEX idx_user (user_id),
    INDEX idx_status_expire (status, expire_at),
    FOREIGN KEY (order_id) REFERENCES orders(id) ON DELETE CASCADE,
    FOREIGN KEY (sku_id) REFERENCES product_skus(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='库存冻结/预占表';

-- 首页运营配置表
CREATE TABLE IF NOT EXISTS home_configs (
    id              VARCHAR(36) PRIMARY KEY,
    region_code     VARCHAR(20) NOT NULL COMMENT '区域编码',
    config_type     VARCHAR(20) NOT NULL COMMENT '配置类型:banner/featured/list',
    config_key      VARCHAR(50) COMMENT '配置键',
    config_value    JSON NOT NULL COMMENT '配置值',
    sort_order      INT DEFAULT 0,
    created_at      DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at      DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY uk_region_type_key (region_code, config_type, config_key),
    INDEX idx_region (region_code),
    INDEX idx_type (config_type)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='首页运营配置表';

-- 搜索历史表
CREATE TABLE IF NOT EXISTS search_history (
    id              VARCHAR(36) PRIMARY KEY,
    user_id         VARCHAR(36) NOT NULL,
    keyword         VARCHAR(100) NOT NULL,
    created_at      DATETIME DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_user_key (user_id, keyword),
    INDEX idx_created (created_at),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='搜索历史表';

-- 初始化测试数据
INSERT INTO categories (id, name, level, sort_order) VALUES
('c1', '农产品', 1, 1),
('c2', '生鲜水果', 1, 2),
('c3', '粮油米面', 1, 3),
('c4', '零食饮料', 1, 4);

INSERT INTO products (id, name, description, main_image, category_id, price, original_price, stock, sales_count, is_on_shelf, show_on_home, available_regions, tags, promo_text) VALUES
('p1', '红富士苹果', '新鲜红富士苹果脆甜多汁', 'https://via.placeholder.com/400', 'c2', 59.90, 79.90, 100, 150, 1, 1, '["EC","NC","SC"]', '["新鲜","脆甜"]', '限时优惠'),
('p2', '东北大米', '优质东北大米5kg装', 'https://via.placeholder.com/400', 'c3', 39.90, 49.90, 200, 80, 1, 1, '["NE","NC"]', '["优质"]', ''),
('p3', '土鸡蛋', '农家土鸡蛋30枚', 'https://via.placeholder.com/400', 'c2', 45.00, 55.00, 50, 200, 1, 1, '["EC","NE"]', '["农家","新鲜"]', '热销'),
('p4', '新疆葡萄干', '新疆特产葡萄干500g', 'https://via.placeholder.com/400', 'c4', 25.00, 30.00, 300, 50, 1, 1, '["NW","SW"]', '["特产"]', '');

INSERT INTO product_skus (id, product_id, sku_name, specs, price, original_price, stock) VALUES
('sku1', 'p1', '5斤装', '{"weight": "5斤"}', 59.90, 79.90, 50),
('sku2', 'p1', '10斤装', '{"weight": "10斤"}', 99.90, 129.90, 50),
('sku3', 'p2', '5kg装', '{"weight": "5kg"}', 39.90, 49.90, 100),
('sku4', 'p2', '10kg装', '{"weight": "10kg"}', 69.90, 89.90, 100),
('sku5', 'p3', '30枚装', '{"count": 30}', 45.00, 55.00, 50),
('sku6', 'p4', '500g装', '{"weight": "500g"}', 25.00, 30.00, 300);
