package workflows

import (
	"github.com/pnnguyen58/go-temporal/core/activities"
	"go.temporal.io/sdk/workflow"
	"sdk/proto/example"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"go.temporal.io/sdk/testsuite"
)

type UnitTestSuite struct {
	suite.Suite
	testsuite.WorkflowTestSuite

	env *testsuite.TestWorkflowEnvironment
}

func (s *UnitTestSuite) SetupTest() {
	s.env = s.NewTestWorkflowEnvironment()
}

func (s *UnitTestSuite) AfterTest(suiteName, testName string) {
	s.env.AssertExpectations(s.T())
}

func TestUnitTestSuite(t *testing.T) {
	suite.Run(t, new(UnitTestSuite))
}

func (s *UnitTestSuite) Test_CreateExampleWorkflow_Success() {
	s.env.RegisterActivity(activities.CreateExample)

	s.env.ExecuteWorkflow(CreateExampleWorkflow, &example.ExampleCreateRequest{})

	s.True(s.env.IsWorkflowCompleted())
	s.NoError(s.env.GetWorkflowError())
}

func ProgressWorkflow(ctx workflow.Context, percent int) error {
	logger := workflow.GetLogger(ctx)

	err := workflow.SetQueryHandler(ctx, "getProgress", func(input []byte) (int, error) {
		return percent, nil
	})
	if err != nil {
		logger.Info("SetQueryHandler failed.", "Error", err)
		return err
	}

	for percent = 0; percent<100; percent++ {
		// Important! Use `workflow.Sleep()`, not `time.Sleep()`, because Temporal's
		// test environment doesn't stub out `time.Sleep()`.
		_ = workflow.Sleep(ctx, time.Second*1)
	}

	return nil
}

func (s *UnitTestSuite) Test_ProgressWorkflow() {
	value := 0

	// After 10 seconds plus padding, progress should be 10.
	// Note that `RegisterDelayedCallback()` doesn't actually make your test wait for 10 seconds!
	// Temporal's test framework advances time internally, so this test should take < 1 second.
	s.env.RegisterDelayedCallback(func() {
		res, err := s.env.QueryWorkflow("getProgress")
		s.NoError(err)
		err = res.Get(&value)
		s.NoError(err)
		s.Equal(10, value)
	}, time.Second*10+time.Millisecond*1)

	s.env.ExecuteWorkflow(ProgressWorkflow, 0)

	s.True(s.env.IsWorkflowCompleted())

	// Once the workflow is completed, progress should always be 100
	res, err := s.env.QueryWorkflow("getProgress")
	s.NoError(err)
	err = res.Get(&value)
	s.NoError(err)
	s.Equal(value, 100)
}

func (s *UnitTestSuite) Test_Workflow() {
	s.env = s.NewTestWorkflowEnvironment()
	s.env.RegisterActivity(activities.CreateExample)

	s.env.ExecuteWorkflow(CreateExampleWorkflow, &example.ExampleCreateRequest{})

	s.True(s.env.IsWorkflowCompleted())
	s.NoError(s.env.GetWorkflowError())

	var result example.ExampleCreateResponse
	_ = s.env.GetWorkflowResult(&result)
}
