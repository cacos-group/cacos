package logger

import (
	"github.com/cacos-group/cacos/internal/conf"
	"github.com/cacos-group/cacos/pkg/zaplog"
)

func NewLog(cfg *conf.Config) zaplog.Logger {
	logCfg := cfg.Log
	return zaplog.New(zaplog.Config{LogName: logCfg.LogName})
}
