package logger

import (
	"fmt"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Build 按配置构建最小可用日志实例。
func Build(level string) (*zap.Logger, error) {
	normalizedLevel := strings.ToLower(strings.TrimSpace(level))
	if normalizedLevel == "" {
		normalizedLevel = "info"
	}

	var parsedLevel zapcore.Level
	if err := parsedLevel.Set(normalizedLevel); err != nil {
		return nil, fmt.Errorf("parse log level: %w", err)
	}

	var cfg zap.Config
	if normalizedLevel == "debug" {
		cfg = zap.NewDevelopmentConfig()
	} else {
		cfg = zap.NewProductionConfig()
	}
	cfg.Level = zap.NewAtomicLevelAt(parsedLevel)

	return cfg.Build()
}
