package eats

import (
	"errors"
	"time"

	"go.uber.org/cadence"
	"go.uber.org/zap"
)

func waitForRestaurant(ctx cadence.Context, signalName string, eta time.Duration) error {

	// wait until it is time to dispatch a courier
	s := cadence.NewSelector(ctx)

	ctx1, cancel := cadence.WithCancel(ctx)
	etaTimer := cadence.NewTimer(ctx1, eta)
	s.AddFuture(etaTimer, func(f cadence.Future) {
		f.Get(ctx, nil)
	})

	var signalVal string
	signalChan := cadence.GetSignalChannel(ctx, signalName)
	s.AddReceive(signalChan, func(c cadence.Channel, more bool) {
		c.Receive(ctx, &signalVal)
		cancel()
		cadence.GetLogger(ctx).Info("Received order status signal!", zap.String("value", signalVal))
	})
	s.Select(ctx)

	if len(signalVal) > 0 && signalVal != "ORDER_READY" {
		cadence.GetLogger(ctx).Error("Received non-ready signal!", zap.String("Signal", signalVal))
		return errors.New("signalVal")
	}

	return nil
}
