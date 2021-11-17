package strategy

import (
	"fmt"
	"strings"
)

const (
	_addEventLogSql    = "INSERT INTO `%s` (`event_type`, `params`, `major`, `event_time`) VALUES(?,?,?,?)"
	_existsEventLogSql = "SELECT count(*) as count FROM `%s` WHERE `event_type` = ? AND `params` = ?"
)

func GenTableName(namespace string, sharding int) string {
	return fmt.Sprintf("event_log_%s", namespace)
}

func GenRedoTableName(eventLogTableName string) string {
	return strings.Replace(eventLogTableName, "event_log", "redo_log", 1)
}
