package main

import (
	"go.uber.org/cadence"

	"github.com/venkat1109/cadence-codelab/common"
	_ "github.com/venkat1109/cadence-codelab/cron/activity"
	_ "github.com/venkat1109/cadence-codelab/cron/workflow"
)

const (
	decisionTaskList   = "cron-decider"
	hostgroup1TaskList = "hostgroup-1"
	hostgroup2TaskList = "hostgroup-2"
)

func main() {
	runtime := common.NewRuntime()

	// Configure worker options.
	decisionWorkerOptions := cadence.WorkerOptions{
		MetricsScope:          runtime.Scope,
		Logger:                runtime.Logger,
		DisableActivityWorker: true,
	}

	// start decision worker on a separate task list
	runtime.StartWorkers(runtime.Config.DomainName, decisionTaskList, decisionWorkerOptions)

	// now start activity workers for hostgroup1 and hostgroup2
	activityWorkerOptions := cadence.WorkerOptions{
		MetricsScope:          runtime.Scope,
		Logger:                runtime.Logger,
		DisableWorkflowWorker: true,
	}
	runtime.StartWorkers(runtime.Config.DomainName, hostgroup1TaskList, activityWorkerOptions)
	runtime.StartWorkers(runtime.Config.DomainName, hostgroup2TaskList, activityWorkerOptions)

	select {}
}
