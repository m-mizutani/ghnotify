package types

import (
	"context"

	"github.com/m-mizutani/ghnotify/pkg/utils"
	"github.com/m-mizutani/zlog"
)

type Context struct {
	context.Context
	logger *zlog.Logger
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

func WithLogger(logger *zlog.Logger) ContextOption {
	return func(c *Context) {
		c.logger = logger
	}
}

func (x *Context) Log() *zlog.Logger {
	return x.logger
}
