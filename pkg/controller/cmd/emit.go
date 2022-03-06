package cmd

import (
	"io"
	"os"
	"path/filepath"

	"github.com/m-mizutani/ghnotify/pkg/domain/types"
	"github.com/m-mizutani/goerr"
	"github.com/urfave/cli/v2"
)

func cmdEmit(cfg *globalConfig) *cli.Command {
	var (
		eventFile string
		eventType string
	)
	return &cli.Command{
		Name:  "emit",
		Usage: "Read a local file and handle as event data",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "event-file",
				Usage:       "Event data JSON file, `-` means stdin",
				Aliases:     []string{"f"},
				EnvVars:     []string{"GHNOTIFY_EVENT_FILE"},
				Destination: &eventFile,
				Required:    true,
			},
			&cli.StringFlag{
				Name:        "event-type",
				Usage:       "GitHub event type",
				Aliases:     []string{"t"},
				EnvVars:     []string{"GHNOTIFY_EVENT_TYPE"},
				Destination: &eventType,
				Required:    true,
			},
		},
		Action: func(c *cli.Context) error {
			uc, err := cfg.newUsecase()
			if err != nil {
				return err
			}

			var data []byte
			if eventFile != "-" {
				raw, err := os.ReadFile(filepath.Clean(eventFile))
				if err != nil {
					return goerr.Wrap(err)
				}
				data = raw
			} else {
				raw, err := io.ReadAll(os.Stdin)
				if err != nil {
					return goerr.Wrap(err)
				}
				data = raw
			}

			ctx := types.NewContext(types.WithCtx(c.Context))
			if err := uc.HandleGitHubEvent(ctx, eventType, data); err != nil {
				return err
			}
			return nil
		},
	}
}
