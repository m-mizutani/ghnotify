package cmd

import (
	"net/http"
	"strings"

	"github.com/m-mizutani/ghnotify/pkg/controller/cmd/option"
	"github.com/m-mizutani/ghnotify/pkg/domain/model"
	"github.com/m-mizutani/ghnotify/pkg/domain/types"
	"github.com/m-mizutani/ghnotify/pkg/infra"
	"github.com/m-mizutani/ghnotify/pkg/infra/notify"
	"github.com/m-mizutani/ghnotify/pkg/usecase"
	"github.com/m-mizutani/ghnotify/pkg/utils"
	"github.com/m-mizutani/goerr"
	"github.com/m-mizutani/opac"
	"github.com/urfave/cli/v2"
	"google.golang.org/api/idtoken"
)

type globalConfig struct {
	SlackWebhookURL     string
	SlackAPIToken       string `masq:"secret"`
	GitHubWebhookSecret string `masq:"secret"`

	LocalPolicy   string
	LocalPackage  string
	RemoteURL     string
	RemoteHeaders []string

	remoteHeaders cli.StringSlice

	remoteIDTokenClient bool
}

type idTokenClient struct{}

func (x *idTokenClient) Do(req *http.Request) (*http.Response, error) {
	client, err := idtoken.NewClient(req.Context(), req.URL.String())
	if err != nil {
		return nil, goerr.Wrap(err, "failed idtoken.NewClient for GCP IAP").With("req", req)
	}

	return client.Do(req)
}

func (x *globalConfig) newUsecase() (*usecase.Usecase, error) {
	cfg := &model.Config{
		WebhookSecret: x.GitHubWebhookSecret,
	}

	// Configure slack client
	var slackClient notify.SlackClient
	if x.SlackWebhookURL != "" {
		slackClient = notify.NewSlackWebhook(x.SlackWebhookURL)
	} else if x.SlackAPIToken != "" {
		slackClient = notify.NewSlackAPI(x.SlackAPIToken)
	} else {
		return nil, goerr.Wrap(types.ErrInvalidConfig, "no slack config")
	}

	// Configure policy client
	var opacClient opac.Client
	if x.LocalPolicy != "" {
		options := []opac.LocalOption{opac.WithDir(x.LocalPolicy)}
		if x.LocalPackage != "" {
			options = append(options, opac.WithPackage(x.LocalPackage))
		}
		client, err := opac.NewLocal(options...)
		if err != nil {
			return nil, goerr.Wrap(err)
		}
		opacClient = client
	} else if x.RemoteURL != "" {
		options := []opac.RemoteOption{}
		for _, hdr := range x.RemoteHeaders {
			parts := strings.SplitN(hdr, ":", 2)
			if len(parts) != 2 {
				return nil, goerr.Wrap(types.ErrInvalidConfig, "invalid header format").With("hdr", hdr)
			}
			options = append(options, opac.WithHTTPHeader(
				strings.TrimSpace(parts[0]),
				strings.TrimSpace(parts[1]),
			))
		}

		if x.remoteIDTokenClient {
			options = append(options, opac.WithHTTPClient(&idTokenClient{}))
		}

		client, err := opac.NewRemote(x.RemoteURL, options...)
		if err != nil {
			return nil, goerr.Wrap(err)
		}
		opacClient = client
	} else {
		return nil, goerr.Wrap(types.ErrInvalidConfig, "no policy config")
	}

	clients := infra.New(
		infra.WithSlack(slackClient),
		infra.WithOPAC(opacClient),
	)
	return usecase.New(cfg, clients), nil
}

func Run(argv []string) error {
	var (
		cfg globalConfig

		logLevel  option.LogLevel
		logFormat option.LogFormat
		logOutput option.LogOutput
	)

	app := &cli.App{
		Name:  "ghnotify",
		Usage: "General GitHub event notification tool to Slack with Open Policy Agent and Rego",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "local-policy",
				Usage:       "Rego local policy directory",
				EnvVars:     []string{"GHNOTIFY_LOCAL_POLICY"},
				Destination: &cfg.LocalPolicy,
			},
			&cli.StringFlag{
				Name:        "local-package",
				Usage:       "Rego local policy package",
				EnvVars:     []string{"GHNOTIFY_LOCAL_PACKAGE"},
				Destination: &cfg.LocalPackage,
				Value:       "github.notify",
			},
			&cli.StringFlag{
				Name:        "remote-url",
				Usage:       "OPA server URL",
				EnvVars:     []string{"GHNOTIFY_REMOTE_URL"},
				Destination: &cfg.RemoteURL,
			},
			&cli.StringSliceFlag{
				Name:        "remote-header",
				Usage:       "HTTP Header (format: `HeaderName: HeaderValue`)",
				EnvVars:     []string{"GHNOTIFY_REMOTE_HEADER"},
				Destination: &cfg.remoteHeaders,
			},
			&cli.BoolFlag{
				Name:        "remote-idtoken-client",
				Usage:       "Enable IDToken client",
				EnvVars:     []string{"GHNOTIFY_REMOTE_IDTOKEN_CLIENT"},
				Destination: &cfg.remoteIDTokenClient,
			},

			&cli.StringFlag{
				Name:        "slack-webhook",
				Usage:       "Slack incoming webhook",
				EnvVars:     []string{"GHNOTIFY_SLACK_WEBHOOK"},
				Destination: &cfg.SlackWebhookURL,
			},
			&cli.StringFlag{
				Name:        "slack-api-token",
				Usage:       "Slack API token",
				EnvVars:     []string{"GHNOTIFY_SLACK_API_TOKEN"},
				Destination: &cfg.SlackAPIToken,
			},

			&cli.GenericFlag{
				Name:        "log-level",
				Category:    "logging",
				Aliases:     []string{"l"},
				Usage:       "Log level [debug|info|warn|error]",
				EnvVars:     []string{"GHNOTIFY_LOG_LEVEL"},
				Destination: &logLevel,
				Value:       &logLevel,
			},
			&cli.GenericFlag{
				Name:        "log-format",
				Category:    "logging",
				Aliases:     []string{"f"},
				Usage:       "Log format [text|json]",
				EnvVars:     []string{"GHNOTIFY_LOG_FORMAT"},
				Destination: &logFormat,
			},
			&cli.GenericFlag{
				Name:        "log-output",
				Category:    "logging",
				Aliases:     []string{"o"},
				Usage:       "Log output [stdout|stderr]",
				EnvVars:     []string{"GHNOTIFY_LOG_OUTPUT"},
				Destination: &logOutput,
			},
		},
		Commands: []*cli.Command{
			cmdServe(&cfg),
			cmdEmit(&cfg),
		},
		Before: func(ctx *cli.Context) error {
			cfg.RemoteHeaders = cfg.remoteHeaders.Value()

			utils.RenewLogger(logOutput.Writer(), logLevel, logFormat)

			utils.Logger.With("config", cfg).Debug("Starting...")

			return nil
		},
	}

	if err := app.Run(argv); err != nil {
		utils.Logger.Error(err.Error())
		return err
	}
	return nil
}
