DROP TABLE IF EXISTS `micro_numbers`;
CREATE TABLE `micro_numbers` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `currency_id` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '币种id',
  `number` decimal(20,8) NOT NULL DEFAULT '0.00000000' COMMENT '数量',
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=50 DEFAULT CHARSET=utf8;

DROP TABLE IF EXISTS `micro_orders`;
CREATE TABLE `micro_orders` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `user_id` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '用户id',
  `match_id` int(11) unsigned NOT NULL COMMENT '交易对id',
  `currency_id` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '支付的币种',
  `type` tinyint(4) NOT NULL DEFAULT '0' COMMENT '买卖类型1.买涨,2.买跌',
  `is_insurance` tinyint(4) NOT NULL DEFAULT '0' COMMENT '订单险种:0.无,1.正向，2反向。',
  `seconds` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '秒数',
  `number` decimal(20,8) unsigned NOT NULL DEFAULT '0.00000000' COMMENT '下单数量',
  `open_price` decimal(20,8) unsigned NOT NULL DEFAULT '0.00000000' COMMENT '开仓价',
  `end_price` decimal(20,8) unsigned NOT NULL DEFAULT '0.00000000' COMMENT '收盘价',
  `fee` decimal(20,8) unsigned NOT NULL COMMENT '手续费',
  `profit_ratio` decimal(20,2) unsigned NOT NULL COMMENT '收益率',
  `fact_profits` decimal(20,8) NOT NULL COMMENT '最终收益',
  `status` tinyint(4) NOT NULL DEFAULT '0' COMMENT '状态:1.交易中,2.平仓中,3.已平仓',
  `pre_profit_result` tinyint(4) NOT NULL DEFAULT '0' COMMENT '预设盈利状态:-1.亏损,0.未设置,1.盈利',
  `profit_result` tinyint(4) NOT NULL DEFAULT '0' COMMENT '盈利结果:-1.亏损,0.平,1.盈利',
  `created_at` timestamp NULL DEFAULT NULL COMMENT '提交日期',
  `updated_at` timestamp NULL DEFAULT NULL COMMENT '更新日期',
  `handled_at` timestamp NULL DEFAULT NULL COMMENT '平仓时间',
  `complete_at` timestamp NULL DEFAULT NULL COMMENT '完成时间',
  `return_at` timestamp NULL DEFAULT NULL COMMENT '返还手续费的时间',
  `agent_path` varchar(2048) DEFAULT '1' COMMENT '代理商关系',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB  DEFAULT CHARSET=utf8;

DROP TABLE IF EXISTS `micro_seconds`;
CREATE TABLE `micro_seconds` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `seconds` int(11) unsigned NOT NULL COMMENT '秒数',
  `min_num` int(11) NOT NULL DEFAULT '0' COMMENT '最小下注限制（0）不限',
  `status` tinyint(4) NOT NULL DEFAULT '0' COMMENT '状态:0.禁用,1.启用',
  `profit_ratio` decimal(10,2) NOT NULL DEFAULT '0.00' COMMENT '收益率',
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;