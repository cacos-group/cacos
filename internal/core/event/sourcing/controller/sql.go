package controller

import (
	"strings"
)

const (
	_addWalLogSql = "INSERT INTO `%s` (`payload`, `major`, `event_time`) VALUES(?,?,?)"
	//todo payload md5 去重
	_existsEventLogSql = "SELECT count(*) as count FROM `%s` WHERE `payload` = ?"
	_addEventLogSql    = "INSERT INTO %s (`trx_id`, `payload`, `major`, `event_time`) SELECT `trx_id`, `payload`, `major`, `event_time` FROM %s WHERE `trx_id` = ?"
	_delWalLogSql      = "DELETE FROM %s WHERE `trx_id` = ?"
)

func GenTableName(namespace string, sharding int) string {
	return "event_logs"
}

func GenWalLogTableName(eventLogTableName string) string {
	return strings.Replace(eventLogTableName, "event_log", "wal_log", 1)
}
