package utils_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/m-mizutani/ghnotify/pkg/utils"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/assert"
)

func TestMetricsMiddleware(t *testing.T) {
	utils.ResetMetrics()

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	metricsHandler := utils.MetricsMiddleware(handler)

	t.Run("metrics path should pass through", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/metrics", nil)
		assert.NoError(t, err)

		recorder := httptest.NewRecorder()
		metricsHandler.ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusOK, recorder.Code)
		assert.Equal(t, "OK", recorder.Body.String())

		metricName := "ghnotify_http_request_total"
		count, err := testutil.GatherAndCount(prometheus.DefaultGatherer, metricName)
		assert.Equal(t, 0, count)

		metricName = "ghnotify_http_status_total"
		count, err = testutil.GatherAndCount(prometheus.DefaultGatherer, metricName)
		assert.Equal(t, 0, count)

	})

	t.Run("numbers of responses", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/test", nil)
		assert.NoError(t, err)

		recorder := httptest.NewRecorder()
		metricsHandler.ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusOK, recorder.Code)
		assert.Equal(t, "OK", recorder.Body.String())

		count, err := testutil.GatherAndCount(prometheus.DefaultGatherer, "ghnotify_http_status_total")
		assert.Equal(t, 1, count)
	})
}

func TestNumberOfRequests(t *testing.T) {
	utils.ResetMetrics()

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	metricsHandler := utils.MetricsMiddleware(handler)

	methods := []string{http.MethodGet, http.MethodPost, http.MethodPut}
	for _, method := range methods {
		t.Run("method "+method, func(t *testing.T) {
			req, err := http.NewRequest(method, "/test", nil)
			assert.NoError(t, err)

			recorder := httptest.NewRecorder()
			metricsHandler.ServeHTTP(recorder, req)
			assert.Equal(t, http.StatusOK, recorder.Code)

			metricFamilies, err := prometheus.DefaultGatherer.Gather()
			assert.NoError(t, err)

			var found bool
			for _, mf := range metricFamilies {
				if *mf.Name == "ghnotify_http_request_total" {
					for _, m := range mf.Metric {
						for _, l := range m.Label {
							if l.GetName() == "method" && l.GetValue() == method {
								assert.Equal(t, float64(1), *m.Counter.Value)
								found = true
								break
							}
						}
					}
				}
			}
			assert.True(t, found, "Did not find counter for method %s", method)
		})
	}
}
