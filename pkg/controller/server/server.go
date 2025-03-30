package server

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/m-mizutani/ghnotify/pkg/domain/types"
	"github.com/m-mizutani/ghnotify/pkg/usecase"
	"github.com/m-mizutani/ghnotify/pkg/utils"
	"github.com/m-mizutani/goerr"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"golang.org/x/exp/slog"
)

type Server struct {
	uc  *usecase.Usecase
	mux *chi.Mux
}

func New(uc *usecase.Usecase) *Server {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	r.Post("/webhook/github", utils.MetricsMiddleware(serveGitHubWebhook(uc)).ServeHTTP)

	r.Get("/health", handleHealthCheckRequest())

	r.Get("/metrics", utils.MetricsMiddleware(promhttp.Handler()).ServeHTTP)

	return &Server{
		uc:  uc,
		mux: r,
	}
}

func (x *Server) Listen(addr string) error {
	utils.Logger.Info("start listening", slog.String("addr", addr))
	server := &http.Server{Addr: addr, Handler: x.mux}

	errCh := make(chan error, 1)
	go func() {
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			errCh <- goerr.Wrap(err)
		}
	}()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGTERM, os.Interrupt)

	select {
	case err := <-errCh:
		return err

	case sig := <-sigCh:
		utils.Logger.With("signal", sig).Info("recv signal and shutting down...")
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := server.Shutdown(ctx); err != nil {
			return goerr.Wrap(err)
		}
	}

	return nil
}

func (x *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	x.mux.ServeHTTP(w, r)
}

func handleError(w http.ResponseWriter, err error) {
	code := http.StatusInternalServerError
	switch {
	case errors.Is(err, types.ErrInvalidWebhookRequest):
		code = http.StatusBadRequest
	}

	w.WriteHeader(code)
	if _, err := w.Write([]byte(err.Error())); err != nil {
		utils.Logger.Error("fail to write error response", utils.ErrLog(err))
	}
}

func toCtx(r *http.Request) *types.Context {
	ctx := r.Context()
	if c, ok := ctx.(*types.Context); ok {
		return c
	}
	return types.NewContext(types.WithCtx(ctx))
}

func serveGitHubWebhook(uc *usecase.Usecase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := handleGitHubWebhook(uc, r); err != nil {
			handleError(w, err)
			return
		}

		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte("ok")); err != nil {
			utils.Logger.Error("fail to write response", utils.ErrLog(err))
		}
	}
}
