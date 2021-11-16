package strategy

import "fmt"

const (
	_addEventLogSql    = "INSERT INTO `%s` (`event_type`, `params`, `major`) VALUES(?,?,?)"
	_existsEventLogSql = "SELECT count(*) as count FROM `%s` WHERE `event_type` = ? AND `params` = ?"
)

func GenTableName(namespace string, sharding int) string {
	return fmt.Sprintf("event_log_%s", namespace)
}
