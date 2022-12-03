DROP TABLE IF EXISTS `agent_admin`;
CREATE TABLE `agent_admin` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `agent_id` int(11) NOT NULL,
  `username` varchar(50) NOT NULL DEFAULT '',
  `password` varchar(255) NOT NULL DEFAULT '',
  `role_id` tinyint(11) NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT="后台代理商管理员表";

DROP TABLE IF EXISTS `agent_role`;
CREATE TABLE `agent_role` (
  `id` int(10) NOT NULL AUTO_INCREMENT,
  `agent_id` int(10) NOT NULL,
  `name` varchar(50) NOT NULL,
  `is_super` tinyint(3) NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT="后台代理商角色";

DROP TABLE IF EXISTS `agent_role_permission`;
CREATE TABLE `agent_role_permission` (
  `id` int(10) NOT NULL AUTO_INCREMENT,
  `agent_id` int(10) NOT NULL,
  `role_id` int(11) NOT NULL DEFAULT '0',
  `module` varchar(50) NOT NULL DEFAULT '',
  `action` varchar(50) NOT NULL DEFAULT '',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT="后台代理商角色权限表";

DROP TABLE IF EXISTS `agent`;
CREATE TABLE `agent` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `user_id` int(11) unsigned NOT NULL COMMENT '用户id',
  `username` varchar(30) NOT NULL DEFAULT '' COMMENT '登录代理商后台的帐号',
  `password` varchar(255) NOT NULL DEFAULT '' COMMENT '登录代理商后台的密码',
  `parent_agent_id` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '父级代理商ID，0表示该用户是一级代理商',
  `level` tinyint(5) unsigned NOT NULL DEFAULT '0' COMMENT '代理商等级,0:超级管理员； 1：一级代理商；2:二级代理商；3:三级代理商；4:四级代理商',
  `agent_path` varchar(255) NOT NULL DEFAULT '0' COMMENT '代理商关系，用,拼接成字符串',
  `is_admin` tinyint(3) NOT NULL DEFAULT '0' COMMENT '是否为超级管理员',
  `is_lock` tinyint(3) NOT NULL DEFAULT '0' COMMENT '该代理商是否锁定',
  `is_addson` tinyint(3) NOT NULL DEFAULT '1' COMMENT '是否拥有开设下级代理商的权限',
  `pro_loss` decimal(10,2) NOT NULL DEFAULT '0.00' COMMENT '头寸比例',
  `pro_ser` decimal(10,2) NOT NULL DEFAULT '0.00' COMMENT '手续费比例',
  `status` tinyint(3) NOT NULL DEFAULT '1' COMMENT '1正常用户',
  `reg_time` int(11) NOT NULL DEFAULT '0' COMMENT '代理商注册时间',
  `lock_time` int(11) NOT NULL DEFAULT '0' COMMENT '代理商锁定时间',
  `money` decimal(20,8) DEFAULT '0.00000000' COMMENT '代理商帐户',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `username` (`username`) USING BTREE,
  KEY `user_id` (`user_id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT="后台代理商账号表";

DROP TABLE IF EXISTS `agent_log`;
CREATE TABLE `agent_log` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `agent_id` int(11) unsigned DEFAULT '0' COMMENT '代理商ID',
  `type` tinyint(3) unsigned DEFAULT '0' COMMENT '类型',
  `value` decimal(15,2) DEFAULT '0.00' COMMENT '操作值',
  `info` varchar(255) DEFAULT '' COMMENT '操作详情',
  `relate_id` int(11) unsigned DEFAULT '0' COMMENT '关联id',
  `add_time` int(11) unsigned DEFAULT '0' COMMENT '日志添加时间',
  `status` tinyint(3) unsigned DEFAULT '1' COMMENT '状态',
  PRIMARY KEY (`id`) USING BTREE,
  KEY `agent_id` (`agent_id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT="后台代理商日志表";

DROP TABLE IF EXISTS `agent_money_log`;
CREATE TABLE `agent_money_log` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `agent_id` int(11) NOT NULL DEFAULT '0' COMMENT '所属代理商',
  `type` tinyint(4) NOT NULL DEFAULT '1' COMMENT '类型。1.代理商头寸，2代理商手续费',
  `relate_id` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '关联id。比如订单id等',
  `before` decimal(20,8) NOT NULL DEFAULT '0.00000000' COMMENT '原余额',
  `change` decimal(20,8) NOT NULL DEFAULT '0.00000000' COMMENT '本次变动',
  `after` decimal(20,8) NOT NULL DEFAULT '0.00000000' COMMENT '现余额',
  `memo` varchar(100) NOT NULL DEFAULT '' COMMENT '备注',
  `created_time` int(11) NOT NULL COMMENT '变动时间',
  `son_user_id` int(11) NOT NULL DEFAULT '0' COMMENT '贡献收益的用户id',
  `status` tinyint(4) unsigned NOT NULL DEFAULT '0' COMMENT '是否体现',
  `legal_id` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '法币id',
  `updated_time` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '提现到账时间',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT="后台代理商提现日志表";

