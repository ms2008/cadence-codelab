package courier

import (
	"context"

	"go.uber.org/cadence"
	"go.uber.org/zap"
)

func init() {
	cadence.RegisterActivity(DeliverOrderActivity)
}

// DeliverOrderActivity implements the devliver order activity.
func DeliverOrderActivity(ctx context.Context, orderID string) (string, error) {

	logger := cadence.GetActivityLogger(ctx)
	activityInfo := cadence.GetActivityInfo(ctx)

	// register token with external service
	err := deliver(orderID, string(activityInfo.TaskToken))
	if err != nil {
		logger.Info("Failed to dispatch corier order.", zap.Error(err))
		return "", err
	}

	return "", cadence.ErrActivityResultPending
}

func deliver(orderID string, taskToken string) error {
	url := "http://localhost:8090/courier?action=c_token&id=" + orderID + "&task_token=" + taskToken
	return sendPatch(url)
}
