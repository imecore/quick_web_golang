CREATE DATABASE IF NOT EXISTS quick_web;

USE quick_web;

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for user
-- ----------------------------
DROP TABLE IF EXISTS `user`;
CREATE TABLE `user`
(
    `id`            int(11) unsigned NOT NULL AUTO_INCREMENT,
    `username`      varchar(36) NOT NULL DEFAULT '',
    `password`      char(36) NOT NULL DEFAULT '',
    `salt`          char(8) NOT NULL DEFAULT '',
    `created_at`    datetime DEFAULT CURRENT_TIMESTAMP,
    `updated_at`    datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`) USING BTREE,
    UNIQUE KEY `username` (`username`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci
  ROW_FORMAT = DYNAMIC;

SET FOREIGN_KEY_CHECKS = 1;