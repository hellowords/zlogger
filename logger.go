package logger

import (
	"fmt"
	"sync"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	Drivers sync.Map
	Options
}

type Options struct {
	LogPath    string `json:"log_path" yaml:"log_path"`
	FileName   string `json:"file_name"`
	MaxSize    int    `json:"max_size"`
	MaxAge     int    `json:"max_age"`
	MaxBackups int    `json:"max_backups"`
}

func NewLogger(options Options) *Logger {
	l := &Logger{Options: options}
	l.Drivers.Store("default", l.getDefaultDriver())

	return l
}

// Store 使用定制日志文件名
func (l *Logger) Store(fileName string) *zap.Logger {
	if fileName == "" {
		if load, ok := l.Drivers.Load("default"); ok {
			return load.(*zap.Logger)
		}
		panic("l.Drivers.Load Error")
	}
	store, _ := l.Drivers.LoadOrStore(fileName, l.setDriver(fileName))
	return store.(*zap.Logger)
}

// Default 默认日志类
func (l *Logger) Default() *zap.Logger {
	load, _ := l.Drivers.Load("default")
	return load.(*zap.Logger)
}

func (l *Logger) getDefaultDriver() *zap.Logger {
	fileName := "app.log"
	if l.FileName != "" {
		fileName = l.FileName
	}
	defaultWriter := getLogWriter(fileName, l.MaxSize, l.MaxBackups, l.MaxAge)
	defaultEncoder := getEncoder()
	core := zapcore.NewCore(defaultEncoder, defaultWriter, zapcore.DebugLevel)

	return zap.New(core, zap.AddCaller())
}

func (l *Logger) setDriver(driver string) *zap.Logger {
	driverWriter := getLogWriter(fmt.Sprintf("%s.log", driver), l.MaxSize, l.MaxBackups, l.MaxAge)
	driverEncoder := getEncoder()
	core := zapcore.NewCore(driverEncoder, driverWriter, zap.DebugLevel)

	return zap.New(core, zap.AddCaller())
}

// Sync 退出时将所有drivers中的日志信息刷出缓冲区
func (l *Logger) Sync() {
	l.Drivers.Range(func(key, value interface{}) bool {
		if key == nil {
			return false
		}
		_ = value.(*zap.Logger).Sync()
		return true
	})
}

func getEncoder() zapcore.Encoder {
	config := zap.NewProductionEncoderConfig()
	config.EncodeTime = zapcore.ISO8601TimeEncoder
	config.TimeKey = "time"
	config.EncodeLevel = zapcore.CapitalLevelEncoder
	config.EncodeDuration = zapcore.SecondsDurationEncoder
	config.EncodeCaller = zapcore.ShortCallerEncoder
	return zapcore.NewJSONEncoder(config)
}

func getLogWriter(filename string, maxSize, maxBackup, maxAge int) zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   filename,
		MaxSize:    maxSize,
		MaxBackups: maxBackup,
		MaxAge:     maxAge,
	}
	return zapcore.AddSync(lumberJackLogger)
}
