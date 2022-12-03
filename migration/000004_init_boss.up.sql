DROP TABLE IF EXISTS `boss_account`;
CREATE TABLE `boss_account` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `uid` int(11) DEFAULT NULL COMMENT '用户id',
  `p_uid` int(11) DEFAULT NULL COMMENT '父亲id',
  `status` tinyint(1) DEFAULT '1' COMMENT '状态，1已申请 2:已激活',
  `invite_code` varchar(30) DEFAULT NULL COMMENT '邀请码',
  `total_invited` mediumint(9) DEFAULT '0' COMMENT '邀请成为boss用户数',
  `total_active` mediumint(9) DEFAULT '0' COMMENT '下线激活总人数',
  `total_profit` decimal(25,8) DEFAULT '0.00000000' COMMENT '总收益',
  `balance` decimal(25,8) DEFAULT '0.00000000' COMMENT '账户余额',
  `parent_id_array` varchar(255) DEFAULT NULL COMMENT 'json 祖宗uid数组 从父亲开始',
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 ROW_FORMAT=COMPACT COMMENT="老板账号表";

DROP TABLE IF EXISTS `boss_account_log`;
CREATE TABLE `boss_account_log` (
  `id` int(11) unsigned NOT NULL,
  `boss_account_id` int(11) DEFAULT NULL,
  `key` varchar(100) DEFAULT NULL,
  `amount` decimal(25,8) DEFAULT NULL,
  `description` varchar(255) DEFAULT NULL COMMENT '描述',
  `created_at` datetime DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8 ROW_FORMAT=COMPACT COMMENT="老板账号日志表";

DROP TABLE IF EXISTS `boss_account_relation`;
CREATE TABLE `boss_account_relation` (
  `id` int(11) unsigned NOT NULL,
  `uid` int(11) unsigned DEFAULT NULL,
  `relatived_uid` int(11) unsigned DEFAULT NULL,
  `type` tinyint(1) unsigned DEFAULT NULL COMMENT '1:parent 2:child',
  `generation` tinyint(2) unsigned DEFAULT NULL COMMENT '第几代关系，父子为第一代',
  `created_at` datetime DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8 ROW_FORMAT=COMPACT;

DROP TABLE IF EXISTS `boss_account_return_fail`;
CREATE TABLE `boss_account_return_fail` (
  `id` int(11) NOT NULL,
  `account_id` int(11) DEFAULT NULL,
  `fail_type` tinyint(1) DEFAULT NULL COMMENT '1:未激活，2:拉不够人',
  `relative_account_id` int(11) DEFAULT NULL,
  `generation` tinyint(4) DEFAULT NULL COMMENT '冗余字段，方便处理',
  `amount` decimal(25,8) DEFAULT NULL COMMENT '反佣数',
  `status` tinyint(1) DEFAULT '1' COMMENT '1:未处理，2:已处理',
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8 ROW_FORMAT=COMPACT;

