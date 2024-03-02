package workflows

import (
	"go.temporal.io/sdk/workflow"
	"go.uber.org/multierr"
)

type Compensations struct {
	compensations []any
	arguments     [][]any
}

func (s *Compensations) AddCompensation(activity any, parameters ...any) {
	s.compensations = append(s.compensations, activity)
	s.arguments = append(s.arguments, parameters)
}

func (s *Compensations) Compensate(ctx workflow.Context, inParallel bool) error {
	var err error
	if !inParallel {
		// Compensate in Last-In-First-Out order, to undo in the reverse order that activities were applied.
		for i := len(s.compensations) - 1; i >= 0; i-- {
			errCompensation := workflow.ExecuteActivity(ctx, s.compensations[i], s.arguments[i]...).Get(ctx, nil)
			if errCompensation != nil {
				workflow.GetLogger(ctx).Error("Executing compensation failed", "Error", errCompensation)
			}
			err = multierr.Append(err, errCompensation)
		}
	} else {
		selector := workflow.NewSelector(ctx)
		for i := 0; i < len(s.compensations); i++ {
			execution := workflow.ExecuteActivity(ctx, s.compensations[i], s.arguments[i]...)
			selector.AddFuture(execution, func(f workflow.Future) {
				if errCompensation := f.Get(ctx, nil); errCompensation != nil {
					workflow.GetLogger(ctx).Error("Executing compensation failed", "Error", errCompensation)
					err = multierr.Append(err, errCompensation)
				}
			})
		}
		for range s.compensations {
			selector.Select(ctx)
		}

	}
	return err
}
