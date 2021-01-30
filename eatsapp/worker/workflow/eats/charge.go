package eats

import (
	"time"

	"github.com/venkat1109/cadence-codelab/eatsapp/worker/activity/eats"
	"go.uber.org/cadence"
	"go.uber.org/zap"
)

func chargeOrder(ctx cadence.Context, orderID string) error {

	ao := cadence.ActivityOptions{
		ScheduleToStartTimeout: time.Minute * 5,
		StartToCloseTimeout:    time.Minute * 15,
	}
	ctx = cadence.WithActivityOptions(ctx, ao)
	err := cadence.ExecuteActivity(ctx, eats.ChargeOrderActivity, orderID).Get(ctx, nil)
	if err != nil {
		cadence.GetLogger(ctx).Error("Failed to charge customer", zap.Error(err))
		return err
	}

	return nil
}
