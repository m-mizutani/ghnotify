package server

import (
	"net/http"

	"github.com/m-mizutani/ghnotify/pkg/utils"
)

func handleHealthCheckRequest() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte("OK")); err != nil {
			utils.Logger.Error("fail to write response", utils.ErrLog(err))
		}
	}
}
