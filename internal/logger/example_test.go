package logger

import (
	"log/slog"
	"testing"
)

func TestLogger(t *testing.T) {
	logger := New(Config{
		Level:  slog.LevelInfo,
		Format: "json",
	})

	// Тестируем все методы
	logger.Debug("Debug message")
	logger.Debugf("Debug with format: %s", "test")
	logger.Debugw("Debug with fields", "key", "value")

	logger.Info("Info message")
	logger.Infof("Info with format: %d", 42)
	logger.Infow("Info with fields", "count", 10, "status", "ok")

	logger.Warn("Warning message")
	logger.Warnf("Warning with format: %v", true)
	logger.Warnw("Warning with fields", "alert", "high")

	logger.Error("Error message")
	logger.Errorf("Error with format: %s", "failed")
	logger.Errorw("Error with fields", "error_code", 500, "message", "internal error")
}
