package ctx

import (
	"context"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

var Ctx *context.Context

func Emit(event string, data ...any) {
	if Ctx != nil {
		runtime.EventsEmit(*Ctx, event, data...)
	}
}
