-- ---------------------------------------------------------------
-- 衣橱 yichu · 0001 初始化表结构
-- 目标库：MySQL 8.0+，表字符集 utf8mb4，排序规则跟随数据库默认
-- 与 internal/model 下 GORM AutoMigrate 的结果保持一致
-- 执行：mysql --default-character-set=utf8mb4 -u yichu -p yichu < 0001_init.up.sql
-- ---------------------------------------------------------------

CREATE TABLE IF NOT EXISTS `users` (
  `id`            bigint unsigned NOT NULL AUTO_INCREMENT,
  `username`      varchar(64)     NOT NULL,
  `password_hash` varchar(100)    NOT NULL,
  `nickname`      varchar(50)     NOT NULL DEFAULT '',
  `created_at`    datetime(3)     NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_username` (`username`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='单账号管理员';

CREATE TABLE IF NOT EXISTS `categories` (
  `id`         bigint unsigned NOT NULL AUTO_INCREMENT,
  `parent_id`  bigint unsigned NOT NULL DEFAULT 0 COMMENT '0 表示一级分类',
  `name`       varchar(50)     NOT NULL,
  `sort`       bigint          NOT NULL DEFAULT 0,
  `is_system`  tinyint(1)      NOT NULL DEFAULT 0 COMMENT '种子数据标记，不可删除',
  `created_at` datetime(3)     NOT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_parent` (`parent_id`),
  KEY `idx_name` (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='衣物分类（两级）';

CREATE TABLE IF NOT EXISTS `items` (
  `id`               bigint unsigned NOT NULL AUTO_INCREMENT,
  `name`             varchar(100)    NOT NULL,
  `category_id`      bigint unsigned NOT NULL,
  `sub_category_id`  bigint unsigned DEFAULT NULL,
  `brand`            varchar(50)     DEFAULT NULL,
  `color`            json            DEFAULT NULL COMMENT '颜色 code 数组',
  `size`             varchar(20)     DEFAULT NULL,
  `season`           json            DEFAULT NULL COMMENT '季节 code 数组',
  `occasion`         json            DEFAULT NULL COMMENT '场合 code 数组',
  `purchase_date`    date            DEFAULT NULL,
  `purchase_price`   decimal(10,2)   DEFAULT NULL COMMENT 'hide_price 时服务端不返回',
  `purchase_channel` varchar(50)     DEFAULT NULL,
  `purchase_url`     varchar(500)    DEFAULT NULL,
  `status`           varchar(16)     NOT NULL DEFAULT 'wearing',
  `storage_location` varchar(50)     DEFAULT NULL,
  `remark`           varchar(500)    DEFAULT NULL,
  `wear_count`       int unsigned    NOT NULL DEFAULT 0,
  `last_worn_at`     datetime(3)     DEFAULT NULL,
  `created_at`       datetime(3)     NOT NULL,
  `updated_at`       datetime(3)     NOT NULL,
  `deleted_at`       datetime(3)     DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_name` (`name`),
  KEY `idx_category` (`category_id`),
  KEY `idx_sub_category` (`sub_category_id`),
  KEY `idx_brand` (`brand`),
  KEY `idx_purchase_date` (`purchase_date`),
  KEY `idx_status` (`status`),
  KEY `idx_created_at` (`created_at`),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='衣物主表（软删除）';

CREATE TABLE IF NOT EXISTS `item_images` (
  `id`         bigint unsigned NOT NULL AUTO_INCREMENT,
  `item_id`    bigint unsigned NOT NULL,
  `url`        varchar(500)    NOT NULL,
  `thumb_url`  varchar(500)    NOT NULL,
  `sort`       bigint          NOT NULL DEFAULT 0 COMMENT '升序，第一张为封面',
  `created_at` datetime(3)     NOT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_item_id` (`item_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='衣物图片';

CREATE TABLE IF NOT EXISTS `tags` (
  `id`         bigint unsigned NOT NULL AUTO_INCREMENT,
  `name`       varchar(50)     NOT NULL,
  `created_at` datetime(3)     NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_tag_name` (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='自定义标签';

CREATE TABLE IF NOT EXISTS `item_tags` (
  `item_id` bigint unsigned NOT NULL,
  `tag_id`  bigint unsigned NOT NULL,
  PRIMARY KEY (`item_id`, `tag_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='衣物-标签关联（many2many）';

CREATE TABLE IF NOT EXISTS `wear_logs` (
  `id`         bigint unsigned NOT NULL AUTO_INCREMENT,
  `item_id`    bigint unsigned NOT NULL,
  `worn_date`  date            NOT NULL,
  `occasion`   varchar(20)     DEFAULT NULL,
  `remark`     varchar(255)    DEFAULT NULL,
  `created_at` datetime(3)     NOT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_wear_item` (`item_id`),
  KEY `idx_worn_date` (`worn_date`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='穿着记录';

CREATE TABLE IF NOT EXISTS `settings` (
  `key`   varchar(64)  NOT NULL,
  `value` varchar(255) NOT NULL DEFAULT '',
  PRIMARY KEY (`key`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='运行时配置（覆盖 config.yaml）';
