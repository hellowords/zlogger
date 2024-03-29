package zlogger

import (
	"fmt"
	"os"
	"strings"
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
	Mode       string `json:"mode"`
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

// getDefaultDriver 获取默认logDriver
func (l *Logger) getDefaultDriver() *zap.Logger {
	fileName := "app"
	var path = l.LogPath
	if ok := strings.HasSuffix(l.LogPath, "/"); ok {
		path = l.LogPath[0 : len(l.LogPath)-1]
	}
	if l.FileName != "" {
		fileName = l.FileName
	}
	defaultWriter := getLogWriter(fmt.Sprintf("%s/%s.log", path, fileName), l.MaxSize, l.MaxBackups, l.MaxAge)
	defaultEncoder := getEncoder()
	writeSyncer := zapcore.NewMultiWriteSyncer(defaultWriter)
	if l.Options.Mode == "dev" {
		writeSyncer = zapcore.NewMultiWriteSyncer(defaultWriter, zapcore.AddSync(os.Stdout))
	}
	core := zapcore.NewCore(defaultEncoder, writeSyncer, zapcore.DebugLevel)

	return zap.New(core, zap.AddCaller())
}

// setDriver 设置新的logDriver
func (l *Logger) setDriver(name string) *zap.Logger {
	var path = l.LogPath
	if ok := strings.HasSuffix(l.LogPath, "/"); ok {
		path = l.LogPath[0 : len(l.LogPath)-1]
	}
	driverWriter := getLogWriter(fmt.Sprintf("%s/%s.log", path, name), l.MaxSize, l.MaxBackups, l.MaxAge)
	driverEncoder := getEncoder()
	writeSyncer := zapcore.NewMultiWriteSyncer(driverWriter)
	if l.Options.Mode == "dev" {
		writeSyncer = zapcore.NewMultiWriteSyncer(driverWriter, zapcore.AddSync(os.Stdout))
	}
	core := zapcore.NewCore(driverEncoder, writeSyncer, zap.DebugLevel)

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
