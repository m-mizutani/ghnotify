package cmd_test

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/m-mizutani/ghnotify/pkg/controller/cmd"
	"github.com/stretchr/testify/assert"
)

func TestRun(t *testing.T) {
	type EnvTest struct {
		port  string
		value string
	}

	testCases := map[string]map[string]EnvTest{
		"GHNOTIFY_LOG_LEVEL": {
			"debug": EnvTest{"4081", "debug"},
			"info":  EnvTest{"4082", "info"},
			"warn":  EnvTest{"4083", "warn"},
			"error": EnvTest{"4084", "error"},
		},
		"GHNOTIFY_LOG_FORMAT": {
			"text": EnvTest{"4085", "console"},
			"json": EnvTest{"4086", "json"},
		},
		"GHNOTIFY_LOG_OUTPUT": {
			"stdout": EnvTest{"4087", "stdout"},
			"stderr": EnvTest{"4088", "stderr"},
		},
	}

	for envKey, variations := range testCases {
		for variationName, variation := range variations {
			testName := fmt.Sprintf("%s as %s", envKey, variationName)
			t.Run(testName, func(t *testing.T) {
				os.Setenv(envKey, variation.value)
				defer os.Unsetenv(envKey)

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
						"0.0.0.0:" + variation.port,
					}
					err := cmd.Run(argv)
					errCh <- err
				}()

				select {
				case err := <-errCh:
					assert.NoError(t, err)
				case <-time.After(time.Duration(0.01 * float64(time.Second))):
					t.Log("cmd.Run exited without error")
				}
			})
		}
	}
}
