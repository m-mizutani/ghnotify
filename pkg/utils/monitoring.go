package utils

import (
	"net/http"
	"net/http/httptest"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
)

var httpStatusCounterVec = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "ghnotify_http_status_total",
		Help: "Total number of HTTP responses, by status code.",
	},
	[]string{"code"},
)

func init() {
	prometheus.MustRegister(httpStatusCounterVec)
}

func MetricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/metrics" {
			next.ServeHTTP(w, r)
			return
		}

		recorder := httptest.NewRecorder()
		next.ServeHTTP(recorder, r)

		for k, v := range recorder.HeaderMap {
			w.Header()[k] = v
		}
		w.WriteHeader(recorder.Code)
		w.Write(recorder.Body.Bytes())

		httpStatus := strconv.Itoa(recorder.Code)
		httpStatusCounterVec.WithLabelValues(httpStatus).Inc()
	})
}
