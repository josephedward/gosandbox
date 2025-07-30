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

func TestScrape(t *testing.T) {
	type args struct {
		connect core.Connection
	}
	tests := []struct {
		name string
		args args
		want rod.Elements
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := acloud.Scrape(tt.args.connect); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Scrape() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCopy(t *testing.T) {
	type args struct {
		elems rod.Elements
	}
	tests := []struct {
		name    string
		args    args
		want    acloud.SandboxCredential
		wantErr bool
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := acloud.Copy(tt.args.elems)
			if (err != nil) != tt.wantErr {
				t.Errorf("Copy() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Copy() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestKeyVals(t *testing.T) {
	type args struct {
		creds acloud.SandboxCredential
	}
	tests := []struct {
		name  string
		args  args
		want  []string
		want1 []string
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := acloud.KeyVals(tt.args.creds)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("KeyVals() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("KeyVals() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestACloudProvider_Login(t *testing.T) {
	type fields struct {
		ACloudEnv          cli.ACloudEnv
		Connection         core.Connection
		SandboxCredential acloud.SandboxCredential
	}
	type args struct {
		username string
		password string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &acloud.ACloudProvider{
				ACloudEnv:          tt.fields.ACloudEnv,
				Connection:         tt.fields.Connection,
				SandboxCredential: tt.fields.SandboxCredential,
			}
			if err := p.Login(tt.args.username, tt.args.password); (err != nil) != tt.wantErr {
				t.Errorf("ACloudProvider.Login() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestACloudProvider_Policies(t *testing.T) {
	type fields struct {
		ACloudEnv          cli.ACloudEnv
		Connection         core.Connection
		SandboxCredential acloud.SandboxCredential
	}
	tests := []struct {
		name         string
		fields       fields
		wantPolicies []proxy.Policy
		wantErr      bool
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &acloud.ACloudProvider{
				ACloudEnv:          tt.fields.ACloudEnv,
				Connection:         tt.fields.Connection,
				SandboxCredential: tt.fields.SandboxCredential,
			}
			gotPolicies, err := p.Policies()
			if (err != nil) != tt.wantErr {
				t.Errorf("ACloudProvider.Policies() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotPolicies, tt.wantPolicies) {
				t.Errorf("ACloudProvider.Policies() = %v, want %v", gotPolicies, tt.wantPolicies)
			}
		})
	}
}

func TestACloudProvider_DocumentDownload(t *testing.T) {
	type fields struct {
		ACloudEnv          cli.ACloudEnv
		Connection         core.Connection
		SandboxCredential acloud.SandboxCredential
	}
	type args struct {
		downloadKey string
		policies    []proxy.Policy
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &acloud.ACloudProvider{
				ACloudEnv:          tt.fields.ACloudEnv,
				Connection:         tt.fields.Connection,
				SandboxCredential: tt.fields.SandboxCredential,
			}
			if err := p.DocumentDownload(tt.args.downloadKey, tt.args.policies); (err != nil) != tt.wantErr {
				t.Errorf("ACloudProvider.DocumentDownload() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
