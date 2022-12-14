package test

import (
	"gosandbox/acloud"
	"gosandbox/cli"
	"gosandbox/core"
	"reflect"
	"testing"

	"github.com/manifoldco/promptui"
)

func Test_main(t *testing.T) {
	tests := []struct {
		name string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			main()
		})
	}
}

func TestGetTemplates(t *testing.T) {
	tests := []struct {
		name string
		want *promptui.SelectTemplates
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetTemplates(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetTemplates() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetSearcher(t *testing.T) {
	type args struct {
		options []cli.PromptOptions
	}
	tests := []struct {
		name string
		args args
		want func(input string, index int) bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetSearcher(tt.args.options); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetSearcher() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSelect(t *testing.T) {
	type args struct {
		promptTitle string
		options     []cli.PromptOptions
	}
	tests := []struct {
		name string
		args args
		want *promptui.Select
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Select(tt.args.promptTitle, tt.args.options); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Select() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestExecute(t *testing.T) {
	type args struct {
		p acloud.ACloudProvider
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Execute(tt.args.p)
		})
	}
}

func TestOpenAWSConsole(t *testing.T) {
	type args struct {
		creds acloud.SandboxCredential
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			OpenAWSConsole(tt.args.creds)
		})
	}
}

func TestSandboxToGithub(t *testing.T) {
	type args struct {
		creds acloud.SandboxCredential
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SandboxToGithub(tt.args.creds)
		})
	}
}

func TestDownloadTextFile(t *testing.T) {
	type args struct {
		creds acloud.SandboxCredential
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			DownloadTextFile(tt.args.creds)
		})
	}
}

func TestAppendCreds(t *testing.T) {
	type args struct {
		creds acloud.SandboxCredential
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			AppendCreds(tt.args.creds)
		})
	}
}

func TestGetSandboxCreds(t *testing.T) {
	type args struct {
		cliEnv cli.ACloudEnv
		p      *acloud.ACloudProvider
	}
	tests := []struct {
		name    string
		args    args
		want    acloud.ACloudProvider
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetSandboxCreds(tt.args.cliEnv, tt.args.p)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetSandboxCreds() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetSandboxCreds() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEnvLocation(t *testing.T) {
	tests := []struct {
		name       string
		wantCliEnv cli.ACloudEnv
		wantErr    bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotCliEnv, err := EnvLocation()
			if (err != nil) != tt.wantErr {
				t.Errorf("EnvLocation() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotCliEnv, tt.wantCliEnv) {
				t.Errorf("EnvLocation() = %v, want %v", gotCliEnv, tt.wantCliEnv)
			}
		})
	}
}
