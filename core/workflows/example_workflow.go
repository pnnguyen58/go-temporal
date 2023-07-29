package workflows

import (
	"github.com/pnnguyen58/go-temporal/core/activities"
	"go.temporal.io/sdk/workflow"
	"go.uber.org/multierr"
	"sdk/proto/example"
	"time"
)

// CreateExampleWorkflow workflows definition
func CreateExampleWorkflow(ctx workflow.Context, flowInput *example.ExampleCreateRequest) (*example.ExampleCreateResponse, error) {
	// Workflow has to check input valid or not
	//inputErr := flowInput.CheckValid()
	//if inputErr != nil {
	//	return nil,
	//		temporal.NewNonRetryableApplicationError("Invalid flow input", common.ErrInvalidInput, inputErr, nil)
	//}

	ao := workflow.ActivityOptions{
		StartToCloseTimeout: 30 * time.Second,
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	// This is how you log
	// workflows.GetLogger(ctx).Info("jobInput.Inputs", flowInput.Inputs)

	result := &example.ExampleCreateResponse{}
	err := workflow.ExecuteActivity(ctx, activities.CreateExample, flowInput).Get(ctx, result)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			errCompensation := workflow.ExecuteActivity(ctx, activities.CreateExampleCompensation, flowInput).
				Get(ctx, nil)
			err = multierr.Append(err, errCompensation)
		}
	}()
	workflow.GetLogger(ctx).Info("Workflow completed.")

	return result, err
}