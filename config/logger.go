package config

import "C"
import (
	"context"
	"encoding/json"
	"go.uber.org/zap"
	"log"
)

func LoadLoggerConfig(ctx context.Context) (*zap.Config, error) {
	data, err := getSecret(ctx, SecretSuite{
		Region:     C.Secret.Region,
		SecretName: C.Secret.Logger,
	})
	if err != nil {
		log.Println(err.Error())
	}
	cfg := &zap.Config{}
	err = json.Unmarshal(data, cfg)
	if err != nil {
		log.Println(err.Error())
	}
	return cfg, err
}
