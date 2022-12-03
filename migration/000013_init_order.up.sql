DROP TABLE IF EXISTS `ordergetpay`;
CREATE TABLE `ordergetpay` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `uid` int(11) DEFAULT NULL COMMENT '用户ID',
  `mer_order_no` varchar(50) DEFAULT NULL COMMENT '订单号',
  `addtime` int(11) DEFAULT NULL COMMENT '申请时间',
  `msgbs` decimal(15,0) DEFAULT '0',
  `addtime1` datetime DEFAULT NULL,
  `acc_no` varchar(32) DEFAULT NULL COMMENT '收款账号',
  `acc_name` varchar(255) DEFAULT NULL COMMENT '收款户名',
  `bank_code` varchar(255) DEFAULT NULL COMMENT '银行编码',
  `order_amount` decimal(20,2) DEFAULT NULL COMMENT '金额',
  `sendmsg` varchar(255) DEFAULT NULL COMMENT '报错反馈信息',
  `stat` tinyint(1) DEFAULT '0' COMMENT '0为申请中，1为放款成功。2为拒绝',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='墨西哥代付申请表';

DROP TABLE IF EXISTS `orderpay`;
CREATE TABLE `orderpay` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `userid` int(11) DEFAULT NULL COMMENT '用户ID',
  `order_amount` decimal(15,2) DEFAULT NULL,
  `mer_order_no` varchar(30) DEFAULT NULL COMMENT '订单号',
  `addtime` int(11) DEFAULT NULL COMMENT '下单时间',
  `status` tinyint(1) DEFAULT '0' COMMENT '0为未支付，1为支付',
  `xmtime` datetime DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;