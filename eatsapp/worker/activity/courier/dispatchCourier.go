package courier

import (
	"context"
	"net/http"
	"net/url"

	"go.uber.org/cadence"
	"go.uber.org/zap"
)

func init() {
	cadence.RegisterActivity(DispatchCourierActivity)
}

// DispatchCourierActivity implements the dispatch courier activity.
func DispatchCourierActivity(ctx context.Context, orderID string) (string, error) {
	logger := cadence.GetActivityLogger(ctx)
	activityInfo := cadence.GetActivityInfo(ctx)

	// register token with external service
	err := dispatch(orderID, string(activityInfo.TaskToken))
	if err != nil {
		logger.Info("Failed to send order to courier.", zap.Error(err))
		return "", err
	}

	// ErrActivityResultPending is returned from activity's execution to indicate the activity is not completed when it returns.
	// activity will be completed asynchronously when Client.CompleteActivity() is called.
	logger.Info("Successfully sent order to courier.", zap.String("Order ID", orderID))
	return "", cadence.ErrActivityResultPending
}

func dispatch(orderID string, taskToken string) error {
	formData := url.Values{}
	formData.Add("id", orderID)
	formData.Add("task_token", taskToken)

	url := "http://localhost:8090/courier"
	_, err := http.PostForm(url, formData)
	if err != nil {
		return err
	}

	return nil
}
