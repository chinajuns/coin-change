DROP TABLE IF EXISTS `lh_bank_account`;
CREATE TABLE `lh_bank_account` (
  `id` int(11) NOT NULL,
  `uid` int(11) DEFAULT NULL,
  `p_uid` int(11) DEFAULT NULL,
  `df_balance` decimal(25,8) DEFAULT '0.00000000',
  `usdt_balance` decimal(25,8) DEFAULT '0.00000000',
  `total_profit` decimal(25,8) DEFAULT '0.00000000',
  `total_deposit_amount` decimal(25,8) DEFAULT NULL,
  `team_deposit_amount` decimal(25,8) DEFAULT NULL COMMENT '团队存款总量',
  `m_level` tinyint(4) DEFAULT '0' COMMENT '直推级别',
  `vip_level` tinyint(4) DEFAULT '0' COMMENT '团队级别',
  `status` tinyint(1) DEFAULT '1' COMMENT '1:正常',
  `updated_at` datetime DEFAULT NULL,
  `created_at` datetime DEFAULT NULL,
  `vip_log` varchar(255) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8 ROW_FORMAT=COMPACT;

DROP TABLE IF EXISTS `lh_bank_account_log`;
CREATE TABLE `lh_bank_account_log` (
  `id` int(11) NOT NULL,
  `account_id` int(11) DEFAULT NULL,
  `type` tinyint(1) DEFAULT NULL COMMENT '1:usdt账户 2:df-one',
  `amount` decimal(25,8) DEFAULT NULL COMMENT '数量',
  `description` varchar(100) DEFAULT NULL,
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8 ROW_FORMAT=COMPACT;

DROP TABLE IF EXISTS `lh_bank_team_member`;
CREATE TABLE `lh_bank_team_member` (
  `id` int(11) NOT NULL,
  `uid` int(11) DEFAULT NULL,
  `leader_uid` int(11) DEFAULT NULL,
  `generation` int(11) DEFAULT NULL COMMENT '第几代团队',
  `created_at` datetime DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8 ROW_FORMAT=COMPACT;