package workflows

import (
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.temporal.io/sdk/testsuite"

	"github.com/pnnguyen58/go-temporal/core/activities"
	"sdk/proto/example"
)

func Test_Compensations(t *testing.T) {

	tests := []struct {
		name string

		env *testsuite.TestWorkflowEnvironment

		acts   []any
		input  []any
		output []any
		outErr []error
	}{
		{
			name: "CreateExampleWorkflow",
			acts: []any{
				activities.CreateExample,
				activities.CreateExample,
			},
			input: []any{
				&example.ExampleCreateRequest{},
				&example.ExampleCreateRequest{},
			},
			output: []any{
				&example.ExampleCreateResponse{},
				&example.ExampleCreateResponse{},
			},
		},
	}

	for _, tt := range tests {
		// Set up the test suite and testing execution environment
		testSuite := &testsuite.WorkflowTestSuite{}
		tt.env = testSuite.NewTestWorkflowEnvironment()
		// Mock activity implementation
		for i, a := range tt.acts {
			tt.env.OnActivity(a, mock.Anything, tt.input[i]).Return(tt.output[i], tt.outErr[i]).Once()
		}
		t.Run(tt.name, func(t *testing.T) {
			_tt := tt
			t.Parallel()
			_tt.env.ExecuteWorkflow(CreateExampleWorkflow, tt.input)
			require.True(t, _tt.env.IsWorkflowCompleted())
			require.Error(t, _tt.env.GetWorkflowError())
		})
	}
}
