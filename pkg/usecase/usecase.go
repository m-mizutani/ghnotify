package usecase

import (
	"github.com/google/go-github/v43/github"
	"github.com/m-mizutani/ghnotify/pkg/domain/model"
	"github.com/m-mizutani/ghnotify/pkg/domain/types"
	"github.com/m-mizutani/ghnotify/pkg/infra"
)

type Usecase struct {
	config  *model.Config
	clients *infra.Clients
}

func New(cfg *model.Config, clients *infra.Clients) *Usecase {
	return &Usecase{
		config:  cfg,
		clients: clients,
	}
}

func (x *Usecase) ValidateWebhook(signature string, body []byte) error {
	if x.config.WebhookSecret == "" {
		return nil
	}

	if err := github.ValidateSignature(signature, body, []byte(x.config.WebhookSecret)); err != nil {
		return types.ErrInvalidWebhookRequest.Wrap(err)
	}

	return nil
}
