DROP TABLE IF EXISTS `algebra`;
CREATE TABLE `algebra` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(25) NOT NULL,
  `algebra` int(25) NOT NULL DEFAULT '0' COMMENT '代数',
  `rate` decimal(25,4) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

DROP TABLE IF EXISTS `area_code`;
CREATE TABLE `area_code` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(255) NOT NULL,
  `lang` varchar(10) DEFAULT NULL,
  `area_code` varchar(255) NOT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8 ROW_FORMAT=COMPACT COMMENT='国家区号';

DROP TABLE IF EXISTS `auto_list`;
CREATE TABLE `auto_list` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `buy_user_id` int(11) NOT NULL DEFAULT '0' COMMENT '买方user_id',
  `sell_user_id` int(11) NOT NULL DEFAULT '0' COMMENT '卖方user_id',
  `currency_id` int(11) NOT NULL DEFAULT '0' COMMENT '币种id',
  `legal_id` int(11) NOT NULL DEFAULT '0' COMMENT '法币id',
  `min_price` decimal(20,5) NOT NULL DEFAULT '0.00000' COMMENT '最低',
  `max_price` decimal(20,5) NOT NULL DEFAULT '0.00000' COMMENT '最高',
  `min_number` decimal(20,5) NOT NULL DEFAULT '0.00000',
  `max_number` decimal(20,5) NOT NULL DEFAULT '0.00000',
  `need_second` int(11) NOT NULL DEFAULT '0',
  `create_time` int(11) NOT NULL DEFAULT '0',
  `is_start` tinyint(4) NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

DROP TABLE IF EXISTS `bank`;
CREATE TABLE `bank` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(60) CHARACTER SET utf8 COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

DROP TABLE IF EXISTS `c2c_deal`;
CREATE TABLE `c2c_deal` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `legal_deal_send_id` int(11) NOT NULL DEFAULT '0',
  `user_id` int(11) NOT NULL DEFAULT '0',
  `seller_id` int(11) NOT NULL DEFAULT '0' COMMENT '发布方用户id',
  `number` decimal(20,5) NOT NULL DEFAULT '0.00000',
  `is_sure` tinyint(4) NOT NULL DEFAULT '0' COMMENT '0 未完成 1 已完成  2取消  3已付款',
  `create_time` int(11) NOT NULL DEFAULT '0',
  `update_time` int(11) NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci ROW_FORMAT=DYNAMIC;

DROP TABLE IF EXISTS `c2c_deal_send`;
CREATE TABLE `c2c_deal_send` (
  `id` int(10) unsigned NOT NULL,
  `seller_id` int(11) NOT NULL DEFAULT '0' COMMENT '发布方用户id',
  `currency_id` int(11) NOT NULL DEFAULT '0',
  `type` enum('buy','sell') CHARACTER SET utf8 COLLATE utf8_unicode_ci NOT NULL DEFAULT 'sell',
  `way` enum('bank','we_chat','ali_pay') CHARACTER SET utf8 COLLATE utf8_unicode_ci NOT NULL DEFAULT 'bank',
  `price` decimal(20,5) NOT NULL DEFAULT '0.00000',
  `total_number` decimal(20,5) NOT NULL DEFAULT '0.00000',
  `surplus_number` decimal(20,5) NOT NULL DEFAULT '0.00000',
  `min_number` decimal(20,5) NOT NULL DEFAULT '0.00000' COMMENT '最小购买量',
  `is_done` tinyint(4) NOT NULL DEFAULT '0' COMMENT ' 0 未完成  1完成    2 24小时未交易取消',
  `create_time` int(11) NOT NULL DEFAULT '0'
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci ROW_FORMAT=DYNAMIC;

DROP TABLE IF EXISTS `candy_transfer`;
CREATE TABLE `candy_transfer` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `from_user_id` int(10) unsigned NOT NULL COMMENT '转出用户id',
  `to_user_id` int(10) unsigned NOT NULL COMMENT '转入用户id',
  `transfer_qty` decimal(20,6) unsigned NOT NULL COMMENT '转账数量',
  `transfer_rate` decimal(20,2) unsigned NOT NULL COMMENT '手续费率(百分比)',
  `transfer_fee` decimal(20,6) unsigned NOT NULL COMMENT '手续费',
  `create_time` int(10) unsigned NOT NULL COMMENT '转账时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 ROW_FORMAT=DYNAMIC COMMENT='糖果转账';

DROP TABLE IF EXISTS `chain_hashes`;
CREATE TABLE `chain_hashes` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `code` varchar(50) NOT NULL COMMENT '币种代码',
  `txid` varchar(1024) NOT NULL COMMENT '链上交易hash',
  `amount` decimal(20,8) NOT NULL DEFAULT '0.00000000' COMMENT '转账金额',
  `sender` varchar(1024) NOT NULL COMMENT '转出地址',
  `recipient` varchar(1024) NOT NULL COMMENT '转入地址',
  `status` tinyint(4) NOT NULL DEFAULT '0',
  `created_at` timestamp NULL DEFAULT NULL COMMENT '创建时间',
  `updated_at` timestamp NULL DEFAULT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8 ROW_FORMAT=DYNAMIC;

DROP TABLE IF EXISTS `charge_req`;
CREATE TABLE `charge_req` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `uid` int(11) NOT NULL,
  `amount` decimal(15,8) NOT NULL,
  `user_account` varchar(500) NOT NULL,
  `status` tinyint(4) NOT NULL,
  `currency_id` int(11) NOT NULL,
  `remark` varchar(200) NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  `is_bank` tinyint(1) NOT NULL DEFAULT '0' COMMENT '1银行卡充值  0普通充值',
  `daozhang_num` decimal(20,8) NOT NULL COMMENT '到账金额',
  `img` varchar(255) DEFAULT NULL COMMENT '图片',
  PRIMARY KEY (`id`),
  UNIQUE KEY `id` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=719 DEFAULT CHARSET=utf8;

DROP TABLE IF EXISTS `chat_log`;
CREATE TABLE `chat_log` (
  `id` int(10) unsigned NOT NULL,
  `type` tinyint(3) NOT NULL DEFAULT '1' COMMENT '1,文字 2,图片 3,视频,4,日期',
  `content` varchar(1024) NOT NULL,
  `from_user` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '发送人',
  `to_user` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '接收人',
  `created_at` timestamp NULL DEFAULT NULL COMMENT '发送时间',
  `updated_at` timestamp NULL DEFAULT NULL,
  `trade_type` tinyint(3) unsigned NOT NULL DEFAULT '0' COMMENT '交易类型0.c2c',
  `trade_id` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '交易的id',
  `readed` tinyint(3) unsigned NOT NULL DEFAULT '0' COMMENT '是否已读0,未读。1，已读'
) ENGINE=InnoDB DEFAULT CHARSET=utf8 ROW_FORMAT=DYNAMIC;

DROP TABLE IF EXISTS `coin_trade`;
CREATE TABLE `coin_trade` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `u_id` int(11) unsigned DEFAULT NULL,
  `currency_id` int(11) DEFAULT NULL,
  `legal_id` int(11) DEFAULT NULL,
  `type` tinyint(4) unsigned DEFAULT NULL COMMENT '1:buy 2:sell',
  `target_price` decimal(20,8) DEFAULT NULL COMMENT '目标价格',
  `trade_price` decimal(20,8) DEFAULT NULL COMMENT '交易当前价格',
  `trade_amount` decimal(20,8) DEFAULT NULL COMMENT '币数量',
  `charge_fee` decimal(10,6) DEFAULT NULL COMMENT '手续费',
  `status` tinyint(4) unsigned DEFAULT '1' COMMENT '状态 1 交易中 2 已完成',
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 ROW_FORMAT=COMPACT;

DROP TABLE IF EXISTS `conversion`;
CREATE TABLE `conversion` (
  `id` int(11) unsigned NOT NULL,
  `form_currency_id` int(11) NOT NULL,
  `to_currency_id` int(11) NOT NULL,
  `num` decimal(20,4) NOT NULL,
  `fee` decimal(20,4) NOT NULL,
  `sj_num` decimal(20,4) NOT NULL,
  `create_time` int(11) NOT NULL,
  `user_id` int(11) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8 ROW_FORMAT=DYNAMIC;


DROP TABLE IF EXISTS `exception_log`;
CREATE TABLE `exception_log` (
  `id` int(11) NOT NULL,
  `title` varchar(100) DEFAULT NULL,
  `message` longtext,
  `created_at` datetime DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8 ROW_FORMAT=COMPACT;

DROP TABLE IF EXISTS `failed_jobs`;
CREATE TABLE `failed_jobs` (
  `id` bigint(20) unsigned NOT NULL,
  `connection` text CHARACTER SET utf8 COLLATE utf8_unicode_ci NOT NULL,
  `queue` text CHARACTER SET utf8 COLLATE utf8_unicode_ci NOT NULL,
  `payload` longtext CHARACTER SET utf8 COLLATE utf8_unicode_ci NOT NULL,
  `exception` longtext CHARACTER SET utf8 COLLATE utf8_unicode_ci NOT NULL,
  `failed_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

DROP TABLE IF EXISTS `false_data`;
CREATE TABLE `false_data` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `address` varchar(255) CHARACTER SET utf8 COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `price` varchar(20) CHARACTER SET utf8 COLLATE utf8_unicode_ci NOT NULL DEFAULT '' COMMENT 'log',
  `time` int(11) NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

DROP TABLE IF EXISTS `feedback`;
CREATE TABLE `feedback` (
  `id` int(10) unsigned NOT NULL,
  `user_id` int(11) NOT NULL DEFAULT '0',
  `content` text CHARACTER SET utf8 COLLATE utf8_unicode_ci NOT NULL,
  `img` varchar(255) CHARACTER SET utf8 COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `create_time` int(11) NOT NULL DEFAULT '0',
  `reply_time` int(11) NOT NULL DEFAULT '0',
  `reply_content` text CHARACTER SET utf8 COLLATE utf8_unicode_ci NOT NULL,
  `is_reply` tinyint(4) NOT NULL DEFAULT '0'
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

DROP TABLE IF EXISTS `flash_against`;
CREATE TABLE `flash_against` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `user_id` int(11) NOT NULL,
  `left_currency_id` int(11) NOT NULL,
  `right_currency_id` int(11) NOT NULL,
  `num` decimal(25,8) NOT NULL,
  `fee` decimal(25,8) DEFAULT NULL,
  `absolute_quantity` decimal(25,8) NOT NULL,
  `market_price` decimal(25,8) NOT NULL DEFAULT '0.00000000' COMMENT '当时的行情价格',
  `price` decimal(25,8) NOT NULL DEFAULT '0.00000000' COMMENT '用户输入的兑换价格价格',
  `status` tinyint(2) NOT NULL DEFAULT '0' COMMENT '状态 0 审核中 ：1成功 ：2 失败',
  `review_time` int(11) DEFAULT NULL,
  `create_time` int(11) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=10 DEFAULT CHARSET=utf8;

DROP TABLE IF EXISTS `gd_order`;
CREATE TABLE `gd_order` (
  `id` int(11) NOT NULL,
  `uid` int(11) NOT NULL,
  `gd_user_id` int(11) NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  `status` tinyint(1) NOT NULL DEFAULT '1',
  `value` int(11) NOT NULL,
  `day_max_value` int(11) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

DROP TABLE IF EXISTS `gd_user`;
CREATE TABLE `gd_user` (
  `id` int(11) NOT NULL,
  `uid` int(11) NOT NULL,
  `total_profit_rate` decimal(10,2) NOT NULL,
  `three_week_profit` decimal(10,2) NOT NULL,
  `total_day` mediumint(9) NOT NULL,
  `trade_count` mediumint(9) NOT NULL,
  `total_follower` int(11) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

DROP TABLE IF EXISTS `historical_data`;
CREATE TABLE `historical_data` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `type` varchar(10) CHARACTER SET utf8 COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `start_time` int(11) NOT NULL DEFAULT '0',
  `end_time` int(11) NOT NULL DEFAULT '0',
  `data` varchar(500) CHARACTER SET utf8 COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

DROP TABLE IF EXISTS `huobi_symbols`;
CREATE TABLE `huobi_symbols` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `base-currency` varchar(50) NOT NULL DEFAULT '',
  `quote-currency` varchar(50) NOT NULL DEFAULT '',
  `price-precision` int(11) NOT NULL DEFAULT '0',
  `amount-precision` int(11) NOT NULL DEFAULT '0',
  `symbol-partition` varchar(50) NOT NULL DEFAULT '',
  `symbol` varchar(50) NOT NULL DEFAULT '',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

DROP TABLE IF EXISTS `insurance_claim_applies`;
CREATE TABLE `insurance_claim_applies` (
  `id` int(11) NOT NULL,
  `user_id` int(10) unsigned NOT NULL COMMENT '用户id',
  `user_insurance_id` int(10) unsigned NOT NULL COMMENT '用户保险单',
  `insurance_type` tinyint(1) unsigned NOT NULL COMMENT '保险类型',
  `apply_status` tinyint(2) unsigned NOT NULL DEFAULT '0' COMMENT '赔付状态。0，申请中，1，已赔付。2，已拒绝',
  `compensate` float(10,2) NOT NULL DEFAULT '0.00' COMMENT '赔付金额',
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `operator` varchar(32) NOT NULL DEFAULT '' COMMENT '操作人,auto,自动。',
  `refuse_reason` varchar(200) NOT NULL COMMENT '拒绝理由'
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

DROP TABLE IF EXISTS `insurance_rules`;
CREATE TABLE `insurance_rules` (
  `id` int(11) NOT NULL,
  `insurance_type_id` smallint(5) unsigned NOT NULL COMMENT '险种id',
  `amount` float(10,2) unsigned NOT NULL DEFAULT '0.00' COMMENT '金额',
  `place_an_order_max` float(10,2) unsigned NOT NULL DEFAULT '0.00' COMMENT '下单最大金额限制',
  `existing_number` smallint(5) unsigned NOT NULL DEFAULT '0' COMMENT '可同时存在最大订单数'
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='不同金额保险的不同规则。';

DROP TABLE IF EXISTS `insurance_types`;
CREATE TABLE `insurance_types` (
  `id` int(11) NOT NULL,
  `name` varchar(20) NOT NULL,
  `currency_id` smallint(5) unsigned NOT NULL COMMENT '币种',
  `type` tinyint(2) unsigned NOT NULL DEFAULT '1' COMMENT '保险类型1，正向。2，反向。',
  `min_amount` float(10,2) unsigned NOT NULL DEFAULT '0.00' COMMENT '最低购买额',
  `max_amount` float(10,2) unsigned NOT NULL DEFAULT '0.00' COMMENT '最大购买额',
  `insurance_assets` float(10,2) unsigned NOT NULL DEFAULT '0.00' COMMENT '保险资产占比%',
  `profit_termination_condition` float(10,2) unsigned NOT NULL DEFAULT '0.00' COMMENT '盈利比例解约条件%',
  `defective_claims_condition` float(10,2) unsigned NOT NULL DEFAULT '0.00' COMMENT '亏损比例理赔条件%（正向）',
  `defective_claims_condition2` float(10,2) unsigned NOT NULL DEFAULT '0.00' COMMENT '亏损理赔条件2(反向)',
  `claims_times_daily` tinyint(2) unsigned NOT NULL DEFAULT '0' COMMENT '每日赔付次数',
  `auto_claim` tinyint(1) unsigned NOT NULL DEFAULT '0' COMMENT '自动赔付0，否。1，是',
  `claim_rate` float(10,2) unsigned NOT NULL DEFAULT '100.00' COMMENT '赔付比例',
  `claim_direction` tinyint(1) unsigned NOT NULL DEFAULT '1' COMMENT '赔付去向。1，受保资产。2，可用资产。',
  `status` tinyint(1) unsigned NOT NULL DEFAULT '1' COMMENT '状态1，开启.0，关闭',
  `is_t_add_1` tinyint(1) unsigned NOT NULL DEFAULT '1' COMMENT '是否T加1，1，是。0，不是。'
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='险种的类型。';

DROP TABLE IF EXISTS `jobs`;
CREATE TABLE `jobs` (
  `id` int(11) NOT NULL,
  `queue` varchar(255) CHARACTER SET utf8 COLLATE utf8_unicode_ci NOT NULL,
  `payload` longtext CHARACTER SET utf8 COLLATE utf8_unicode_ci NOT NULL,
  `attempts` tinyint(3) unsigned NOT NULL,
  `reserved_at` int(10) unsigned DEFAULT NULL,
  `available_at` int(10) unsigned NOT NULL,
  `created_at` int(10) unsigned NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

DROP TABLE IF EXISTS `js_news`;
CREATE TABLE `js_news` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `title` varchar(100) DEFAULT NULL,
  `summary` varchar(100) DEFAULT NULL,
  `content` longtext,
  `published_time` datetime DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

DROP TABLE IF EXISTS `lbx_hashes`;
CREATE TABLE `lbx_hashes` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `wallet_id` int(11) NOT NULL COMMENT '钱包id',
  `txid` varchar(1024) NOT NULL COMMENT '链上交易hash',
  `type` tinyint(4) NOT NULL DEFAULT '0' COMMENT '业务类型:0.归拢,1.提币 2.打入手续费',
  `amount` decimal(20,8) NOT NULL COMMENT '数量',
  `status` tinyint(4) NOT NULL DEFAULT '0' COMMENT '0 未处理  1处理成功   2处理失败',
  `created_at` timestamp NULL DEFAULT NULL COMMENT '创建时间',
  `updated_at` timestamp NULL DEFAULT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

DROP TABLE IF EXISTS `legal_deal`;
CREATE TABLE `legal_deal` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `legal_deal_send_id` int(11) NOT NULL DEFAULT '0',
  `user_id` int(11) NOT NULL DEFAULT '0',
  `seller_id` int(11) NOT NULL DEFAULT '0',
  `number` decimal(20,5) NOT NULL DEFAULT '0.00000',
  `is_sure` tinyint(4) NOT NULL DEFAULT '0' COMMENT '0未确认 1已确认 2已取消 3已付款',
  `pay_orders_img` varchar(255) CHARACTER SET utf8 COLLATE utf8_unicode_ci NOT NULL COMMENT '付款凭证',
  `create_time` int(11) NOT NULL DEFAULT '0',
  `update_time` int(11) NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

DROP TABLE IF EXISTS `legal_deal_send`;
CREATE TABLE `legal_deal_send` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `seller_id` int(11) NOT NULL DEFAULT '0',
  `currency_id` int(11) NOT NULL DEFAULT '0',
  `type` enum('buy','sell') CHARACTER SET utf8 COLLATE utf8_unicode_ci NOT NULL DEFAULT 'sell',
  `way` enum('bank','we_chat','ali_pay') CHARACTER SET utf8 COLLATE utf8_unicode_ci NOT NULL DEFAULT 'bank',
  `price` decimal(20,5) NOT NULL DEFAULT '0.00000',
  `total_number` decimal(20,5) NOT NULL DEFAULT '0.00000',
  `surplus_number` decimal(20,5) NOT NULL DEFAULT '0.00000',
  `min_number` decimal(20,5) NOT NULL DEFAULT '0.00000',
  `max_number` decimal(20,5) DEFAULT '0.00000',
  `is_done` tinyint(4) NOT NULL DEFAULT '0' COMMENT '0  1已完成  2撤回',
  `create_time` int(11) NOT NULL DEFAULT '0',
  `is_shelves` int(4) DEFAULT '1' COMMENT '1:上架   2下架',
  `is_sendback` int(11) DEFAULT '1' COMMENT '1:未撤回  2：撤回',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

DROP TABLE IF EXISTS `needle`;
CREATE TABLE `needle` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `open` decimal(25,8) DEFAULT NULL COMMENT '开',
  `high` decimal(25,8) DEFAULT NULL COMMENT '高',
  `low` decimal(25,8) DEFAULT NULL COMMENT '低',
  `close` decimal(25,8) DEFAULT NULL COMMENT '收',
  `symbol` varchar(255) DEFAULT NULL COMMENT '交易对名称 BTC/USDT',
  `base` varchar(24) DEFAULT NULL COMMENT '基础币种BTC',
  `target` varchar(24) DEFAULT NULL COMMENT '交易币种USDT',
  `itime` int(10) unsigned DEFAULT NULL COMMENT '精确到秒的时间戳',
  `updated_at` datetime DEFAULT NULL,
  `created_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

DROP TABLE IF EXISTS `prize_pool`;
CREATE TABLE `prize_pool` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `scene` int(11) NOT NULL DEFAULT '0' COMMENT '奖励场景',
  `reward_type` tinyint(4) NOT NULL DEFAULT '0' COMMENT '奖励类型:0.糖果,1.数字货币',
  `reward_currency` int(11) NOT NULL DEFAULT '0' COMMENT '奖励币种,不是数字货币传0',
  `currency_type` int(11) NOT NULL DEFAULT '0' COMMENT '货币类型:1.法币,2.币币交易,3.杠杆交易',
  `reward_qty` decimal(20,8) NOT NULL DEFAULT '0.00000000' COMMENT '奖励数量',
  `from_user_id` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '触发用户',
  `to_user_id` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '受奖励用户',
  `status` tinyint(4) NOT NULL DEFAULT '0' COMMENT '状态:[0.未领取,1.已领取,2.已过期]',
  `sign` int(11) NOT NULL DEFAULT '0' COMMENT '扩展标识,与奖励类型配合使用,用于区分具体的奖励,非必填项',
  `extra_data` varchar(512) NOT NULL DEFAULT '' COMMENT '附加数据',
  `memo` varchar(255) NOT NULL DEFAULT '' COMMENT '备注',
  `create_time` int(11) DEFAULT NULL COMMENT '奖励时间',
  `expire_time` int(11) DEFAULT NULL COMMENT '过期时间',
  `receive_time` int(11) DEFAULT NULL COMMENT '领取时间',
  `error_info` varchar(512) DEFAULT NULL COMMENT '错误信息',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='奖金池';

DROP TABLE IF EXISTS `robot`;
CREATE TABLE `robot` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `currency_id` int(11) NOT NULL DEFAULT '0',
  `legal_id` int(11) NOT NULL DEFAULT '0',
  `buy_user_id` int(11) NOT NULL DEFAULT '0',
  `sell_user_id` int(11) NOT NULL DEFAULT '0',
  `create_time` int(11) NOT NULL DEFAULT '0',
  `status` int(11) NOT NULL DEFAULT '0',
  `second` int(11) NOT NULL DEFAULT '0',
  `sell` int(11) NOT NULL DEFAULT '0',
  `buy` int(11) NOT NULL DEFAULT '0',
  `number_max` decimal(20,5) NOT NULL DEFAULT '0.00000',
  `number_min` decimal(20,5) NOT NULL DEFAULT '0.00000',
  `float_number_down` decimal(20,5) NOT NULL DEFAULT '0.00000',
  `float_number_up` decimal(20,5) NOT NULL DEFAULT '0.00000',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

DROP TABLE IF EXISTS `seller`;
CREATE TABLE `seller` (
  `id` int(10) unsigned NOT NULL,
  `user_id` int(11) NOT NULL DEFAULT '0',
  `name` varchar(60) CHARACTER SET utf8 COLLATE utf8_unicode_ci DEFAULT '',
  `seller_balance` decimal(20,5) DEFAULT '0.00000',
  `lock_seller_balance` decimal(20,5) DEFAULT '0.00000',
  `wechat_nickname` varchar(60) CHARACTER SET utf8 COLLATE utf8_unicode_ci DEFAULT '',
  `wechat_account` varchar(60) CHARACTER SET utf8 COLLATE utf8_unicode_ci DEFAULT '',
  `ali_nickname` varchar(60) CHARACTER SET utf8 COLLATE utf8_unicode_ci DEFAULT '',
  `ali_account` varchar(60) CHARACTER SET utf8 COLLATE utf8_unicode_ci DEFAULT '',
  `bank_id` int(11) DEFAULT '0',
  `bank_account` varchar(60) CHARACTER SET utf8 COLLATE utf8_unicode_ci DEFAULT '',
  `bank_address` varchar(60) CHARACTER SET utf8 COLLATE utf8_unicode_ci DEFAULT '',
  `create_time` int(11) NOT NULL DEFAULT '0',
  `currency_id` int(11) NOT NULL DEFAULT '0',
  `mobile` varchar(60) CHARACTER SET utf8 COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `alipay_qr_code` varchar(255) CHARACTER SET utf8 COLLATE utf8_unicode_ci DEFAULT '',
  `wechat_qr_code` varchar(255) CHARACTER SET utf8 COLLATE utf8_unicode_ci DEFAULT '',
  `status` int(2) NOT NULL DEFAULT '0' COMMENT '审核状态0：未通过  1：通过',
  KEY `user_id` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

DROP TABLE IF EXISTS `settings`;
CREATE TABLE `settings` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `key` varchar(50) CHARACTER SET utf8 COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `value` varchar(1000) CHARACTER SET utf8 COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `notes` varchar(255) CHARACTER SET utf8 COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  PRIMARY KEY (`id`),
  KEY `key` (`key`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

DROP TABLE IF EXISTS `tokens`;
CREATE TABLE `tokens` (
  `id` int(10) NOT NULL AUTO_INCREMENT,
  `token` varchar(255) CHARACTER SET utf8 COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `time_out` int(11) NOT NULL DEFAULT '0',
  `user_id` int(11) NOT NULL DEFAULT '0',
  `lang` varchar(10) CHARACTER SET utf8 COLLATE utf8_unicode_ci DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `tokenindex` (`token`)
) ENGINE=InnoDB  DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

DROP TABLE IF EXISTS `ztpay_log`;
CREATE TABLE `ztpay_log` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `unique_key` varchar(100) DEFAULT NULL,
  `body` longtext,
  `created_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 ROW_FORMAT=COMPACT;