package cmd_test

import (
	"os"
	"testing"
	"time"

	"github.com/m-mizutani/ghnotify/pkg/controller/cmd"
	"github.com/stretchr/testify/assert"
)

func TestRun(t *testing.T) {
	tests := []struct {
		name     string
		envKey   string
		envValue string
		port     string
	}{
		{
			name:     "configure using GHNOTIFY_LOG_LEVEL",
			envKey:   "GHNOTIFY_LOG_LEVEL",
			envValue: "debug",
			port:     "4080",
		},
		{
			name:     "configure using GHNOTIFY_LOG_FORMAT",
			envKey:   "GHNOTIFY_LOG_FORMAT",
			envValue: "json",
			port:     "4081",
		},
		{
			name:     "configure using GHNOTIFY_LOG_OUTPUT",
			envKey:   "GHNOTIFY_LOG_OUTPUT",
			envValue: "stdout",
			port:     "4082",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Setenv(tt.envKey, tt.envValue)
			defer os.Unsetenv(tt.envKey)

			errCh := make(chan error)

			go func() {
				argv := []string{
					"ghnotify",
					"--slack-api-token",
					"dummy",
					"--remote-url",
					"localhost:1234",
					"serve",
					"--addr",
					"0.0.0.0:" + tt.port,
				}
				err := cmd.Run(argv)
				errCh <- err
			}()

			select {
			case err := <-errCh:
				assert.NoError(t, err)
			case <-time.After(time.Second * 1):
				t.Log("cmd.Run exited without error")
			}
		})
	}
}
