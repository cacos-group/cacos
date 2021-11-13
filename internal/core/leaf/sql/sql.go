package sql

const (
	_updateSql = "UPDATE leaf_alloc SET max_id=max_id+step WHERE biz_tag=?"
	_selectSql = "SELECT max_id, step FROM leaf_alloc where biz_tag = ?"
	_insertSql = "INSERT INTO `leaf_alloc` (`biz_tag`, `max_id`, `step`, `description`) VALUES (?, ?, ?, ?) ON DUPLICATE KEY UPDATE `step` = VALUES(`step`);"
)
