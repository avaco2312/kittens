package middleware

import (
	"net/http"
)

// LimitHandler is middleware which limits the current number of active
// connections that this handler can sustain.
// Once the current connections equal the max http.StatusTooManyRequests is
// returned
type LimitHandler struct {
	connections chan struct{}
	handler     http.Handler
}

// NewLimitHandler creates a new instance of the LimitHandler for the
// given parameters.
func NewLimitHandler(connections int, next http.Handler) *LimitHandler {
	cons := make(chan struct{}, connections)
	return &LimitHandler{
		connections: cons,
		handler:     next,
	}
}

func (l *LimitHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	select {
	case l.connections <- struct{}{}:
		l.handler.ServeHTTP(rw, r)
		<-l.connections
	default:
		http.Error(rw, http.StatusText(http.StatusTooManyRequests), http.StatusTooManyRequests)
	}
}
