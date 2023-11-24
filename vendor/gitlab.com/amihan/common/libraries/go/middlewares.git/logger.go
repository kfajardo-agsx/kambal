package middlewares

import (
	"net/http"
	"time"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

type loggerMiddleware struct {
	logger *log.Logger
	next   http.Handler
}

// Logger adds logging functions to http calls using the provided logrus logger.
// If the logger is nil, it will use logrus standard logger.
func Logger(logger *log.Logger) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return &loggerMiddleware{
			logger: logger,
			next:   next,
		}
	}
}

// ServeHTTP runs the logger middleware to log requests and responses. For security, the only the method, path,
// header and query counts of the request is logged. For the response, only the duration of the entire entire
// request is logged. (there is also no way of logging the response code/body without hijacking)
func (m *loggerMiddleware) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	method := req.Method
	path := req.URL.Path
	headerCount := len(req.Header)
	queryCount := len(req.URL.Query())
	m.logger.Infof("Request :: [%s] %s :: Header Count: %d :: Query Count: %d", method, path, headerCount, queryCount)
	start := time.Now()
	m.next.ServeHTTP(res, req)
	end := time.Now()
	m.logger.Infof("Response :: [%s] %s :: Duration: %v", method, path, end.Sub(start))
}
