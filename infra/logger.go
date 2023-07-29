package infra

import (
	"context"
	"github.com/pnnguyen58/go-temporal/config"
	"go.uber.org/zap"
)

func NewLogger(ctx context.Context) (*zap.Logger, error) {
	return zap.NewDevelopment()
	// TODO: fix to add config for logger
	cfg := config.LoadLoggerConfig(ctx)
	logger := zap.Must(cfg.Build())
	defer logger.Sync()

	logger.Info("logger construction succeeded")
	return logger, nil
}
