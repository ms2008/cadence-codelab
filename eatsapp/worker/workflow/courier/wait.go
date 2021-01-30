package courier

import (
	"errors"

	"go.uber.org/cadence"
	"go.uber.org/zap"
)

func waitForRestaurantPickupConfirmation(ctx cadence.Context, signalName string) error {
	// listen for signal from restaurant
	s := cadence.NewSelector(ctx)

	var signalVal string
	signalChan := cadence.GetSignalChannel(ctx, signalName)
	s.AddReceive(signalChan, func(c cadence.Channel, more bool) {
		c.Receive(ctx, &signalVal)
		cadence.GetLogger(ctx).Info("Received order status signal!", zap.String("value", signalVal))
	})
	s.Select(ctx)

	if len(signalVal) > 0 && signalVal != "ORDER_PICKED_UP" {
		cadence.GetLogger(ctx).Error("Recieved wrong signal!", zap.String("Signal", signalVal))
		return errors.New("signalVal")
	}

	return nil
}
