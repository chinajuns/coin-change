DROP TABLE IF EXISTS `market`;
CREATE TABLE `market` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(60) CHARACTER SET utf8 COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `symbol` varchar(60) CHARACTER SET utf8 COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `website_slug` varchar(60) CHARACTER SET utf8 COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `rank` int(11) NOT NULL DEFAULT '0',
  `circulating_supply` bigint(20) unsigned NOT NULL DEFAULT '0',
  `total_supply` bigint(20) unsigned NOT NULL DEFAULT '0',
  `max_supply` bigint(20) unsigned NOT NULL DEFAULT '0',
  `quotes` text CHARACTER SET utf8 COLLATE utf8_unicode_ci NOT NULL,
  `last_updated` int(10) unsigned NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

DROP TABLE IF EXISTS `market_day`;
CREATE TABLE `market_day` (
  `id` int(11) NOT NULL,
  `currency_id` int(11) NOT NULL DEFAULT '0',
  `legal_id` int(11) NOT NULL DEFAULT '0',
  `start_price` decimal(20,5) NOT NULL DEFAULT '0.00000',
  `end_price` decimal(20,5) NOT NULL DEFAULT '0.00000',
  `highest` decimal(20,5) NOT NULL DEFAULT '0.00000',
  `mminimum` decimal(20,5) NOT NULL DEFAULT '0.00000',
  `number` decimal(20,5) NOT NULL DEFAULT '0.00000',
  `times` varchar(36) CHARACTER SET utf8 COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `mar_id` varchar(100) CHARACTER SET utf8 COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `type` int(11) NOT NULL DEFAULT '0'
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

DROP TABLE IF EXISTS `market_hour`;
CREATE TABLE `market_hour` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `currency_id` int(11) NOT NULL DEFAULT '0',
  `legal_id` int(11) NOT NULL DEFAULT '0',
  `start_price` decimal(20,5) NOT NULL DEFAULT '0.00000',
  `end_price` decimal(20,5) NOT NULL DEFAULT '0.00000',
  `highest` decimal(20,5) NOT NULL DEFAULT '0.00000',
  `mminimum` decimal(20,5) NOT NULL DEFAULT '0.00000',
  `day_time` int(11) NOT NULL DEFAULT '0',
  `type` tinyint(4) NOT NULL DEFAULT '0',
  `number` decimal(20,5) NOT NULL DEFAULT '0.00000',
  `mar_id` varchar(100) CHARACTER SET utf8 COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `period` varchar(100) CHARACTER SET utf8 COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `sign` tinyint(4) NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`) USING BTREE,
  KEY `currency_id` (`currency_id`,`legal_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci ROW_FORMAT=COMPACT;