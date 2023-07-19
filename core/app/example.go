package app

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/pnnguyen58/go-temporal/common"
	"github.com/pnnguyen58/go-temporal/core/workflows"
	"github.com/spf13/viper"
	"go.temporal.io/sdk/client"
	"go.uber.org/zap"
	"sdk/proto/example"
	"time"
)

type ExampleApp interface {
	CreateOrchestration(context.Context, *example.ExampleCreateRequest) (*example.ExampleCreateResponse, error)
	GetOrchestration(context.Context, *example.ExampleGetRequest) (*example.ExampleGetResponse, error)
	GetAllOrchestration(context.Context, *example.ExampleGetAllRequest) (*example.ExampleGetAllResponse, error)
	UpdateOrchestration(context.Context, *example.ExampleUpdateRequest) (*example.ExampleUpdateResponse, error)
	DeleteOrchestration(context.Context, *example.ExampleDeleteRequest) (*example.ExampleDeleteResponse, error)
}

type exampleApp struct {
	logger *zap.Logger
	temporal client.Client
}

func NewExampleApp(logger *zap.Logger) ExampleApp {
	return &exampleApp{logger: logger}
}

func (ea *exampleApp) CreateOrchestration(ctx context.Context, req *example.ExampleCreateRequest) (*example.ExampleCreateResponse, error) {
	attributes := map[string]interface{}{
		// "Invoker": "http",
	}

	taskQueueName := viper.GetString(common.EnvTempoTaskQueueName)

	workflowOptions := client.StartWorkflowOptions{
		ID:               taskQueueName + "_" + uuid.New().String(),
		TaskQueue:        taskQueueName,
		SearchAttributes: attributes,
	}

	we, err := ea.temporal.ExecuteWorkflow(ctx, workflowOptions, workflows.CreateExampleWorkflowV1, req)
	if err != nil {
		return nil, err
	}

	ctxWithTimeout, cancelHandler := context.WithTimeout(context.Background(), time.Second*60)
	defer cancelHandler()

	res := &example.ExampleCreateResponse{}
	err = we.Get(ctxWithTimeout, &res)
	if err != nil {
		return nil, err
	}
	ea.logger.Info(fmt.Sprintf("execute workflow ID: %v successfully", we.GetID()))
	return res, nil
}

func (ea exampleApp) GetOrchestration(ctx context.Context, request *example.ExampleGetRequest) (*example.ExampleGetResponse, error) {
	return nil, nil
}

func (ea exampleApp) GetAllOrchestration(ctx context.Context, request *example.ExampleGetAllRequest) (*example.ExampleGetAllResponse, error) {
	return nil, nil
}

func (ea exampleApp) UpdateOrchestration(ctx context.Context, request *example.ExampleUpdateRequest) (*example.ExampleUpdateResponse, error) {
	return nil, nil
}

func (ea exampleApp) DeleteOrchestration(ctx context.Context, request *example.ExampleDeleteRequest) (*example.ExampleDeleteResponse, error) {
	return nil, nil
}