CREATE TABLE `event_logs` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `trx_id` bigint(20) NOT NULL DEFAULT 0 COMMENT '事务id',
  `payload` varbinary(255) NOT NULL DEFAULT '' COMMENT 'event',
  `major` char(32) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT 'key的md5值',
  `event_time` timestamp NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT current_timestamp(),
  PRIMARY KEY (`id`),
  KEY `idx_trx_id` (`trx_id`),
  KEY `idx_major` (`major`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE `wal_logs` (
                              `id` bigint(20) NOT NULL AUTO_INCREMENT,
                              `trx_id` bigint(20) NOT NULL DEFAULT 0 COMMENT '事务id',
                              `payload` varbinary(255) NOT NULL DEFAULT '' COMMENT 'event',
                              `major` char(32) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT 'key的md5值',
                              `event_time` timestamp NOT NULL,
                              `created_at` timestamp NOT NULL DEFAULT current_timestamp(),
                              PRIMARY KEY (`id`),
                              KEY `idx_trx_id` (`trx_id`),
                              KEY `idx_major` (`major`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;