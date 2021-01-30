package main

import (
	"time"

	"github.com/pborman/uuid"
	"go.uber.org/cadence"

	"github.com/venkat1109/cadence-codelab/common"
	"github.com/venkat1109/cadence-codelab/helloworld/workflow"
)

func main() {
	workflowOptions := cadence.StartWorkflowOptions{
		ID:                              "helloworld_" + uuid.New(),
		TaskList:                        workflow.TaskListName,
		ExecutionStartToCloseTimeout:    time.Minute,
		DecisionTaskStartToCloseTimeout: time.Minute,
	}
	// runtime is a helper routine for this codelab
	runtime := common.NewRuntime()
	runtime.StartWorkflow(workflowOptions, workflow.HelloWorld, "Cadence")
}
