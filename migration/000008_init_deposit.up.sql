DROP TABLE IF EXISTS `lh_deposit_config`;
CREATE TABLE `lh_deposit_config` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(255) DEFAULT NULL,
  `currency_id` int(11) NOT NULL,
  `interest_rate` decimal(10,4) NOT NULL,
  `save_min` decimal(20,8) NOT NULL DEFAULT '0.00000000',
  `day` int(10) NOT NULL,
  `created_at` datetime NOT NULL,
  `type` tinyint(1) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

DROP TABLE IF EXISTS `lh_deposit_order`;
CREATE TABLE `lh_deposit_order` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `user_id` int(11) DEFAULT NULL,
  `currency_id` int(11) NOT NULL,
  `amount` int(11) NOT NULL COMMENT '解冻数量',
  `day_rate` decimal(10,5) NOT NULL COMMENT '日利率',
  `total_interest` decimal(20,8) DEFAULT '0.00000000' COMMENT '总利息',
  `last_settle_time` date DEFAULT NULL COMMENT '上次结息时间',
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  `end_at` date DEFAULT NULL COMMENT '结束日期',
  `start_at` date DEFAULT NULL COMMENT '开始时间',
  `status` tinyint(4) DEFAULT '1' COMMENT '1:进行中 2：已结束',
  `is_return_reward` tinyint(1) DEFAULT '0' COMMENT '是否执行反佣逻辑',
  `is_cancel` tinyint(1) DEFAULT '0' COMMENT '是否毁约',
  `config_id` int(11) NOT NULL,
  `type` tinyint(1) DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8 ROW_FORMAT=COMPACT;

DROP TABLE IF EXISTS `lh_deposit_order_bak`;
CREATE TABLE `lh_deposit_order_bak` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `user_id` int(11) DEFAULT NULL,
  `currency_id` int(11) NOT NULL,
  `amount` int(11) NOT NULL COMMENT '解冻数量',
  `day_rate` decimal(10,5) NOT NULL COMMENT '日利率',
  `total_interest` decimal(20,8) DEFAULT '0.00000000' COMMENT '总利息',
  `last_settle_time` date DEFAULT NULL COMMENT '上次结息时间',
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  `end_at` date DEFAULT NULL COMMENT '结束日期',
  `start_at` date DEFAULT NULL COMMENT '开始时间',
  `status` tinyint(4) DEFAULT '1' COMMENT '1:进行中 2：已结束',
  `is_return_reward` tinyint(1) DEFAULT '0' COMMENT '是否执行反佣逻辑',
  `is_cancel` tinyint(1) DEFAULT '0' COMMENT '是否毁约',
  `config_id` int(11) NOT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8 ROW_FORMAT=COMPACT;

DROP TABLE IF EXISTS `lh_deposit_order_log`;
CREATE TABLE `lh_deposit_order_log` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `interest_amount` decimal(25,8) DEFAULT NULL COMMENT '利息数',
  `lh_order_id` int(11) DEFAULT NULL,
  `user_id` int(11) DEFAULT NULL,
  `created_at` datetime DEFAULT NULL COMMENT '什么时候结的',
  `interest_day` date DEFAULT NULL COMMENT '哪天的利息',
  `currency_id` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `orderIndex` (`lh_order_id`,`interest_day`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8 ROW_FORMAT=COMPACT;

DROP TABLE IF EXISTS `lh_loan_order`;
CREATE TABLE `lh_loan_order` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `bank_account_id` int(11) DEFAULT NULL,
  `amount` int(11) NOT NULL COMMENT '存款数量',
  `day_rate` decimal(10,5) NOT NULL COMMENT '日利率',
  `total_interest` decimal(20,8) DEFAULT '0.00000000' COMMENT '总利息',
  `total_return` decimal(20,8) DEFAULT NULL COMMENT '总还款数',
  `last_settle_time` date DEFAULT NULL COMMENT '上次结息时间',
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  `end_at` date DEFAULT NULL COMMENT '结束日期',
  `start_at` date DEFAULT NULL COMMENT '开始时间',
  `status` tinyint(4) DEFAULT '1' COMMENT '1:进行中 2：已还清',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 ROW_FORMAT=COMPACT;