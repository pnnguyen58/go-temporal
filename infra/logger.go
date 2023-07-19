package infra

import (
	"context"
	"go.uber.org/zap"
)

func NewLogger(ctx context.Context, cfg *zap.Config) (*zap.Logger, error) {
	return zap.NewDevelopment()
	// TODO: fix to add config for logger
}
