package context

import (
	"context"
	"time"
)

// getValueCtx to get ctx value only suitable for goroutine
type getValueCtx struct {
	ctx context.Context
}

func NewGetValueCtx(ctx context.Context) context.Context {
	return &getValueCtx{
		ctx: ctx,
	}
}

func (getValueCtx) Deadline() (deadline time.Time, ok bool) {
	return
}

func (getValueCtx) Done() <-chan struct{} {
	return nil
}

func (getValueCtx) Err() error {
	return nil
}

func (g getValueCtx) Value(key any) any {
	return g.ctx.Value(key)
}
