package test

import (
	"os"
	"reflect"
	"testing"

	"github.com/go-rod/rod"
	"gosandbox/acloud"
	"gosandbox/cli"
	"gosandbox/core"
	"gosandbox/proxy"
)

func TestSandbox_Integration(t *testing.T) {
	if os.Getenv("INTEGRATION_TESTS") == "" {
		t.Skip("skipping integration test")
	}

	loginEnv, err := cli.LoadEnv()
	if err != nil {
		t.Fatalf("failed to load env: %v", err)
	}

	login := core.WebsiteLogin{
		Url:      loginEnv.Url,
		Username: loginEnv.Username,
		Password: loginEnv.Password,
	}

	connection, err := core.Login(login)
	if err != nil {
		t.Fatalf("Login() error = %v", err)
	}
	defer connection.Browser.MustClose()

	t.Run("scrape sandbox credentials", func(t *testing.T) {
		elems, err := acloud.Sandbox(connection, "test-sandbox-screenshot")
		if err != nil {
			t.Fatalf("Sandbox() error = %v, wantErr %v", err, false)
		}

		if len(elems) == 0 {
			t.Error("Sandbox() returned no elements, expected credentials")
		}
	})
}

