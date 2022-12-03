DROP TABLE IF EXISTS `users`;
CREATE TABLE `users` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `area_code_id` int(11) NOT NULL DEFAULT '1' COMMENT '国家区号 1默认大陆',
  `area_code` int(10) NOT NULL COMMENT '区号',
  `account_number` varchar(30) CHARACTER SET utf8 COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `type` tinyint(4) NOT NULL DEFAULT '0',
  `phone` varchar(60) CHARACTER SET utf8 COLLATE utf8_unicode_ci DEFAULT NULL,
  `agent_id` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '0表示不是代理商，1以上表示该代理商id',
  `agent_note_id` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '代理商节点id。当该用户是代理商时该值等于上级代理商Id，当该用户不是代理商时该值等于节点代理商id',
  `parent_id` int(11) DEFAULT '0',
  `email` varchar(60) CHARACTER SET utf8 COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `password` varchar(255) CHARACTER SET utf8 COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `pay_password` varchar(255) CHARACTER SET utf8 COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `time` int(11) NOT NULL DEFAULT '0',
  `head_portrait` varchar(400) CHARACTER SET utf8 COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `extension_code` varchar(10) CHARACTER SET utf8 COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `status` int(11) NOT NULL DEFAULT '0' COMMENT '1，已锁定',
  `gesture_password` varchar(255) CHARACTER SET utf8 COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `is_auth` varchar(255) CHARACTER SET utf8 COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `nickname` varchar(255) CHARACTER SET utf8 COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `wallet_address` varchar(255) CHARACTER SET utf8 COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `is_blacklist` tinyint(4) NOT NULL DEFAULT '0',
  `parents_path` text CHARACTER SET utf8 COLLATE utf8_unicode_ci COMMENT '上级推荐人节点',
  `push_status` int(4) DEFAULT '0' COMMENT '0:未实名认证 1:实名认证  2:直推3人 3:直推5人 4:直推10人  5:直推30人  6:直推50人',
  `candy_number` decimal(10,4) DEFAULT '0.0000' COMMENT '糖果数量',
  `zhitui_real_number` int(11) DEFAULT '0' COMMENT '实名认证过的直推人数',
  `real_teamnumber` int(11) unsigned DEFAULT '0' COMMENT '实名认证通过的团队人数',
  `top_upnumber` decimal(20,6) DEFAULT '0.000000' COMMENT '团队业绩充值金额',
  `is_realname` int(4) NOT NULL DEFAULT '1' COMMENT '1:未实名认证过  2：实名认证过',
  `is_atelier` int(11) NOT NULL DEFAULT '0' COMMENT '是否工作室',
  `new_isreal_time` int(11) DEFAULT '0' COMMENT '最新通过的下级实名认证时间',
  `today_real_teamnumber` int(11) DEFAULT '0' COMMENT '今日新增团队实名认证人数',
  `today_LegalDealCancel_num` int(4) DEFAULT '0' COMMENT '今天c2c订单已经取消次数',
  `LegalDealCancel_num__update_time` int(11) DEFAULT NULL COMMENT 'c2c取消单子更新时间',
  `risk` tinyint(1) NOT NULL DEFAULT '0' COMMENT '-1.亏损,0.正常,1.盈利',
  `lock_time` int(11) NOT NULL DEFAULT '0' COMMENT '锁定时间',
  `level` int(25) DEFAULT '0' COMMENT '代数',
  `fund` decimal(30,8) DEFAULT '0.00000000' COMMENT '秒合约资产',
  `is_service` tinyint(4) unsigned DEFAULT '0' COMMENT '是否是客服',
  `agent_path` varchar(2048) CHARACTER SET utf8 COLLATE utf8_unicode_ci DEFAULT '1' COMMENT '代理商关系',
  `wallet_pwd` varchar(255) CHARACTER SET utf8 COLLATE utf8_unicode_ci NOT NULL,
  `country_code` varchar(10) CHARACTER SET utf8 COLLATE utf8_unicode_ci DEFAULT NULL,
  `label` varchar(30) CHARACTER SET utf8 COLLATE utf8_unicode_ci DEFAULT NULL,
  `nationality` varchar(30) CHARACTER SET utf8 COLLATE utf8_unicode_ci DEFAULT NULL,
  `last_login_ip` varchar(30) CHARACTER SET utf8 COLLATE utf8_unicode_ci DEFAULT NULL,
  `score` decimal(10,2) DEFAULT '0.00',
  `remark` varchar(255) CHARACTER SET utf8 COLLATE utf8_unicode_ci DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci COMMENT="用户表";

DROP TABLE IF EXISTS `user_algebra`;
CREATE TABLE `user_algebra` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `user_id` int(11) NOT NULL,
  `touch_user_id` int(11) NOT NULL COMMENT '触发者',
  `algebra` int(11) NOT NULL,
  `value` decimal(25,8) NOT NULL,
  `info` varchar(255) NOT NULL,
  `created_at` timestamp NOT NULL,
  `updated_at` timestamp NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

DROP TABLE IF EXISTS `user_cash_info`;
CREATE TABLE `user_cash_info` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `user_id` int(11) NOT NULL DEFAULT '0',
  `bank_id` int(10) NOT NULL COMMENT '银行id',
  `bank_name` varchar(60) CHARACTER SET utf8 COLLATE utf8_unicode_ci NOT NULL DEFAULT '' COMMENT '银行',
  `bank_branch` varchar(60) CHARACTER SET utf8 COLLATE utf8_unicode_ci NOT NULL COMMENT '支行',
  `bank_account` varchar(60) CHARACTER SET utf8 COLLATE utf8_unicode_ci NOT NULL DEFAULT '' COMMENT '账户',
  `real_name` varchar(60) CHARACTER SET utf8 COLLATE utf8_unicode_ci NOT NULL DEFAULT '' COMMENT '真实姓名',
  `alipay_account` varchar(60) CHARACTER SET utf8 COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `wechat_nickname` varchar(60) CHARACTER SET utf8 COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `wechat_account` varchar(60) CHARACTER SET utf8 COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `create_time` int(11) NOT NULL DEFAULT '0',
  `alipay_qr_code` varchar(255) CHARACTER SET utf8 COLLATE utf8_unicode_ci DEFAULT '',
  `wechat_qr_code` varchar(255) CHARACTER SET utf8 COLLATE utf8_unicode_ci DEFAULT '',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=2489 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

DROP TABLE IF EXISTS `user_chat`;
CREATE TABLE `user_chat` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `from_user_id` int(11) NOT NULL DEFAULT '0',
  `to_user_id` int(11) NOT NULL DEFAULT '0',
  `content` varchar(500) CHARACTER SET utf8 COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `offline` tinyint(4) NOT NULL DEFAULT '0',
  `type` varchar(60) CHARACTER SET utf8 COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `add_time` int(11) NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

DROP TABLE IF EXISTS `user_profiles`;
CREATE TABLE `user_profiles` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `user_id` int(11) NOT NULL DEFAULT '0' COMMENT '用户id',
  `name` varchar(50) NOT NULL DEFAULT '' COMMENT '真实姓名',
  `card_id` varchar(32) NOT NULL DEFAULT '' COMMENT '证件号码',
  `front_pic` varchar(255) NOT NULL DEFAULT '' COMMENT '身份证正面',
  `reverse_pic` varchar(255) NOT NULL DEFAULT '' COMMENT '身份证反面',
  `hand_pic` varchar(255) NOT NULL DEFAULT '' COMMENT '手持身份证',
  `created_at` timestamp NOT NULL   COMMENT '创建时间',
  `updated_at` timestamp NOT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 ROW_FORMAT=DYNAMIC;

DROP TABLE IF EXISTS `user_real`;
CREATE TABLE `user_real` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `user_id` int(11) NOT NULL DEFAULT '0',
  `name` varchar(60) CHARACTER SET utf8 COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `card_id` varchar(60) CHARACTER SET utf8 COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `type` int(11) NOT NULL DEFAULT '1' COMMENT '1初级认证 2高级认证',
  `review_status` tinyint(4) NOT NULL DEFAULT '1' COMMENT '1,未审核2,已审核',
  `front_pic` varchar(255) CHARACTER SET utf8 COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `reverse_pic` varchar(255) CHARACTER SET utf8 COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `hand_pic` varchar(60) CHARACTER SET utf8 COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `create_time` int(11) NOT NULL DEFAULT '0',
  `review_time` int(11) NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`),
  KEY `user_id` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

DROP TABLE IF EXISTS `users_insurances`;
CREATE TABLE `users_insurances` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `user_id` int(10) unsigned NOT NULL,
  `insurance_type_id` mediumint(5) unsigned NOT NULL COMMENT '险种类型id',
  `amount` float(10,2) unsigned NOT NULL DEFAULT '0.00' COMMENT '投保金额（受保资产）',
  `insurance_amount` float(10,2) unsigned NOT NULL DEFAULT '0.00' COMMENT '购买时的保险资产',
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
  `yielded_at` timestamp NULL DEFAULT NULL COMMENT '生币时间',
  `rescinded_at` timestamp NULL DEFAULT NULL COMMENT '解约时间',
  `rescinded_type` tinyint(4) unsigned DEFAULT '0' COMMENT '解约类型0，自动解约。1，手动解约',
  `status` tinyint(1) unsigned NOT NULL DEFAULT '1' COMMENT '状态，1，生效中。0，已失效。',
  `claim_status` tinyint(1) unsigned NOT NULL DEFAULT '0' COMMENT '是否申请理赔中1,是',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 ROW_FORMAT=DYNAMIC;

DROP TABLE IF EXISTS `users_wallet`;
CREATE TABLE `users_wallet` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `user_id` int(11) NOT NULL DEFAULT '0',
  `currency` int(11) NOT NULL DEFAULT '0',
  `address` varchar(50) CHARACTER SET utf8 COLLATE utf8_unicode_ci DEFAULT '',
  `address_2` varchar(50) CHARACTER SET utf8 COLLATE utf8_unicode_ci DEFAULT NULL COMMENT '第二类型地址usdt',
  `legal_balance` decimal(25,8) NOT NULL DEFAULT '0.00000000' COMMENT '法币交易余额',
  `lock_legal_balance` decimal(25,8) NOT NULL DEFAULT '0.00000000',
  `change_balance` decimal(25,8) NOT NULL DEFAULT '0.00000000' COMMENT '币币交易余额',
  `lock_change_balance` decimal(25,8) NOT NULL DEFAULT '0.00000000',
  `lever_balance` decimal(25,8) NOT NULL DEFAULT '0.00000000' COMMENT '杠杆交易余额',
  `lever_balance_add_allnum` decimal(25,8) NOT NULL DEFAULT '0.00000000' COMMENT '资产兑换累加产生的杠杆值(作为入金的一部分）',
  `lock_lever_balance` decimal(25,8) NOT NULL DEFAULT '0.00000000',
  `micro_balance` decimal(25,8) NOT NULL DEFAULT '0.00000000' COMMENT '微盘',
  `lock_micro_balance` decimal(25,8) NOT NULL DEFAULT '0.00000000' COMMENT '锁定微盘',
  `insurance_balance` decimal(25,8) NOT NULL DEFAULT '0.00000000' COMMENT '受保资产',
  `lock_insurance_balance` decimal(25,8) NOT NULL DEFAULT '0.00000000' COMMENT '锁定受保资产  保险资产',
  `status` int(11) NOT NULL DEFAULT '0',
  `create_time` int(11) NOT NULL DEFAULT '0',
  `old_balance` decimal(20,8) NOT NULL DEFAULT '0.00000000',
  `private` varchar(1024) CHARACTER SET utf8 COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `cost` decimal(20,5) NOT NULL DEFAULT '0.00000' COMMENT '持仓成本',
  `gl_time` varchar(50) CHARACTER SET utf8 COLLATE utf8_unicode_ci DEFAULT NULL,
  `txid` varchar(1024) CHARACTER SET utf8 COLLATE utf8_unicode_ci NOT NULL DEFAULT '' COMMENT '交易哈希',
  PRIMARY KEY (`id`),
  KEY `user_id` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci COMMENT="用户钱包";

DROP TABLE IF EXISTS `users_wallet_out`;
CREATE TABLE `users_wallet_out` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `user_id` int(11) NOT NULL DEFAULT '0',
  `currency` int(11) NOT NULL DEFAULT '0',
  `address` varchar(50) CHARACTER SET utf8 COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `nettype` varchar(20) CHARACTER SET utf8 COLLATE utf8_unicode_ci DEFAULT NULL COMMENT '通道',
  `number` decimal(20,8) NOT NULL DEFAULT '0.00000000',
  `create_time` int(11) NOT NULL DEFAULT '0',
  `rate` decimal(20,4) NOT NULL DEFAULT '0.0000',
  `status` tinyint(4) NOT NULL DEFAULT '0',
  `notes` text CHARACTER SET utf8 COLLATE utf8_unicode_ci NOT NULL,
  `real_number` decimal(20,8) NOT NULL DEFAULT '0.00000000',
  `txid` varchar(1024) CHARACTER SET utf8 COLLATE utf8_unicode_ci NOT NULL DEFAULT '' COMMENT '链上哈希',
  `verificationcode` varchar(100) CHARACTER SET utf8 COLLATE utf8_unicode_ci DEFAULT '',
  `update_time` int(11) DEFAULT NULL,
  `is_bank` tinyint(1) NOT NULL COMMENT ' 1银行卡  0普通',
  `tibi_rmb` decimal(20,8) NOT NULL DEFAULT '0.00000000' COMMENT '提币金额 rmb',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci COMMENT="用户钱包提币表";

DROP TABLE IF EXISTS `wallet_log`;
CREATE TABLE `wallet_log` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `user_id` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '用户id',
  `from_user_id` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '触发者id',
  `account_log_id` int(11) NOT NULL DEFAULT '0' COMMENT '关联account_log',
  `wallet_id` int(11) NOT NULL DEFAULT '0' COMMENT '钱包id',
  `balance_type` int(11) NOT NULL DEFAULT '0' COMMENT '余额类型:1.法币,2.币币,3.杆杠',
  `lock_type` tinyint(3) unsigned NOT NULL DEFAULT '0' COMMENT '是否锁定',
  `before` decimal(20,8) NOT NULL DEFAULT '0.00000000' COMMENT '原余额',
  `change` decimal(20,8) NOT NULL DEFAULT '0.00000000' COMMENT '变动余额',
  `after` decimal(20,8) NOT NULL DEFAULT '0.00000000' COMMENT '当前余额',
  `memo` varchar(255) CHARACTER SET utf8 COLLATE utf8_unicode_ci NOT NULL DEFAULT '' COMMENT '备注',
  `extra_sign` int(11) NOT NULL DEFAULT '0' COMMENT '扩展标识',
  `extra_data` varchar(1024) CHARACTER SET utf8 COLLATE utf8_unicode_ci NOT NULL DEFAULT '' COMMENT '扩展数据',
  `create_time` int(11) NOT NULL DEFAULT '0' COMMENT '发生时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci COMMENT="钱包日志";
