package server

import (
	"io"
	"net/http"

	"github.com/m-mizutani/ghnotify/pkg/domain/types"
	"github.com/m-mizutani/ghnotify/pkg/usecase"
	"github.com/m-mizutani/goerr"
)

func handleGitHubWebhook(uc *usecase.Usecase, r *http.Request) error {
	ctx := toCtx(r)

	eventType := r.Header.Get("X-GitHub-Event")
	if eventType == "" {
		return goerr.Wrap(types.ErrInvalidWebhookRequest)
	}
	eventBody, err := io.ReadAll(r.Body)
	if err != nil {
		return goerr.Wrap(err)
	}

	signature := r.Header.Get("X-Hub-Signature-256")
	if err := uc.ValidateWebhook(signature, eventBody); err != nil {
		return err
	}

	if err := uc.HandleGitHubEvent(ctx, eventType, eventBody); err != nil {
		return err
	}

	return nil
}
