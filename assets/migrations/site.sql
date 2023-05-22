SET NAMES utf8;
SET time_zone = '+00:00';
SET foreign_key_checks = 0;
SET sql_mode = 'NO_AUTO_VALUE_ON_ZERO';

DROP TABLE IF EXISTS `sites`;
CREATE TABLE `sites` (
  `ID` int(11) NOT NULL,
  `name` varchar(200) NOT NULL,
  `accessTime` varchar(500) NOT NULL,
  PRIMARY KEY (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

DROP TABLE IF EXISTS `admins`;
CREATE TABLE `admins` (
  `ID` int(11) NOT NULL,
  `login` varchar(200) NOT NULL,
  `password` varchar(500) NOT NULL,
  PRIMARY KEY (`login`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

DROP TABLE IF EXISTS `sessions`;
CREATE TABLE `sessions` (
    `id` varchar(200) NOT NULL,
    `userid` varchar(200) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8;