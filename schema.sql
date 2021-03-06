DROP TABLE IF EXISTS `position`;
CREATE TABLE `position` (
  `id` mediumint(9) NOT NULL AUTO_INCREMENT,
  `lat` varchar(100) COLLATE utf8_unicode_ci NOT NULL,
  `lon` varchar(100) COLLATE utf8_unicode_ci NOT NULL,
  `alt` varchar(100) COLLATE utf8_unicode_ci DEFAULT NULL,
  `at` datetime NOT NULL,
  `raw` longtext COLLATE utf8_unicode_ci NOT NULL,
  `origin` varchar(50) COLLATE utf8_unicode_ci NOT NULL,
  `created_at` datetime DEFAULT NULL,
  `batt` varchar(100) COLLATE utf8_unicode_ci DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;