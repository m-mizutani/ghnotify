package utils

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	httpRequestCounterVec = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "ghnotify_http_request_total",
			Help: "Total number of HTTP requests.",
		},
		[]string{"method"},
	)
	httpStatusCounterVec = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "ghnotify_http_status_total",
			Help: "Total number of HTTP responses, by status code.",
		},
		[]string{"code"},
	)
	httpResponseDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "ghnotify_http_response_duration_seconds",
			Help:    "Histogram of http response durations by status code.",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"code"},
	)
)

func init() {
	prometheus.MustRegister(httpRequestCounterVec)
	prometheus.MustRegister(httpStatusCounterVec)
	prometheus.MustRegister(httpResponseDuration)
}

func ResetMetrics() {
	httpRequestCounterVec.Reset()
	httpStatusCounterVec.Reset()
}

func MetricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/metrics" {
			next.ServeHTTP(w, r)
			return
		}

		startTime := time.Now()
		recorder := httptest.NewRecorder()
		next.ServeHTTP(recorder, r)
		duration := time.Since(startTime).Seconds()

		for k, v := range recorder.HeaderMap {
			w.Header()[k] = v
		}
		w.WriteHeader(recorder.Code)
		w.Write(recorder.Body.Bytes())

		httpStatus := strconv.Itoa(recorder.Code)
		httpRequestCounterVec.WithLabelValues(r.Method).Inc()
		httpStatusCounterVec.WithLabelValues(httpStatus).Inc()
		httpResponseDuration.WithLabelValues(httpStatus).Observe(duration)
	})
}
