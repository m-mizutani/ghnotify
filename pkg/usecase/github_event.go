package usecase

import (
	"encoding/json"
	"fmt"

	"github.com/google/go-github/v43/github"
	"github.com/m-mizutani/ghnotify/pkg/domain/model"
	"github.com/m-mizutani/ghnotify/pkg/domain/types"
	"github.com/m-mizutani/goerr"
	"github.com/slack-go/slack"
)

func (x *Usecase) HandleGitHubEvent(ctx *types.Context, eventType string, body []byte) error {
	eventData, err := github.ParseWebHook(eventType, body)
	if err != nil {
		return goerr.Wrap(err).With("body", string(body))
	}

	input := &model.RegoInput{
		Name:  eventType,
		Event: eventData,
	}
	var result model.RegoResult

	if err := x.clients.OPAC().Query(ctx, input, &result); err != nil {
		return goerr.Wrap(err)
	}
	ctx.Log().With("result", result).Debug("Policy evaluated")

	for _, notify := range result.Notify {
		msg := buildSlackMessage(notify, eventData)
		msg.Channel = notify.Channel

		rawMsg, _ := json.Marshal(msg)
		ctx.Log().With("msg", string(rawMsg)).Debug("notifying to slack")

		if err := x.clients.Slack().Post(ctx, msg); err != nil {
			return err
		}
	}

	return nil
}

func toBlock(field *model.NotifyField) *slack.TextBlockObject {
	text := fmt.Sprintf("*%s*: ", field.Name)
	valueStr := fmt.Sprintf("%+v", field.Value)

	if field.URL != "" {
		text += fmt.Sprintf("<%s|%s>", field.URL, valueStr)
	} else {
		text += valueStr
	}
	return slack.NewTextBlockObject(slack.MarkdownType, text, false, false)
}

func buildCommonFields(event interface{}) []*slack.TextBlockObject {
	body, err := json.Marshal(event)
	if err != nil {
		panic("failed to marshal github event: " + err.Error())
	}

	var data struct {
		Repository  *github.Repository  `json:"repository,omitempty"`
		Issue       *github.Issue       `json:"issue,omitempty"`
		PullRequest *github.PullRequest `json:"pull_request,omitempty"`
		Sender      *github.User        `json:"sender,omitempty"`
	}

	if err := json.Unmarshal(body, &data); err != nil {
		panic("failed to unmarshal github event: " + err.Error())
	}

	var blocks []*slack.TextBlockObject
	if data.Repository != nil {
		blocks = append(blocks, toBlock(&model.NotifyField{
			Name:  "Repo",
			Value: data.Repository.GetFullName(),
			URL:   data.Repository.GetHTMLURL(),
		}))
	}

	if data.Issue != nil {
		blocks = append(blocks, toBlock(&model.NotifyField{
			Name:  "Issue",
			Value: fmt.Sprintf("#%d %s", data.Issue.GetNumber(), data.Issue.GetTitle()),
			URL:   data.Issue.GetHTMLURL(),
		}))
	}

	if data.PullRequest != nil {
		blocks = append(blocks, toBlock(&model.NotifyField{
			Name:  "Pull Request",
			Value: fmt.Sprintf("#%d %s", data.PullRequest.GetNumber(), data.PullRequest.GetTitle()),
			URL:   data.PullRequest.GetHTMLURL(),
		}))
	}

	if data.Sender != nil {
		blocks = append(blocks, toBlock(&model.NotifyField{
			Name:  "Sender",
			Value: data.Sender.GetLogin(),
			URL:   data.Sender.GetHTMLURL(),
		}))
	}

	return blocks
}

func buildSlackMessage(notify *model.Notify, event interface{}) *slack.WebhookMessage {

	color := "#2EB67D"
	if notify.Color != "" {
		color = notify.Color
	}

	var blocks []slack.Block

	if fields := buildCommonFields(event); len(fields) > 0 {
		blocks = append(blocks, slack.NewSectionBlock(nil, fields, nil))
	}

	if notify.Body != "" {
		blocks = append(blocks,
			slack.NewDividerBlock(),
			slack.NewSectionBlock(&slack.TextBlockObject{
				Type: slack.PlainTextType,
				Text: notify.Body,
			}, nil, nil),
		)
	}

	var customFields []*slack.TextBlockObject
	for _, field := range notify.Fields {
		customFields = append(customFields, toBlock(field))
	}
	if len(customFields) > 0 {
		blocks = append(blocks,
			slack.NewDividerBlock(),
			slack.NewSectionBlock(nil, customFields, nil),
		)
	}

	msg := &slack.WebhookMessage{
		Text: notify.Text,
		Attachments: []slack.Attachment{
			{
				Color: color,
				Blocks: slack.Blocks{
					BlockSet: blocks,
				},
			},
		},
	}

	return msg
}
