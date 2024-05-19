package middleware

import (
	"bytes"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
)

type LoggerController struct {
	logger  *slog.Logger
	metrics MetricsInterface
}

type MetricsInterface interface {
	ObserveHTTPDuration(timeSince time.Time, path string, code int, method string)
}

func NewLoggerController(logger *slog.Logger, m MetricsInterface) *LoggerController {
	return &LoggerController{logger: logger, metrics: m}
}

type statusRecorder struct {
	http.ResponseWriter
	Status       int
	ResponseBody string
}

func (r *statusRecorder) WriteHeader(status int) {
	r.Status = status
	r.ResponseWriter.WriteHeader(status)
}

func (r *statusRecorder) Write(body []byte) (int, error) {
	r.ResponseBody = string(body)
	return r.ResponseWriter.Write(body)
}

func (lc *LoggerController) LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		recorder := &statusRecorder{
			ResponseWriter: w,
			Status:         200,
		}

		request, err := io.ReadAll(r.Body)
		if err != nil {
			lc.logger.Info("http_logger: read request body")
		}

		lc.logger.Info(string(request))

		r.Body = io.NopCloser(bytes.NewBuffer(request))

		start := time.Now()

		next.ServeHTTP(recorder, r)

		routeContext := chi.RouteContext(r.Context())
		path := strings.Join(routeContext.RoutePatterns, "")
		lc.logger.Info("Request handled", "path", path, "responce_status", recorder.Status, "method", r.Method)
		lc.metrics.ObserveHTTPDuration(start, path, recorder.Status, r.Method)
		lc.logger.Info("Response returned", "body", recorder.ResponseBody)
	})
}
