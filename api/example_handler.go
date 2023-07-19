package main

import (
	"context"
	"github.com/pnnguyen58/go-temporal/core/app"
	"go.uber.org/zap"
	"sdk/proto/example"
)

type ExampleHandler interface {
	CreateExample(context.Context, *example.ExampleCreateRequest) (*example.ExampleCreateResponse, error)
}

func NewCommonHandler(logger *zap.Logger, app app.ExampleApp) ExampleHandler {
	return &ExampleController{logger: logger, app: app}
}

type ExampleController struct {
	example.UnimplementedExampleServiceServer
	logger *zap.Logger
	app app.ExampleApp
}

func (ec *ExampleController) CreateExample(ctx context.Context, req *example.ExampleCreateRequest) (*example.ExampleCreateResponse, error) {
	return ec.app.CreateOrchestration(ctx, req)
}

