package context

import (
	"context"
	"fmt"
	"runtime"
	"sync"
	"time"
)

// stopperCtx to avoid circular ctx call per-goroutine
type stopperCtx struct {
	ctx         context.Context
	isCalledMap sync.Map
}

func newStopperCtx(ctx context.Context) context.Context {
	return &stopperCtx{
		ctx: ctx,
	}
}

// getGID retrieves the goroutine ID (we'll use this to identify recursive calls).
func getGID() int {
	var buf [64]byte
	n := runtime.Stack(buf[:], false)
	idField := buf[10:n]
	var id int
	for _, b := range idField {
		if b >= '0' && b <= '9' {
			id = id*10 + int(b-'0')
		}
	}
	return id
}

func (s *stopperCtx) Value(key any) any {
	mapKey := fmt.Sprintf("Value-%d", getGID())
	isValueCalled, _ := s.isCalledMap.LoadOrStore(mapKey, false)

	if isValueCalled.(bool) {
		return nil
	}

	s.isCalledMap.Store(mapKey, true)
	defer s.isCalledMap.Store(mapKey, false)

	return s.ctx.Value(key)
}

func (s *stopperCtx) Done() <-chan struct{} {
	mapKey := fmt.Sprintf("Done-%d", getGID())
	isDoneCalled, _ := s.isCalledMap.LoadOrStore(mapKey, false)

	if isDoneCalled.(bool) {
		return nil
	}

	s.isCalledMap.Store(mapKey, true)
	defer s.isCalledMap.Store(mapKey, false)

	return s.ctx.Done()
}

func (s *stopperCtx) Err() error {
	mapKey := fmt.Sprintf("Err-%d", getGID())
	isErrCalled, _ := s.isCalledMap.LoadOrStore(mapKey, false)

	if isErrCalled.(bool) {
		return nil
	}

	s.isCalledMap.Store(mapKey, true)
	defer s.isCalledMap.Store(mapKey, false)

	return s.ctx.Err()
}

func (s *stopperCtx) Deadline() (deadline time.Time, ok bool) {
	mapKey := fmt.Sprintf("Deadline-%d", getGID())
	isDeadlineCalled, _ := s.isCalledMap.LoadOrStore(mapKey, false)

	if isDeadlineCalled.(bool) {
		return time.Time{}, false
	}

	s.isCalledMap.Store(mapKey, true)
	defer s.isCalledMap.Store(mapKey, false)

	return s.ctx.Deadline()
}
