/*
Navicat MySQL Data Transfer

Source Server         : Golang
Source Server Version : 50505
Source Host           : 192.168.73.100:3306
Source Database       : pricedb-kosong

Target Server Type    : MYSQL
Target Server Version : 50505
File Encoding         : 65001

Date: 2017-12-22 07:26:57
*/

SET FOREIGN_KEY_CHECKS=0;

-- ----------------------------
-- Table structure for `product`
-- ----------------------------
DROP TABLE IF EXISTS `product`;
CREATE TABLE `product` (
  `store_id` varchar(3) NOT NULL DEFAULT '',
  `id` varchar(50) NOT NULL DEFAULT '',
  `name` varchar(50) DEFAULT NULL,
  `harga` float DEFAULT NULL,
  `url` varchar(100) DEFAULT NULL,
  `image` varchar(100) DEFAULT NULL,
  `update_date` datetime DEFAULT NULL,
  PRIMARY KEY (`store_id`,`id`),
  FULLTEXT KEY `idx` (`name`),
  CONSTRAINT `product_ibfk_1` FOREIGN KEY (`store_id`) REFERENCES `store` (`id`) ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- ----------------------------
-- Records of product
-- ----------------------------

-- ----------------------------
-- Table structure for `product_his`
-- ----------------------------
DROP TABLE IF EXISTS `product_his`;
CREATE TABLE `product_his` (
  `store_id` varchar(3) NOT NULL DEFAULT '',
  `id` varchar(50) NOT NULL DEFAULT '',
  `name` varchar(50) DEFAULT NULL,
  `harga` float DEFAULT NULL,
  `url` varchar(100) DEFAULT NULL,
  `image` varchar(100) DEFAULT NULL,
  `update_date` datetime NOT NULL DEFAULT '0000-00-00 00:00:00',
  PRIMARY KEY (`store_id`,`id`,`update_date`),
  FULLTEXT KEY `idx_his` (`name`),
  CONSTRAINT `product_his_ibfk_1` FOREIGN KEY (`store_id`) REFERENCES `store` (`id`) ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- ----------------------------
-- Records of product_his
-- ----------------------------

-- ----------------------------
-- Table structure for `store`
-- ----------------------------
DROP TABLE IF EXISTS `store`;
CREATE TABLE `store` (
  `id` varchar(3) NOT NULL DEFAULT '',
  `name` varchar(30) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- ----------------------------
-- Records of store
-- ----------------------------
INSERT INTO `store` VALUES ('alf', 'alfamart');
INSERT INTO `store` VALUES ('idm', 'indomaret');
