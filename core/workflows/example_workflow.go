package workflows

import (
	"time"

	"sdk/proto/example"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"

	"github.com/pnnguyen58/go-temporal/core/activities"
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
		RetryPolicy:         &temporal.RetryPolicy{MaximumAttempts: 2},
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	var err error
	// This is how you log
	// workflows.GetLogger(ctx).Info("jobInput.Inputs", flowInput.Inputs)
	var compensations Compensations
	defer func() {
		if err != nil {
			// activity failed, and workflow context is canceled
			disconnectedCtx, _ := workflow.NewDisconnectedContext(ctx)
			err = compensations.Compensate(disconnectedCtx, true)
		}
	}()

	// Block for 1st activity of the workflow
	// execute the activity and add compensation in case failed
	{
		result := &example.ExampleCreateResponse{}
		compensations.AddCompensation(activities.CreateExampleCompensation, flowInput)
		err = workflow.ExecuteActivity(ctx, activities.CreateExample, flowInput).Get(ctx, result)
		if err != nil {
			return nil, err
		}
	}
	// Block for 2nd activity of the workflow
	// execute the activity and add compensation in case failed
	{
		result := &example.ExampleCreateResponse{}
		compensations.AddCompensation(activities.CreateExampleCompensation, flowInput)
		err = workflow.ExecuteActivity(ctx, activities.CreateExample, flowInput).Get(ctx, result)
		if err != nil {
			return nil, err
		}
	}
	// Add more activities and compensations
	workflow.GetLogger(ctx).Info("Workflow completed.")

	return nil, err
}
