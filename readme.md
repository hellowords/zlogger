# zlogger 

使用 `github.com/uber-go/zap`日志库封装

## 特性

支持多通道日志配置

## 安装

`go get -u github.com/hellowords/zlogger`

## 使用


```go

logger := NewLogger(Options{
	Mode:       "dev" //只有在Mode为dev的时候会将日志同时输出到文件和控制台中
	LogPath:    "log/your/save/path",
	FileName:   "app",
	MaxSize:    10,
	MaxAge:     10,
	MaxBackups: 10,
})

defer logger.Sync()

logger.Default().Info("Hello Drinke9!")
logger.Default().Debug("this is debug", zap.Time("loadTime", time.Now()))
logger.Store("storeTest").Debug("store test")

```

