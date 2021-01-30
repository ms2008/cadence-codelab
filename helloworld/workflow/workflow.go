package workflow

import (
	"time"

	"go.uber.org/cadence"
	"go.uber.org/zap"

	"github.com/venkat1109/cadence-codelab/helloworld/activity"
)

// TaskListName is the name of the task list
// that the decision / activity workers will
// poll for tasks
const TaskListName = "helloWorldTaskList"

// init registers the workflow with the cadence library
// the registration causes the client library to record
// the mapping from the function name to function pointer
// This mapping will be used during task dispatch
func init() {
	cadence.RegisterWorkflow(HelloWorld)
}

// HelloWorld is the implementation of helloworld workflow
// this workflow simply invokes an activity that prints a
// hello world message
func HelloWorld(ctx cadence.Context, name string) error {

	ao := cadence.ActivityOptions{
		ScheduleToStartTimeout: time.Minute,
		StartToCloseTimeout:    time.Minute,
		HeartbeatTimeout:       time.Second * 20,
	}

	ctx = cadence.WithActivityOptions(ctx, ao)

	logger := cadence.GetLogger(ctx)
	logger.Info("Workflow started")

	var helloworldResult string
	future := cadence.ExecuteActivity(ctx, activity.Helloworld, name)
	err := future.Get(ctx, &helloworldResult)
	if err != nil {
		logger.Error("Activity failed.", zap.Error(err))
		return err
	}

	logger.Info("Workflow completed.", zap.String("Result", helloworldResult))
	return nil
}
