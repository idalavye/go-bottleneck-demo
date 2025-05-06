package api

import (
	"context"
	"runtime/trace"
)

// TraceRegion runs the given function within a trace region with the provided name.
func TraceRegion(ctx context.Context, regionName string, fn func()) {
	trace.WithRegion(ctx, regionName, fn)
}

// TraceRegionWithResult runs the given function within a trace region and returns its result.
func TraceRegionWithResult[T any](ctx context.Context, regionName string, fn func() T) T {
	var result T
	trace.WithRegion(ctx, regionName, func() {
		result = fn()
	})
	return result
}
