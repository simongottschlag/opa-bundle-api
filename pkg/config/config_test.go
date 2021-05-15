package config

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

func TestNewConfig(t *testing.T) {
	envVarsToClear := []string{
		"ADDRESS",
		"PORT",
	}

	for _, envVar := range envVarsToClear {
		restore := tempUnsetEnv(envVar)
		defer restore()
	}

	cliClient := newClient(Options{
		DisableExitOnHelp: true,
	})

	baseArgs := []string{"fake-bin"}
	baseWorkingArgs := append(baseArgs)

	cases := []struct {
		client              *Client
		args                []string
		expectedHosts       []string
		expectedErrContains string
		outBuffer           bytes.Buffer
		errBuffer           bytes.Buffer
	}{
		{
			client:              cliClient,
			args:                baseWorkingArgs,
			expectedErrContains: "",
			outBuffer:           bytes.Buffer{},
			errBuffer:           bytes.Buffer{},
		},
	}

	for _, c := range cases {
		c.client.setIO(&bytes.Buffer{}, &c.outBuffer, &c.errBuffer)
		cfg, err := c.client.generateConfig(c.args)
		if err != nil && c.expectedErrContains == "" {
			t.Errorf("Expected err to be nil: %q", err)
		}

		if err == nil && c.expectedErrContains != "" {
			t.Errorf("Expected err to contain '%s' but was nil", c.expectedErrContains)
		}

		if err != nil && c.expectedErrContains != "" {
			if !strings.Contains(err.Error(), c.expectedErrContains) {
				t.Errorf("Expected err to contain '%s' but was: %q", c.expectedErrContains, err)
			}
		}

		if c.expectedErrContains == "" {
			if cfg.Port != 8080 {
				t.Errorf("Expected cfg.Port to be '8080' but was: %d", cfg.Port)
			}
		}
	}
}

func tempUnsetEnv(key string) func() {
	oldEnv := os.Getenv(key)
	os.Unsetenv(key)
	return func() { os.Setenv(key, oldEnv) }
}
