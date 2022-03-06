package server_test

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/go-github/v43/github"
	"github.com/m-mizutani/ghnotify/pkg/controller/server"
	"github.com/m-mizutani/ghnotify/pkg/domain/model"
	"github.com/m-mizutani/ghnotify/pkg/domain/types"
	"github.com/m-mizutani/ghnotify/pkg/infra"
	"github.com/m-mizutani/ghnotify/pkg/infra/notify"
	"github.com/m-mizutani/ghnotify/pkg/usecase"
	"github.com/m-mizutani/opac"
	"github.com/slack-go/slack"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func bind(secret string, data interface{}) (string, io.Reader) {
	raw, err := json.Marshal(data)
	if err != nil {
		panic(err.Error())
	}

	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(raw)
	signature := mac.Sum(nil)

	return "sha256=" + hex.EncodeToString(signature), bytes.NewReader(raw)
}

func TestWebhook(t *testing.T) {
	secret := "blue"
	var calledOPA, calledSlack int
	opaMock := opac.NewMock(func(in interface{}) (interface{}, error) {
		calledOPA++
		input, ok := in.(*model.RegoInput)
		require.True(t, ok)
		assert.Equal(t, "issue_comment", input.Name)

		event, ok := input.Event.(*github.IssueCommentEvent)
		require.True(t, ok)
		assert.Equal(t, "test-repo", event.Repo.GetName())
		return &model.RegoResult{
			Notify: []*model.Notify{
				{
					Channel: "#testing",
					Text:    "hello",
					Color:   "#123456",
					Body:    "body",
					Fields: []*model.NotifyField{
						{
							Name:  "color",
							Value: "blue",
							URL:   "https://example.com",
						},
					},
				},
			},
		}, nil
	})

	slackMock := notify.NewSlackWebhookMock()
	slackMock.PostMock = func(ctx *types.Context, msg *slack.WebhookMessage) error {
		calledSlack++
		return nil
	}

	clients := infra.New(infra.WithOPAC(opaMock), infra.WithSlack(slackMock))
	uc := usecase.New(&model.Config{WebhookSecret: secret}, clients)
	srv := server.New(uc)

	{
		w := httptest.NewRecorder()
		sig, body := bind(secret, github.IssueCommentEvent{
			Repo: &github.Repository{
				Name: github.String("test-repo"),
			},
		})
		r := httptest.NewRequest("POST", "/webhook/github", body)
		r.Header.Add("X-GitHub-Event", "issue_comment")
		r.Header.Add("X-Hub-Signature-256", sig)

		srv.ServeHTTP(w, r)
		assert.Equal(t, http.StatusOK, w.Result().StatusCode)
	}

	assert.Equal(t, 1, calledOPA)
	assert.Equal(t, 1, calledSlack)
}
