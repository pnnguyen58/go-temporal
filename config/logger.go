package config

import (
	"context"
	"encoding/json"
	"go.uber.org/zap"
	"log"
)

func LoadLoggerConfig(ctx context.Context) *zap.Config {
	// TODO: create secret
	data, err := getSecret(ctx, SecretSuite{
		Region:     C.Secret.Region,
		SecretName: C.Secret.Logger,
	})
	if err != nil {
		log.Fatalln(err.Error())
	}
	cfg := &zap.Config{}
	err = json.Unmarshal(data, cfg)
	if err != nil {
		log.Fatalln(err.Error())
	}
	return cfg
}
