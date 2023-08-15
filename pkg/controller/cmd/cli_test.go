package cmd_test

import (
	"os"
	"testing"
	"time"

	"github.com/m-mizutani/ghnotify/pkg/controller/cmd"
	"github.com/stretchr/testify/assert"
)

func TestRun(t *testing.T) {
	t.Run("configure using environmental variable", func(t *testing.T) {
		os.Setenv("GHNOTIFY_LOG_LEVEL", "debug")
		defer os.Unsetenv("GHNOTIFY_LOG_LEVEL")

		errCh := make(chan error)

		go func() {
			argv := []string{
				"ghnotify",
				"--slack-api-token",
				"dummy",
				"--remote-url",
				"localhost:1234",
				"serve",
			}
			err := cmd.Run(argv)
			errCh <- err
		}()

		select {
		case err := <-errCh:
			assert.NoError(t, err)
		case <-time.After(time.Duration(0.1 * float64(time.Second))):
			t.Log("cmd.Run exited without error")
		}
	})
}
