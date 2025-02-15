package middleware

import (
	"net/http"

	"github.com/SapolovichSV/durak/auth/internal/logger"
)

type MiddleWare struct {
	logger logger.Logger
}

func New(logger logger.Logger) *MiddleWare {
	return &MiddleWare{
		logger: logger,
	}
}
func (m *MiddleWare) Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		m.logger.Info(
			"Get HTTP Request",
			"method: ", r.Method,
			"pattern: ", r.Pattern)
		next.ServeHTTP(w, r)
		m.logger.Info(
			"Served HTTP Request")
	})
}
