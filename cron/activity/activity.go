package activity

import (
	"context"
	"time"

	"go.uber.org/cadence"
)

const (
	heartbeatInterval = 10 * time.Second
)

func init() {
	cadence.RegisterActivity(Cron)
}

// Cron implements the cron activity
func Cron(ctx context.Context) error {

	stopC := make(chan struct{})
	doneC := make(chan struct{})

	go doWork(stopC, doneC) // run cron task in background

	for {

		time.Sleep(heartbeatInterval)

		// Heartbeat to cadence with an optional
		// status message piggy backed with it
		// if the workflow is cancelled or terminated
		// cadence will set the context to Done
		cadence.RecordActivityHeartbeat(ctx, "status-report-to-workflow")

		if isDone(ctx, doneC) {
			// activity / workflow is cancelled, return now
			close(stopC)
			return ctx.Err()
		}
	}

	return nil
}

// isDone returns true if the activity is done or cancelled
// activity is done if the background task is done
// activity is cancelled if the context is marked as done
func isDone(ctx context.Context, doneC chan struct{}) bool {
	select {
	case <-ctx.Done():
		return true
	case <-doneC:
		return true
	default:
		return false
	}
}

// doWork is a dummy implementation of a
// cron job. In real world, this method
// could run a actual background task
func doWork(stopC chan struct{}, doneC chan struct{}) {
	select {
	// simulate two mins of work
	case <-time.After(time.Minute):
		close(doneC)
		return
	case <-stopC:
		// stop work and quit if we are asked to
		return
	}
}
