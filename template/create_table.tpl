CREATE TABLE IF NOT EXISTS `${TABLE_NAME}`(
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '主键自增编号',
  `key` varchar(255) DEFAULT NULL COMMENT '缓存的唯一键',
  `value` longblob COMMENT '缓存的值',
  `timeout` bigint DEFAULT NULL COMMENT '超时时间',
  `create_time` datetime DEFAULT NULL COMMENT '创建记录的时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `key_UNIQUE` (`key`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='缓存数据表'