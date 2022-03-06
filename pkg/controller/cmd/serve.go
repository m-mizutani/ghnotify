package cmd

import (
	"github.com/m-mizutani/ghnotify/pkg/controller/server"
	"github.com/urfave/cli/v2"
)

func cmdServe(cfg *globalConfig) *cli.Command {
	var config struct {
		Addr string
	}
	return &cli.Command{
		Name:  "serve",
		Usage: "Run http server",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "addr",
				Usage:       "HTTP server address",
				EnvVars:     []string{"GHNOTIFY_ADDR"},
				Value:       "0.0.0.0:4080",
				Destination: &config.Addr,
			},
			&cli.StringFlag{
				Name:        "webhook-secret",
				Usage:       "GitHub Webhook secret",
				EnvVars:     []string{"GHNOTIFY_WEBHOOK_SECRET"},
				Destination: &cfg.GitHubWebhookSecret,
			},
		},
		Action: func(ctx *cli.Context) error {
			uc, err := cfg.newUsecase()
			if err != nil {
				return err
			}

			srv := server.New(uc)
			if err := srv.Listen(config.Addr); err != nil {
				return err
			}
			return nil
		},
	}
}
