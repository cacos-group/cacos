package sql

const (
	CreateEventLogSql = "CREATE TABLE `%s` (\n`id` bigint(20) NOT NULL AUTO_INCREMENT,\n`trx_id` bigint(20) NOT NULL DEFAULT 0,\n`event_type` varchar(45) COLLATE utf8mb4_unicode_ci NOT NULL,\n`params` varchar(1000) COLLATE utf8mb4_unicode_ci NOT NULL,\n`major` varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '放置key的名字',\n`created_at` timestamp NOT NULL DEFAULT current_timestamp(),\nPRIMARY KEY (`id`),\nKEY `idx_trx_id` (`trx_id`),\nKEY `idx_major` (`major`)\n) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci"
)
