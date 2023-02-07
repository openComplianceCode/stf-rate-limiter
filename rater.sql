# ************************************************************
# Sequel Pro SQL dump
# Version 5446
#
# https://www.sequelpro.com/
# https://github.com/sequelpro/sequelpro
#
# Host: 127.0.0.1 (MySQL 8.0.27)
# Database: osca
# Generation Time: 2022-12-15 07:40:44 +0000
# ************************************************************


/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8 */;
SET NAMES utf8mb4;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;


# Dump of table users
# ------------------------------------------------------------

DROP TABLE IF EXISTS `users`;

CREATE TABLE `users` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `create_time` timestamp NOT NULL,
  `update_time` timestamp NOT NULL,
  `role` enum('everyone','gold','diamond','root') COLLATE utf8mb4_bin NOT NULL DEFAULT 'everyone',
  `gitee_id` varchar(255) COLLATE utf8mb4_bin DEFAULT NULL,
  `gitee_login` varchar(255) COLLATE utf8mb4_bin DEFAULT NULL,
  `gitee_name` varchar(255) COLLATE utf8mb4_bin DEFAULT NULL,
  `gitee_email` varchar(255) COLLATE utf8mb4_bin DEFAULT NULL,
  `gitee_avatar_url` varchar(255) COLLATE utf8mb4_bin DEFAULT NULL,
  `github_id` varchar(255) COLLATE utf8mb4_bin DEFAULT NULL,
  `github_login` varchar(255) COLLATE utf8mb4_bin DEFAULT NULL,
  `github_name` varchar(255) COLLATE utf8mb4_bin DEFAULT NULL,
  `github_email` varchar(255) COLLATE utf8mb4_bin DEFAULT NULL,
  `github_avatar_url` varchar(255) COLLATE utf8mb4_bin DEFAULT NULL,
  `api_token` varchar(255) COLLATE utf8mb4_bin DEFAULT NULL,
  `api_token_generate_time` timestamp NULL DEFAULT NULL,
  `first_name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT NULL,
  `last_name` varchar(255) COLLATE utf8mb4_bin DEFAULT NULL,
  `email_address` varchar(255) COLLATE utf8mb4_bin DEFAULT NULL,
  `company` varchar(255) COLLATE utf8mb4_bin DEFAULT NULL,
  `city` varchar(255) COLLATE utf8mb4_bin DEFAULT NULL,
  `country` varchar(255) COLLATE utf8mb4_bin DEFAULT NULL,
  `postal_code` varchar(255) COLLATE utf8mb4_bin DEFAULT NULL,
  `address` varchar(512) COLLATE utf8mb4_bin DEFAULT NULL,
  `about_me` varchar(512) COLLATE utf8mb4_bin DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;




/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;
/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
