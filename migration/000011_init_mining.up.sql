DROP TABLE IF EXISTS `mining_list`;
CREATE TABLE `mining_list` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(255) NOT NULL,
  `price` decimal(10,2) NOT NULL,
  `sort` int(1) NOT NULL DEFAULT '1',
  `cycle` int(8) NOT NULL,
  `max_exchange` int(5) NOT NULL,
  `cl_sum` decimal(10,2) NOT NULL,
  `cl_h` decimal(10,2) NOT NULL,
  `buy_type` int(5) NOT NULL DEFAULT '1',
  `create_time` int(11) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=10 DEFAULT CHARSET=utf8 COMMENT="矿机表";

DROP TABLE IF EXISTS `my_mining`;
CREATE TABLE `my_mining` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `kj_id` int(11) NOT NULL,
  `user_id` int(11) NOT NULL,
  `prcie` decimal(10,2) NOT NULL,
  `zhouqi` int(5) NOT NULL,
  `ok_zhouqi` int(10) NOT NULL DEFAULT '0',
  `xs_shouyi` int(11) NOT NULL,
  `ok_shouyi` decimal(10,2) NOT NULL DEFAULT '0.00',
  `shouyi` decimal(10,2) NOT NULL,
  `addtime` int(11) NOT NULL,
  `update_time` int(11) NOT NULL,
  `status` int(1) NOT NULL,
  `sort` int(2) DEFAULT '1',
  `yest_shouyi` decimal(10,2) NOT NULL DEFAULT '0.00',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT="用户矿机";