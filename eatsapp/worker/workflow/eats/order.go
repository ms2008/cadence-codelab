package eats

import (
	"time"

	"go.uber.org/cadence"
	"go.uber.org/zap"

	"github.com/venkat1109/cadence-codelab/eatsapp/worker/workflow/restaurant"
)

func placeRestaurantOrder(ctx cadence.Context, orderID string, items []string) (time.Duration, error) {
	execution := cadence.GetWorkflowInfo(ctx).WorkflowExecution

	cwo := cadence.ChildWorkflowOptions{
		// Do not specify WorkflowID if you want cadence to generate a unique ID for child execution
		WorkflowID:                   "PO_" + orderID,
		ExecutionStartToCloseTimeout: time.Minute * 30,
	}

	ctx = cadence.WithChildWorkflowOptions(ctx, cwo)

	var eta time.Duration
	err := cadence.ExecuteChildWorkflow(ctx, restaurant.PlaceOrderWorkflow, execution.RunID, orderID, items).Get(ctx, &eta)
	if err != nil {
		cadence.GetLogger(ctx).Error("PlaceOrder failed.", zap.Error(err))
		return 0, err
	}

	return eta, nil
}
