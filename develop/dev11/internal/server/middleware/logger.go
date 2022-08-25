package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

//Middleware ...
type Middleware struct {
	logger  *logrus.Logger
	Enabled bool
}

//NewMiddleware creates new middleware object
func NewMiddleware(logger *logrus.Logger, enabled bool) *Middleware {
	return &Middleware{
		logger:  logger,
		Enabled: enabled,
	}
}

//Log wraps a request and logs required info
func (m *Middleware) Log(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		next(w, r)

		if !m.Enabled {
			return
		}

		s := fmt.Sprintf("remote: %s | method: %s | URL: %s | time: %s\n",
			r.RemoteAddr, r.Method, r.URL.String(), time.Now())

		m.logger.Info(s)
	}
}
