package notify_test

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/m-mizutani/ghnotify/pkg/domain/types"
	"github.com/m-mizutani/ghnotify/pkg/infra/notify"
	"github.com/slack-go/slack"
	"github.com/stretchr/testify/require"
)

func TestSlackClient(t *testing.T) {
	url, ok := os.LookupEnv(types.EnvSlackWebhook)
	if !ok {
		t.Skip(types.EnvSlackWebhook + " is not set")
	}
	client := notify.NewSlackWebhook(url)

	msg := &slack.WebhookMessage{
		Text: "Hello, mizutani",
		Attachments: []slack.Attachment{
			{
				Color: "#2EB67D",

				Blocks: slack.Blocks{
					BlockSet: []slack.Block{
						slack.SectionBlock{
							Type: slack.MBTSection,
							Fields: []*slack.TextBlockObject{
								slack.NewTextBlockObject(slack.MarkdownType, "*Repo*: <https://github.com|m-mizutani/octovy>", false, false),
								slack.NewTextBlockObject(slack.MarkdownType, "*Issue*: <https://github.com|#41 New feature>", false, false),
								slack.NewTextBlockObject(slack.MarkdownType, "*Author*: <https://github.com|@someone>", false, false),
							},
						},
						slack.NewDividerBlock(),
						slack.SectionBlock{
							Type: slack.MBTSection,
							Text: &slack.TextBlockObject{
								Type: slack.MarkdownType,
								Text: "This implementation is slightly complicated. What do you think, mizutani?",
							},
						},
					},
				},
			},
		},
	}
	raw, err := json.Marshal(msg)
	require.NoError(t, err)
	t.Log(string(raw))
	require.NoError(t, client.Post(types.NewContext(), msg))
}
