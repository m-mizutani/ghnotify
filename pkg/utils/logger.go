package utils

import (
	"github.com/m-mizutani/zlog"
	"github.com/m-mizutani/zlog/filter"
)

var Logger = zlog.New()

func RenewLogger(options ...zlog.Option) {
	options = append(options,
		zlog.WithFilters(filter.Tag("secret")),
	)
	Logger = zlog.New(options...)
}
