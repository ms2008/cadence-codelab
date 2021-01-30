package activity

import (
	"context"

	"go.uber.org/cadence"
)

// init registers the activity with the cadence library
func init() {
	cadence.RegisterActivity(Helloworld)
}

// HelloWorld is the implementation helloworld activity
func Helloworld(ctx context.Context, name string) (string, error) {
	logger := cadence.GetActivityLogger(ctx)
	logger.Info("helloworld activity started")
	return "Hello " + name + "!", nil
}
