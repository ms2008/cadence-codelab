package eats

import (
	"time"

	"go.uber.org/cadence"
	"go.uber.org/zap"

	"github.com/venkat1109/cadence-codelab/eatsapp/worker/workflow/courier"
)

func deliverOrder(ctx cadence.Context, orderID string) error {
	cwo := cadence.ChildWorkflowOptions{
		WorkflowID:                   "DO_" + orderID,
		ExecutionStartToCloseTimeout: time.Minute * 30,
	}
	ctx = cadence.WithChildWorkflowOptions(ctx, cwo)
	err := cadence.ExecuteChildWorkflow(ctx, courier.DeliverOrderWorkflow, orderID).Get(ctx, nil)
	if err != nil {
		cadence.GetLogger(ctx).Error("DeliverOrder failed.", zap.Error(err))
		return err
	}

	return nil
}
