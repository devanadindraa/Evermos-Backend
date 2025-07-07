package context

type contextKey string

const (
	tokenKey     contextKey = "token"
	requestIdKey contextKey = "requestId"
)
