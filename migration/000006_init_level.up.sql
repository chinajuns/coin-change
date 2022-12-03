DROP TABLE IF EXISTS `level`;
CREATE TABLE `level` (
  `id` int(11) NOT NULL,
  `name` varchar(25) NOT NULL DEFAULT '',
  `fill_currency` decimal(25,4) NOT NULL DEFAULT '0.0000' COMMENT '充币数量',
  `direct_drive_count` int(11) NOT NULL DEFAULT '0' COMMENT '直推数量',
  `direct_drive_price` decimal(25,4) NOT NULL DEFAULT '0.0000' COMMENT '直推金额',
  `max_algebra` int(20) NOT NULL DEFAULT '0' COMMENT '最大代数',
  `level` int(25) NOT NULL DEFAULT '0' COMMENT '级别'
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

DROP TABLE IF EXISTS `lever_multiple`;
CREATE TABLE `lever_multiple` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `type` int(11) DEFAULT '1' COMMENT '1倍数  2手数',
  `value` varchar(255) DEFAULT NULL,
  `currency_id` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB  DEFAULT CHARSET=utf8;

DROP TABLE IF EXISTS `lever_tolegal`;
CREATE TABLE `lever_tolegal` (
  `id` int(11) NOT NULL,
  `user_id` int(11) DEFAULT NULL,
  `number` decimal(11,5) DEFAULT NULL COMMENT '杠杆转c2c数量',
  `add_time` int(11) DEFAULT NULL,
  `type` int(11) DEFAULT '1' COMMENT '1:c2c转杠杆  2杠杆转c2c',
  `status` int(11) DEFAULT '1' COMMENT '1:未审核   2：审核通过 3:审核不通过'
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='杠杆账户转c2c账户后台审核';

DROP TABLE IF EXISTS `lever_transaction`;
CREATE TABLE `lever_transaction` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `type` tinyint(4) unsigned NOT NULL DEFAULT '0' COMMENT '买卖类型:1.买入,2.卖出',
  `user_id` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '用户id',
  `currency` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '交易id',
  `legal` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '法币id',
  `origin_price` decimal(20,8) unsigned NOT NULL DEFAULT '0.00000000' COMMENT '原始价格',
  `price` decimal(20,8) unsigned NOT NULL DEFAULT '0.00000000' COMMENT '开仓价格(点差处理之后)',
  `update_price` decimal(20,8) unsigned NOT NULL DEFAULT '0.00000000' COMMENT '当前价格',
  `target_profit_price` decimal(20,8) unsigned NOT NULL DEFAULT '0.00000000' COMMENT '止盈价格',
  `stop_loss_price` decimal(20,8) unsigned NOT NULL DEFAULT '0.00000000' COMMENT '止亏价格',
  `share` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '手数',
  `number` decimal(20,5) unsigned NOT NULL DEFAULT '0.00000' COMMENT '手数换算数量(非放大的)',
  `multiple` tinyint(4) unsigned NOT NULL DEFAULT '0' COMMENT '倍数',
  `origin_caution_money` decimal(20,8) unsigned NOT NULL DEFAULT '0.00000000' COMMENT '初始保证金',
  `caution_money` decimal(20,8) unsigned NOT NULL DEFAULT '0.00000000' COMMENT '当前可用保证金',
  `fact_profits` decimal(20,8) NOT NULL DEFAULT '0.00000000' COMMENT '最终盈亏',
  `trade_fee` decimal(20,8) unsigned NOT NULL DEFAULT '0.00000000' COMMENT '交易手续费',
  `overnight` decimal(20,4) unsigned NOT NULL DEFAULT '0.0000' COMMENT '隔夜费率,百分比',
  `overnight_money` decimal(20,8) unsigned NOT NULL DEFAULT '0.00000000' COMMENT '隔夜费金额',
  `status` tinyint(4) NOT NULL DEFAULT '0' COMMENT '交易状态:0.挂单中,1.交易中,2.平仓中,3.已平仓,4.已撤单',
  `settled` tinyint(4) NOT NULL DEFAULT '0' COMMENT '结算状态:0.未结算,1.已结算',
  `create_time` int(11) NOT NULL DEFAULT '0' COMMENT '下单时间',
  `transaction_time` decimal(20,8) NOT NULL DEFAULT '0.00000000' COMMENT '交易时间',
  `update_time` decimal(20,8) NOT NULL DEFAULT '0.00000000' COMMENT '价格刷新时间(毫秒级)',
  `handle_time` decimal(20,8) NOT NULL DEFAULT '0.00000000' COMMENT '平仓时间',
  `complete_time` decimal(20,8) NOT NULL DEFAULT '0.00000000' COMMENT '完成时间',
  `agent_path` varchar(2048) CHARACTER SET utf8 COLLATE utf8_unicode_ci DEFAULT '1' COMMENT '代理商关系',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB  DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;