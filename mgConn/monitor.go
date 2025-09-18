package mgConn

import (
	"context"

	"github.com/d3v-friends/go-tools/fnLogger"
	"go.mongodb.org/mongo-driver/event"
)

func NewMonitor(loggers ...fnLogger.Logger) *event.CommandMonitor {
	var logger = fnLogger.NewLogger(fnLogger.LogLevelInfo)
	if len(loggers) == 1 {
		logger = loggers[0]
	}

	return &event.CommandMonitor{
		Started: func(ctx context.Context, ev *event.CommandStartedEvent) {
			logger.CtxTrace(
				ctx,
				ev.Command,
			)
		},
		Succeeded: func(ctx context.Context, ev *event.CommandSucceededEvent) {
			logger.CtxTrace(
				ctx,
				ev.Reply,
			)
		},
		Failed: func(ctx context.Context, ev *event.CommandFailedEvent) {
			logger.CtxError(
				ctx,
				ev.Failure,
			)
		},
	}
}
