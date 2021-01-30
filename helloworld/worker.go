package main

import (
	"go.uber.org/cadence"

	"github.com/venkat1109/cadence-codelab/common"
	"github.com/venkat1109/cadence-codelab/helloworld/workflow"
)

func main() {
	runtime := common.NewRuntime()
	// Configure worker options.
	workerOptions := cadence.WorkerOptions{
		MetricsScope: runtime.Scope,
		Logger:       runtime.Logger,
	}
	runtime.StartWorkers(runtime.Config.DomainName, workflow.TaskListName, workerOptions)
	select {}
}
