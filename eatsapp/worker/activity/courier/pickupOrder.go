package courier

import (
	"context"

	"go.uber.org/cadence"
	"go.uber.org/zap"
)

func init() {
	cadence.RegisterActivity(PickUpOrderActivity)
}

// PickUpOrderActivity implements the pick-up order activity.
func PickUpOrderActivity(ctx context.Context, execution cadence.WorkflowExecution, orderID string) (string, error) {
	logger := cadence.GetActivityLogger(ctx)

	err := notifyRestaurant(execution, orderID)
	if err != nil {
		logger.Info("Failed to notify restaurant.", zap.Error(err))
		return "", err
	}

	activityInfo := cadence.GetActivityInfo(ctx)

	// register token with courier service
	err = pickup(orderID, string(activityInfo.TaskToken))
	if err != nil {
		logger.Info("Failed to dispatch courier order.", zap.Error(err))
		return "", err
	}

	return "", cadence.ErrActivityResultPending
}

func notifyRestaurant(execution cadence.WorkflowExecution, orderID string) error {
	url := "http://localhost:8090/restaurant?action=p_sig&id=" + orderID +
		"&workflow_id=" + execution.ID + "&run_id=" + execution.RunID
	return sendPatch(url)
}

func pickup(orderID string, taskToken string) error {
	url := "http://localhost:8090/courier?action=p_token&id=" + orderID + "&task_token=" + taskToken
	return sendPatch(url)
}
