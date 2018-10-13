-- ----------------------------
-- Table structure for t_application
-- ----------------------------
CREATE TABLE IF NOT EXISTS `t_application` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(64) NOT NULL,
  `namespace` varchar(64) NOT NULL,
  `created` datetime DEFAULT NULL,
  `updated` datetime DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for t_container_resource
-- ----------------------------
CREATE TABLE IF NOT EXISTS `t_container_resource` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',
  `application_id` int(11) NOT NULL,
  `name` varchar(64) NOT NULL COMMENT '应用名称',
  `cpu_request` float(255,2) unsigned DEFAULT NULL,
  `cpu_limit` float(255,2) unsigned DEFAULT NULL,
  `memory_request` float(255,2) unsigned DEFAULT NULL,
  `memory_limit` float(255,2) unsigned DEFAULT NULL,
  `disk_read_io_request` float(255,2) unsigned DEFAULT NULL,
  `disk_read_io_limit` float(255,2) unsigned DEFAULT NULL,
  `disk_write_io_request` float(255,2) unsigned DEFAULT NULL,
  `disk_write_io_limit` float(255,2) unsigned DEFAULT NULL,
  `network_io_request` float(255,2) unsigned DEFAULT NULL,
  `network_io_limit` float(255,2) unsigned DEFAULT NULL,
  `created` datetime DEFAULT NULL,
  `updated` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_application_name` (`name`,`application_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
