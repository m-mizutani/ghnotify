package types

import (
	"context"

	"github.com/m-mizutani/ghnotify/pkg/utils"
	"golang.org/x/exp/slog"
)

type Context struct {
	context.Context
	logger *slog.Logger
}

type ContextOption func(c *Context)

func NewContext(options ...ContextOption) *Context {
	ctx := &Context{
		Context: context.Background(),
		logger:  utils.Logger,
	}

	for _, opt := range options {
		opt(ctx)
	}
	return ctx
}

func WithCtx(ctx context.Context) ContextOption {
	return func(c *Context) {
		c.Context = ctx
	}
}

func WithLogger(logger *slog.Logger) ContextOption {
	return func(c *Context) {
		c.logger = logger
	}
}

func (x *Context) Log() *slog.Logger {
	return x.logger
}
