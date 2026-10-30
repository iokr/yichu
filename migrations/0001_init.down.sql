-- ---------------------------------------------------------------
-- 衣橱 yichu · 0001 回滚：删除全部业务表
-- 危险操作：会清空所有数据，仅用于开发环境重建
-- ---------------------------------------------------------------

DROP TABLE IF EXISTS `item_tags`;
DROP TABLE IF EXISTS `wear_logs`;
DROP TABLE IF EXISTS `item_images`;
DROP TABLE IF EXISTS `items`;
DROP TABLE IF EXISTS `tags`;
DROP TABLE IF EXISTS `categories`;
DROP TABLE IF EXISTS `settings`;
DROP TABLE IF EXISTS `users`;
