package restaurant

import (
	"context"
	"net/http"
	"net/url"

	"go.uber.org/cadence"
	"go.uber.org/zap"
)

func init() {
	cadence.RegisterActivity(PlaceOrderActivity)
}

// PlaceOrderActivity implements of send order activity.
func PlaceOrderActivity(ctx context.Context, wfRunID string, orderID string, items []string) (string, error) {

	logger := cadence.GetActivityLogger(ctx)
	activityInfo := cadence.GetActivityInfo(ctx)

	// register token with external service
	err := sendOrder(wfRunID, orderID, items, string(activityInfo.TaskToken))
	if err != nil {
		logger.Info("Failed to send order.", zap.Error(err))
		return "", err
	}

	// ErrActivityResultPending is returned from activity's execution to indicate the activity is not completed when it returns.
	// activity will be completed asynchronously when Client.CompleteActivity() is called.
	logger.Info("Successfully sent order.", zap.String("Order ID", orderID), zap.Strings("Items", items))
	return "", cadence.ErrActivityResultPending
}

func sendOrder(wfRunID string, orderID string, items []string, taskToken string) error {
	formData := url.Values{}
	formData.Add("id", orderID)
	formData.Add("workflow_id", orderID)
	formData.Add("run_id", wfRunID)
	formData.Add("task_token", taskToken)
	for _, item := range items {
		formData.Add("item", item)
	}
	url := "http://localhost:8090/restaurant"
	_, err := http.PostForm(url, formData)
	if err != nil {
		return err
	}
	return nil
}
