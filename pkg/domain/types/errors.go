package types

import "github.com/m-mizutani/goerr"

var (
	ErrInvalidWebhookRequest = goerr.New("invalid webhook request")

	ErrInvalidConfig = goerr.New("invalid config")
)
