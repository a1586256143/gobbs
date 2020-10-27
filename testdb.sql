/*
 Navicat Premium Data Transfer

 Source Server         : localhost
 Source Server Type    : MySQL
 Source Server Version : 80020
 Source Host           : 127.0.0.1:3306
 Source Schema         : testdb

 Target Server Type    : MySQL
 Target Server Version : 80020
 File Encoding         : 65001

 Date: 27/10/2020 14:47:31
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for article
-- ----------------------------
DROP TABLE IF EXISTS `article`;
CREATE TABLE `article` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `uid` int DEFAULT NULL,
  `cid` int DEFAULT NULL,
  `title` varchar(40) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL,
  `content` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci,
  `create_time` int DEFAULT NULL,
  `status` tinyint(1) DEFAULT '0',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=24 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- ----------------------------
-- Records of article
-- ----------------------------
BEGIN;
INSERT INTO `article` VALUES (18, 6, 5, '如果让一个人爱上自己', '如果让一个人爱上自己，请加WX 15070091894', 1600841437, 0);
INSERT INTO `article` VALUES (19, 6, 5, 'aaa', '<img src=\"/upload/20201023/20032100222878.png\" alt=\"undefined\">', 1603445248, 0);
INSERT INTO `article` VALUES (20, 6, 5, '招聘启示', '<p style=\"text-align: center;\"><b>《招聘》</b></p><p style=\"text-align: left;\">招聘客服岗位需求如下：</p><p style=\"text-align: left;\">1. 身高1.5米以上</p><p style=\"text-align: left;\">2. 面容较好</p><p style=\"text-align: left;\"><img src=\"http://127.0.0.1:9999/static/layui/images/face/1.gif\" alt=\"[嘻嘻]\"><br></p>', 1603448102, 0);
INSERT INTO `article` VALUES (21, 6, 7, '无聊进群聊天啊', '<p>无聊进群聊天啊无聊进群聊天啊无聊进群聊天啊</p><p><img src=\"/upload/20201023/20160701160326389.png\" alt=\"undefined\"><br></p>', 1603448250, 0);
INSERT INTO `article` VALUES (22, 6, 5, 'dddddd', 'lalaa<img src=\"/upload/20201023/20070615533973.png\" alt=\"undefined\">', 1603448319, 0);
INSERT INTO `article` VALUES (23, 7, 5, '来呀，来呀', '22222', 1603694816, 0);
COMMIT;

-- ----------------------------
-- Table structure for category
-- ----------------------------
DROP TABLE IF EXISTS `category`;
CREATE TABLE `category` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `title` varchar(40) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=8 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- ----------------------------
-- Records of category
-- ----------------------------
BEGIN;
INSERT INTO `category` VALUES (5, '个人情感');
INSERT INTO `category` VALUES (6, '饭后时光');
INSERT INTO `category` VALUES (7, '其它');
COMMIT;

-- ----------------------------
-- Table structure for comments
-- ----------------------------
DROP TABLE IF EXISTS `comments`;
CREATE TABLE `comments` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `uid` int unsigned DEFAULT '0',
  `aid` int unsigned DEFAULT '0',
  `content` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci,
  `create_time` int unsigned DEFAULT '0',
  `status` tinyint unsigned DEFAULT '0',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=7 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- ----------------------------
-- Records of comments
-- ----------------------------
BEGIN;
INSERT INTO `comments` VALUES (1, 6, 18, '路过', 1600842030, 0);
INSERT INTO `comments` VALUES (2, 6, 18, '我是谁，我在哪', 1600842159, 1);
INSERT INTO `comments` VALUES (3, 6, 18, '沙发，沙发', 1600842178, 1);
INSERT INTO `comments` VALUES (4, 6, 18, '大家好，我是张三', 1603440451, 1);
INSERT INTO `comments` VALUES (5, 6, 18, '说的都什么啊', 1603441252, 1);
INSERT INTO `comments` VALUES (6, 6, 18, 'mmmm', 1603544202, 1);
COMMIT;

-- ----------------------------
-- Table structure for friends
-- ----------------------------
DROP TABLE IF EXISTS `friends`;
CREATE TABLE `friends` (
  `id` int NOT NULL AUTO_INCREMENT,
  `uid` int DEFAULT '0' COMMENT '请求发起人',
  `pull_uid` int DEFAULT '0' COMMENT '接收人',
  `remark` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL COMMENT '请求附加信息',
  `status` tinyint unsigned DEFAULT '0' COMMENT '状态 0 待同意 1 已同意 2 已拒绝',
  `create_time` int unsigned DEFAULT '0' COMMENT '发起时间',
  `update_time` int unsigned DEFAULT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci COMMENT='好友表';

-- ----------------------------
-- Records of friends
-- ----------------------------
BEGIN;
INSERT INTO `friends` VALUES (1, 7, 6, '你好，我是colin', 1, 1603680544, 1603684850);
COMMIT;

-- ----------------------------
-- Table structure for gousers
-- ----------------------------
DROP TABLE IF EXISTS `gousers`;
CREATE TABLE `gousers` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(64) DEFAULT NULL,
  `password` varchar(40) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL,
  `telephone` varchar(12) DEFAULT '',
  `avatar` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '头像',
  `last_time` int DEFAULT NULL,
  `create_time` int DEFAULT NULL,
  `views` int DEFAULT '0' COMMENT '总浏览数',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=8 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- ----------------------------
-- Records of gousers
-- ----------------------------
BEGIN;
INSERT INTO `gousers` VALUES (6, 'admin', 'd033e22ae348aeb5660fc2140aec35850c4da997', '17603017302', NULL, 1603696381, 1600833490, 41);
INSERT INTO `gousers` VALUES (7, 'colin', '7c4a8d09ca3762af61e59520943dc26494f8941b', '17603017302', NULL, 1603694642, 1603676720, 2);
COMMIT;

-- ----------------------------
-- Table structure for stat
-- ----------------------------
DROP TABLE IF EXISTS `stat`;
CREATE TABLE `stat` (
  `id` int NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `uid` int DEFAULT NULL COMMENT '用户ID',
  `day` int DEFAULT NULL COMMENT '时间',
  `views` int DEFAULT '0' COMMENT '次数',
  PRIMARY KEY (`id`),
  KEY `uid` (`uid`,`day`)
) ENGINE=InnoDB AUTO_INCREMENT=9 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci COMMENT='访问统计表';

-- ----------------------------
-- Records of stat
-- ----------------------------
BEGIN;
INSERT INTO `stat` VALUES (4, 6, 20201024, 45);
INSERT INTO `stat` VALUES (5, 0, 20201024, 2);
INSERT INTO `stat` VALUES (6, 0, 20201026, 5);
INSERT INTO `stat` VALUES (7, 7, 20201026, 2);
INSERT INTO `stat` VALUES (8, 6, 20201026, 1);
COMMIT;

SET FOREIGN_KEY_CHECKS = 1;
