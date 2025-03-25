package mgConn

import (
	"context"
	"github.com/d3v-friends/go-tools/fnLogger"
	"go.mongodb.org/mongo-driver/event"
	"strings"
)

func NewMonitor(loggers ...fnLogger.Logger) *event.CommandMonitor {
	var logger fnLogger.Logger

	if len(loggers) == 1 {
		logger = loggers[0]
	} else {
		logger = fnLogger.NewLogger(fnLogger.LogLevelInfo)
	}

	return &event.CommandMonitor{
		Started: func(ctx context.Context, ev *event.CommandStartedEvent) {
			logger.CtxInfo(
				ctx,
				map[string]any{
					"command": strings.ReplaceAll(ev.Command.String(), "\n", ""),
				},
			)
		},
		Succeeded: func(ctx context.Context, ev *event.CommandSucceededEvent) {
			logger.CtxTrace(
				ctx,
				map[string]any{
					"reply": strings.ReplaceAll(ev.Reply.String(), "\n", ""),
				},
			)
		},
		Failed: func(ctx context.Context, ev *event.CommandFailedEvent) {
			logger.CtxError(
				ctx,
				map[string]any{
					"failure": ev.Failure,
				},
			)
		},
	}
}
