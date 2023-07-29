package app

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/pnnguyen58/go-temporal/config"
	"github.com/pnnguyen58/go-temporal/core/workflows"
	"go.temporal.io/api/workflowservice/v1"
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
	temporalClient client.Client
	tempoConfig *config.TempoConfig
}

func NewExampleApp(ctx context.Context, logger *zap.Logger, tempoConfig map[string]config.TempoConfig) ExampleApp {
	app := &exampleApp{
		logger: logger,
	}
	app.setConfig(tempoConfig)
	err := app.register(ctx)
	if err != nil {
		logger.Fatal(err.Error())
	}
	return app
}

// Initial steps of new temporal client
func (ea *exampleApp) setConfig(tempoConfig map[string]config.TempoConfig) () {
	key := "tempo-example-app"
	// Set default tempo app config
	ea.tempoConfig = &config.TempoConfig{
		HostPort: config.C.Server.TempoHost,
		Namespace: &config.Namespace {
			Name: key,
			WorkflowExecutionRetentionPeriod: 1720*time.Hour,
		},
		Workflow: &config.Workflow{
			TaskQueueName: key,
			ExecutionTimeout: 300*time.Second,
			RunTimeout: 300*time.Second,
			TaskTimeout: 300*time.Second,
		},
	}
	if cfg, ok := tempoConfig[key]; ok {
		if cfg.HostPort != "" {
			ea.tempoConfig.HostPort = cfg.HostPort
		}
		if cfg.Namespace != nil {
			if cfg.Namespace.WorkflowExecutionRetentionPeriod > 0 {
				ea.tempoConfig.Namespace.WorkflowExecutionRetentionPeriod =
					cfg.Namespace.WorkflowExecutionRetentionPeriod * time.Second
			}
			if cfg.Namespace.Name != "" {
				ea.tempoConfig.Namespace.Name = cfg.Namespace.Name
			}
		}
		if cfg.Workflow != nil {
			ea.tempoConfig.Workflow.SearchAttributes = cfg.Workflow.SearchAttributes
			if cfg.Workflow.TaskQueueName != "" {
				ea.tempoConfig.Workflow.TaskQueueName = cfg.Workflow.TaskQueueName
			}
		}
	}
	return
}

func (ea *exampleApp) register(ctx context.Context) error {
	if cl, err := client.Dial(client.Options{
		HostPort: config.C.Server.TempoHost,
		Namespace: ea.tempoConfig.Namespace.Name,
	}); err != nil {
		return err
	} else {
		ea.temporalClient = cl
	}
	namespace, err := ea.temporalClient.WorkflowService().DescribeNamespace(ctx, &workflowservice.DescribeNamespaceRequest{
		Namespace: ea.tempoConfig.Namespace.Name,
	})
	if namespace != nil && err == nil {
		return nil
	}
	_, err = ea.temporalClient.WorkflowService().RegisterNamespace(ctx, &workflowservice.RegisterNamespaceRequest{
		Namespace:                        ea.tempoConfig.Namespace.Name,
		WorkflowExecutionRetentionPeriod: &ea.tempoConfig.Namespace.WorkflowExecutionRetentionPeriod,
	})
	return err
}

// CreateOrchestration create new example use case
func (ea *exampleApp) CreateOrchestration(ctx context.Context, req *example.ExampleCreateRequest) (*example.ExampleCreateResponse, error) {
	// Get task config
	taskQueueName := ea.tempoConfig.Workflow.TaskQueueName
	taskQueueID := uuid.New().String()
	taskTimeout := ea.tempoConfig.Workflow.TaskTimeout

	// Get workflow config
	attributes := ea.tempoConfig.Workflow.SearchAttributes
	executionTimeout := ea.tempoConfig.Workflow.ExecutionTimeout
	runTimeout := ea.tempoConfig.Workflow.RunTimeout

	workflowOptions := client.StartWorkflowOptions{
		ID:               taskQueueName + "_" + taskQueueID,
		TaskQueue:        taskQueueName,
		SearchAttributes: attributes,
		WorkflowExecutionTimeout: executionTimeout,
		WorkflowRunTimeout: runTimeout,
	}

	we, err := ea.temporalClient.ExecuteWorkflow(ctx, workflowOptions, workflows.CreateExampleWorkflow, req)
	if err != nil {
		ea.logger.Error("execute workflow failed")
		return nil, err
	}

	ctxWithTimeout, cancelHandler := context.WithTimeout(context.Background(), taskTimeout)
	defer cancelHandler()

	res := &example.ExampleCreateResponse{}
	err = we.Get(ctxWithTimeout, &res)
	if err != nil {
		return nil, err
	}
	ea.logger.Info(fmt.Sprintf("execute workflow ID: %v successfully", we.GetID()))
	return res, nil
}

// GetOrchestration get an example use case
func (ea *exampleApp) GetOrchestration(ctx context.Context, request *example.ExampleGetRequest) (*example.ExampleGetResponse, error) {
	return nil, nil
}

// GetAllOrchestration get all example use case
func (ea *exampleApp) GetAllOrchestration(ctx context.Context, request *example.ExampleGetAllRequest) (*example.ExampleGetAllResponse, error) {
	return nil, nil
}

// UpdateOrchestration update an example use case
func (ea *exampleApp) UpdateOrchestration(ctx context.Context, request *example.ExampleUpdateRequest) (*example.ExampleUpdateResponse, error) {
	return nil, nil
}

// DeleteOrchestration delete an example use case
func (ea *exampleApp) DeleteOrchestration(ctx context.Context, request *example.ExampleDeleteRequest) (*example.ExampleDeleteResponse, error) {
	return nil, nil
}