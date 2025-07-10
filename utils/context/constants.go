package context

type ctxKey string

const (
	tokenKey     ctxKey = "token"
	RequestIdKey ctxKey = "requestId"
	FiberCtxKey  ctxKey = "fiberCtx"
)
