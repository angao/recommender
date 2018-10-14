CREATE TABLE IF NOT EXISTS `t_application` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(64) NOT NULL DEFAULT '' COMMENT '应用名称',
  `created` datetime DEFAULT NULL,
  `updated` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `t_application_unique` (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE IF NOT EXISTS `t_container_resource` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',
  `name` varchar(64) NOT NULL COMMENT '容器名称',
  `application_id` int(11) NOT NULL COMMENT '关联应用ID',
  `timeframe_id` int(11) DEFAULT NULL COMMENT '指定时间段ID',
  `cpu_limit` int(11) unsigned DEFAULT NULL,
  `memory_limit` int(11) unsigned DEFAULT NULL,
  `disk_read_io_limit` int(11) unsigned DEFAULT NULL,
  `disk_write_io_limit` int(11) unsigned DEFAULT NULL,
  `network_receive_io_limit` int(11) unsigned DEFAULT NULL,
  `network_transmit_io_limit` int(11) unsigned DEFAULT NULL,
  `created` datetime DEFAULT NULL,
  `updated` datetime DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE IF NOT EXISTS `t_timeframe` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(64) NOT NULL DEFAULT '' COMMENT '名称',
  `start` datetime NOT NULL COMMENT '开始时间',
  `end` datetime NOT NULL COMMENT '结束时间',
  `status` varchar(10) DEFAULT NULL COMMENT '执行状态',
  `description` varchar(255) DEFAULT NULL COMMENT '描述',
  `created` datetime DEFAULT NULL COMMENT '创建时间',
  `updated` datetime DEFAULT NULL COMMENT '修改时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `t_timeframe_unique` (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
