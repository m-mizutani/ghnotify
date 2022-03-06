package notify

import (
	"encoding/json"

	"github.com/m-mizutani/ghnotify/pkg/domain/types"
	"github.com/m-mizutani/goerr"
	"github.com/slack-go/slack"
)

type SlackClient interface {
	Post(ctx *types.Context, msg *slack.WebhookMessage) error
}

type webhookClient struct {
	url string
}

func NewSlackWebhook(url string) *webhookClient {
	return &webhookClient{
		url: url,
	}
}

func (x *webhookClient) Post(ctx *types.Context, msg *slack.WebhookMessage) error {
	if err := slack.PostWebhookContext(ctx, x.url, msg); err != nil {
		raw, _ := json.Marshal(msg)
		return goerr.Wrap(err).With("body", string(raw))
	}
	return nil
}

type webhookMock struct {
	PostMock func(ctx *types.Context, msg *slack.WebhookMessage) error
}

func NewSlackWebhookMock() *webhookMock {
	return &webhookMock{}
}

func (x *webhookMock) Post(ctx *types.Context, msg *slack.WebhookMessage) error {
	return x.PostMock(ctx, msg)
}

type apiClient struct {
	api *slack.Client
}

func NewSlackAPI(token string) *apiClient {
	return &apiClient{
		api: slack.New(token),
	}
}

func (x *apiClient) Post(ctx *types.Context, msg *slack.WebhookMessage) error {
	if _, _, _, err := x.api.SendMessageContext(ctx, msg.Channel,
		slack.MsgOptionText(msg.Text, false),
		slack.MsgOptionAttachments(msg.Attachments...),
	); err != nil {
		raw, _ := json.Marshal(msg)
		return goerr.Wrap(err).With("body", string(raw))
	}

	return nil
}
