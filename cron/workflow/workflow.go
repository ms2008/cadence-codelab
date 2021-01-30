package workflow

import (
	"time"

	"go.uber.org/cadence"
	"go.uber.org/zap"

	"github.com/venkat1109/cadence-codelab/cron/activity"
)

type (
	CronSchedule struct {
		Count      int           // total number of jobs to schedule
		Frequency  time.Duration // frequency at which to schedule jobs
		Hostgroups []string      // schedule a job for each one of these hostgroup
	}
)

const maxJobsPerLoop = 1

func init() {
	cadence.RegisterWorkflow(Cron)
}

// Cron implements the Cron workflow
func Cron(ctx cadence.Context, schedule *CronSchedule) error {

	cadence.GetLogger(ctx).Info("Cron started", zap.Int("Count", schedule.Count),
		zap.Duration("frequency", schedule.Frequency), zap.Strings("groups", schedule.Hostgroups))

	activityCtx := cadence.WithActivityOptions(ctx, cadence.ActivityOptions{
		ScheduleToStartTimeout: time.Minute,
		StartToCloseTimeout:    10 * time.Minute,
		HeartbeatTimeout:       time.Minute,
	})

	return runScheduler(ctx, activityCtx, schedule)
}

// runScheduler runs the cron scheduler
func runScheduler(ctx cadence.Context, activityCtx cadence.Context, schedule *CronSchedule) error {

	var loopCount int

	for schedule.Count > 0 && loopCount < maxJobsPerLoop {

		futures := make([]cadence.Future, len(schedule.Hostgroups))

		for i, hg := range schedule.Hostgroups {
			// create a child context to route this activity
			// task to a specific task list
			childCtx := cadence.WithTaskList(activityCtx, hg)
			futures[i] = cadence.ExecuteActivity(childCtx, activity.Cron)
		}

		cadence.GetLogger(ctx).Info("Waiting for activities to complete")

		for _, future := range futures {
			err := future.Get(ctx, nil)
			if err != nil {
				cadence.GetLogger(ctx).Error("cron job failed", zap.Error(err))
			}
		}

		loopCount++
		schedule.Count--
		time.Sleep(schedule.Frequency)
	}

	if schedule.Count == 0 {
		cadence.GetLogger(ctx).Info("Cron workflow completed")
		return nil
	}

	// ContinueAsNew workflow to limit the history size
	ctx = cadence.WithExecutionStartToCloseTimeout(ctx, 24*time.Hour)
	ctx = cadence.WithWorkflowTaskStartToCloseTimeout(ctx, 20*time.Minute)
	return cadence.NewContinueAsNewError(ctx, Cron, schedule)
}
