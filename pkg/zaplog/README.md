# Logger

## Usage

### Structured logging

```go
logger := log.NewStdLogger(os.Stdout)
// fields & valuer
logger = log.With(logger,
    "service.name", "hellworld",
    "service.version", "v1.0.0",
    "ts", log.DefaultTimestamp,
    "caller", log.DefaultCaller,
)
logger.Log(log.LevelInfo, "key", "value")

// helper
helper := log.NewHelper(logger)
helper.Log(log.LevelInfo, "key", "value")
helper.Info("info message")
helper.Infof("info %s", "message")
helper.Infow("key", "value")

// filter
log := log.NewHelper(log.NewFilter(logger,
	log.FilterLevel(LevelInfo),
	log.FilterKey("foo"),
	log.FilterValue("bar"),
	log.FilterFunc(customFilter),
))
log.Debug("debug zaplog")
log.Info("info zaplog")
log.Warn("warn zaplog")
log.Error("warn zaplog")
```
