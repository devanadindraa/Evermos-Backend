package context

import (
	"context"
	"time"
)

type combinerCtx struct {
	parentCtx context.Context
	childCtx  context.Context
}

func newCombinerCtx(parentCtx context.Context, childCtx context.Context) context.Context {
	return &combinerCtx{
		parentCtx: parentCtx,
		childCtx:  childCtx,
	}
}

func (c *combinerCtx) Deadline() (deadline time.Time, ok bool) {
	deadline, ok = c.childCtx.Deadline()
	if !deadline.IsZero() {
		return deadline, ok
	}
	return c.parentCtx.Deadline()
}

func (c *combinerCtx) Done() <-chan struct{} {
	doneChan := c.childCtx.Done()
	if doneChan != nil {
		return doneChan
	}
	return c.parentCtx.Done()
}

func (c *combinerCtx) Err() error {
	err := c.childCtx.Err()
	if err != nil {
		return err
	}
	return c.parentCtx.Err()
}

func (c *combinerCtx) Value(key any) any {
	val := c.childCtx.Value(key)
	if val != nil {
		return val
	}
	return c.parentCtx.Value(key)
}
