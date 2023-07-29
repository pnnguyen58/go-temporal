package infra

import (
	"context"
	"github.com/pnnguyen58/go-temporal/config"
	"go.temporal.io/sdk/client"
)

func NewTemporalClient(ctx context.Context) (client.Client, error) {
	return client.Dial(client.Options{
		HostPort: config.C.Server.TempoHost,
	})
}
