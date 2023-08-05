package server_test

import (
	"net/http/httptest"
	"testing"

	"github.com/m-mizutani/ghnotify/pkg/controller/server"
	"github.com/m-mizutani/ghnotify/pkg/usecase"
	"github.com/stretchr/testify/assert"
)

func TestHealthCheck(t *testing.T) {
	srv := server.New(&usecase.Usecase{})

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/health", nil)
	srv.ServeHTTP(w, r)

	assert.Equal(t, 200, w.Code)
}
