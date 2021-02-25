package logger

import (
	"fmt"
	"testing"
	"time"

	"go.uber.org/zap"
)

func TestNewLogger(t *testing.T) {
	logger := NewLogger(Options{
		LogPath:    "",
		FileName:   "",
		MaxSize:    0,
		MaxAge:     0,
		MaxBackups: 0,
	})
	fmt.Printf("logger: %#v\n", logger)
}

func TestLogger_Default(t *testing.T) {
	logger := NewLogger(Options{
		LogPath:    "",
		FileName:   "",
		MaxSize:    10,
		MaxAge:     10,
		MaxBackups: 10,
	})

	defer logger.Sync()

	logger.Default().Info("Hello Drinke9!")
	logger.Default().Debug("this is debug", zap.Time("loadTime", time.Now()))
}

func TestLogger_Store(t *testing.T) {
	logger := NewLogger(Options{
		LogPath:    "",
		FileName:   "",
		MaxSize:    10,
		MaxAge:     10,
		MaxBackups: 20,
	})
	defer logger.Sync()

	logger.Store("storeTest").Debug("store test")
}
