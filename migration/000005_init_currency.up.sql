DROP TABLE IF EXISTS `currency`;
CREATE TABLE `currency` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(60) CHARACTER SET utf8 COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `get_address` varchar(60) CHARACTER SET utf8 COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `sort` int(11) NOT NULL DEFAULT '0',
  `logo` varchar(60) CHARACTER SET utf8 COLLATE utf8_unicode_ci DEFAULT '',
  `is_display` tinyint(4) NOT NULL DEFAULT '0',
  `min_number` decimal(23,8) NOT NULL DEFAULT '0.00000000' COMMENT '最小提币数量',
  `max_number` decimal(20,8) NOT NULL DEFAULT '0.00000000' COMMENT '最大提币数量',
  `rate` decimal(20,4) NOT NULL DEFAULT '0.0000' COMMENT '费率',
  `is_lever` tinyint(4) NOT NULL DEFAULT '0' COMMENT '是否杠杆币 0否 1是',
  `is_legal` tinyint(4) NOT NULL DEFAULT '0' COMMENT '是否法币 0否 1是',
  `is_match` tinyint(4) NOT NULL DEFAULT '0' COMMENT '是否撮合交易 0否 1是',
  `is_micro` tinyint(4) NOT NULL DEFAULT '0' COMMENT '是否微交易 0.否1是',
  `insurancable` tinyint(4) unsigned NOT NULL DEFAULT '0' COMMENT '是否可买保险',
  `show_legal` tinyint(4) NOT NULL DEFAULT '0' COMMENT '是否显示法币商家 0否 1是',
  `type` varchar(20) CHARACTER SET utf8 COLLATE utf8_unicode_ci NOT NULL DEFAULT '' COMMENT '基于哪个区块链',
  `black_limt` int(11) NOT NULL DEFAULT '0' COMMENT '币种黑名单限制数量',
  `key` varchar(1024) CHARACTER SET utf8 COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `contract_address` varchar(255) CHARACTER SET utf8 COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `total_account` varchar(255) CHARACTER SET utf8 COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `collect_account` varchar(300) CHARACTER SET utf8 COLLATE utf8_unicode_ci DEFAULT '0' COMMENT '归拢地址',
  `currency_decimals` int(11) DEFAULT NULL,
  `rmb_relation` decimal(23,2) NOT NULL DEFAULT '0.00' COMMENT '折合人民币比例',
  `decimal_scale` int(11) NOT NULL DEFAULT '18' COMMENT '发布小数点',
  `chain_fee` decimal(20,8) NOT NULL DEFAULT '0.00000000',
  `price` decimal(20,8) NOT NULL DEFAULT '0.00000000' COMMENT '价值(美元)',
  `micro_trade_fee` decimal(10,2) NOT NULL DEFAULT '0.00' COMMENT '微交易手续费%',
  `micro_min` decimal(20,0) DEFAULT '0' COMMENT '最小下单数量',
  `micro_max` decimal(20,0) DEFAULT '0' COMMENT '最大下单数量',
  `micro_holdtrade_max` int(11) NOT NULL DEFAULT '0' COMMENT '最大持仓笔数',
  `create_time` int(11) NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`) USING BTREE,
  KEY `name` (`name`) USING BTREE,
  KEY `name_2` (`name`,`is_display`) USING BTREE
) ENGINE=InnoDB  DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci ROW_FORMAT=DYNAMIC COMMENT="币种表";

DROP TABLE IF EXISTS `currency_deposit`;
CREATE TABLE `currency_deposit` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `currency_id` int(11) NOT NULL,
  `day` smallint(6) DEFAULT NULL,
  `save_min` decimal(20,8) DEFAULT '0.00000000',
  `total_interest_rate` tinyint(4) DEFAULT NULL COMMENT '百分比',
  `created_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB  DEFAULT CHARSET=utf8;

DROP TABLE IF EXISTS `currency_deposit_order`;
CREATE TABLE `currency_deposit_order` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `u_id` int(11) NOT NULL,
  `currency_id` int(11) NOT NULL,
  `amount` decimal(20,8) NOT NULL COMMENT '存款数量',
  `day_rate` decimal(10,4) NOT NULL COMMENT '日利率',
  `total_rate` decimal(10,4) DEFAULT NULL,
  `total_interest` decimal(20,8) DEFAULT '0.00000000' COMMENT '总利息',
  `last_settle_time` date DEFAULT NULL COMMENT '上次结息时间',
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  `end_at` date DEFAULT NULL COMMENT '结束日期',
  `start_at` date DEFAULT NULL COMMENT '开始时间',
  `status` tinyint(4) DEFAULT '1' COMMENT '1:进行中 2：已结束',
  `is_cancel` tinyint(1) DEFAULT NULL COMMENT '是否毁约',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

DROP TABLE IF EXISTS `currency_matches`;
CREATE TABLE `currency_matches` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `legal_id` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '法币id',
  `currency_id` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '币种id',
  `is_display` tinyint(4) NOT NULL DEFAULT '0' COMMENT '是否显示',
  `market_from` tinyint(4) NOT NULL DEFAULT '0' COMMENT '0.无,1.交易所,2.火币接口',
  `open_transaction` tinyint(4) NOT NULL DEFAULT '0' COMMENT '开启撮合交易',
  `open_lever` tinyint(4) NOT NULL DEFAULT '0' COMMENT '开启杠杆交易',
  `open_microtrade` tinyint(4) NOT NULL DEFAULT '0' COMMENT '开启微交易',
  `open_coin_trade` tinyint(4) NOT NULL COMMENT '开启币币交易',
  `sort` int(11) NOT NULL DEFAULT '0' COMMENT '排序',
  `micro_trade_fee` decimal(10,2) NOT NULL DEFAULT '0.00' COMMENT '微交易手续费(百分比)',
  `lever_share_num` decimal(20,8) NOT NULL DEFAULT '1.00000000' COMMENT '每手折合数量',
  `spread` decimal(20,4) NOT NULL DEFAULT '0.0000' COMMENT '点差',
  `overnight` decimal(20,4) NOT NULL DEFAULT '0.0000' COMMENT '隔夜费',
  `lever_trade_fee` decimal(20,4) NOT NULL DEFAULT '0.0000' COMMENT '交易手续费(百分比)',
  `lever_min_share` int(11) unsigned NOT NULL DEFAULT '1' COMMENT '杠杆交易最低手数',
  `lever_max_share` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '杠杆交易最高手数',
  `fluctuate_min` decimal(20,8) NOT NULL DEFAULT '0.00000000' COMMENT '行情波动最小值',
  `fluctuate_max` decimal(20,8) NOT NULL DEFAULT '0.00000000' COMMENT '行情波动最大值',
  `risk_group_result` tinyint(4) NOT NULL DEFAULT '0' COMMENT '群控结果',
  `create_time` int(11) NOT NULL DEFAULT '0' COMMENT '创建时间',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB  DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci ROW_FORMAT=DYNAMIC;

DROP TABLE IF EXISTS `currency_quotation`;
CREATE TABLE `currency_quotation` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `match_id` int(11) NOT NULL DEFAULT '0' COMMENT '交易对id',
  `legal_id` int(11) NOT NULL DEFAULT '0',
  `currency_id` int(11) NOT NULL DEFAULT '0',
  `change` varchar(255) CHARACTER SET utf8 COLLATE utf8_unicode_ci NOT NULL DEFAULT '' COMMENT '涨跌幅 带+ - 号',
  `volume` decimal(20,5) NOT NULL DEFAULT '0.00000' COMMENT '成交量',
  `now_price` decimal(20,5) NOT NULL DEFAULT '0.00000' COMMENT '当前价位',
  `add_time` int(11) NOT NULL DEFAULT '0',
  `xm` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB  DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci ROW_FORMAT=COMPACT;


