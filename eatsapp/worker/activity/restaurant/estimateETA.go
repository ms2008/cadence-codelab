package restaurant

import (
	"context"
	"time"

	"go.uber.org/cadence"
	"go.uber.org/zap"
)

func init() {
	cadence.RegisterActivity(EstimateETAActivity)
}

// EstimateETAActivity implements the estimate eta activity.
func EstimateETAActivity(ctx context.Context, orderID string) (time.Duration, error) {
	time.Sleep(time.Second * 5)
	eta := time.Minute * 1
	cadence.GetActivityLogger(ctx).Info("Computed restaurant ready ETA", zap.Duration("ETA", eta))
	return eta, nil
}
