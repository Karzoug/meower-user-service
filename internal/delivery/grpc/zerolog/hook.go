package zerolog

import (
	"github.com/rs/zerolog"

	"github.com/Karzoug/meower-common-go/trace/otlp"
)

// TraceIDHook() returns a zerolog Hook function to log the trace ID attribute;
// it's required using the zerolog Ctx function when logging events to make the context available to the hook.
func TraceIDHook() zerolog.HookFunc {
	return func(e *zerolog.Event, level zerolog.Level, msg string) {
		e.Str("trace_id", otlp.GetTraceID(e.GetCtx()))
	}
}
