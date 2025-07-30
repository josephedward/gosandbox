package test

import (
	"gosandbox/cli"
	"gosandbox/core"
	"gosandbox/proxy"
	"os"
	"reflect"
	"testing"
)

func TestLogin_Integration(t *testing.T) {
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

	t.Run("successful login", func(t *testing.T) {
		connection, err := core.Login(login)
		if err != nil {
			t.Fatalf("Login() error = %v, wantErr %v", err, false)
		}
		defer connection.Browser.MustClose()

		title, err := connection.Page.Title()
		if err != nil {
			t.Fatalf("failed to get page title: %v", err)
		}

		// A bit brittle, but confirms we are past the login page.
		expectedTitle := "A Cloud Guru"
		if title != expectedTitle {
			t.Errorf("Login() wrong page title = %q, want %q", title, expectedTitle)
		}
	})
}


// func TestExecute(t *testing.T) {
// 	tests := []struct {
// 		name    string
// 		want    cli.ACloudEnv
// 		wantErr bool
// 	}{
// 		// TODO: Add test cases.
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			got, err := core.Execute()
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("Execute() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}
// 			if !reflect.DeepEqual(got, tt.want) {
// 				t.Errorf("Execute() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }

// func Test_promptEnvFile(t *testing.T) {
// 	type args struct {
// 		tempEnv cli.ACloudEnv
// 	}
// 	tests := []struct {
// 		name    string
// 		args    args
// 		want    cli.ACloudEnv
// 		wantErr bool
// 	}{
// 		// TODO: Add test cases.
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			got, err := core.promptEnvFile(tt.args.tempEnv)
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("promptEnvFile() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}
// 			if !reflect.DeepEqual(got, tt.want) {
// 				t.Errorf("promptEnvFile() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }

// func Test_promptManual(t *testing.T) {
// 	type args struct {
// 		tempEnv cli.ACloudEnv
// 	}
// 	tests := []struct {
// 		name    string
// 		args    args
// 		want   cli.ACloudEnv
// 		wantErr bool
// 	}{
// 		// TODO: Add test cases.
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			got, err := core.promptManual(tt.args.tempEnv)
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("promptManual() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}
// 			if !reflect.DeepEqual(got, tt.want) {
// 				t.Errorf("promptManual() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }

// func Test_promptGetInput(t *testing.T) {
// 	type args struct {
// 		pc core.PromptContent
// 	}
// 	tests := []struct {
// 		name string
// 		args args
// 		want string
// 	}{
// 		// TODO: Add test cases.
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			if got := core.promptGetInput(tt.args.pc); got != tt.want {
// 				t.Errorf("promptGetInput() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }

// func TestPrintIfErr(t *testing.T) {
// 	type args struct {
// 		err error
// 	}
// 	tests := []struct {
// 		name string
// 		args args
// 	}{
// 		// TODO: Add test cases.
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			cli.PrintIfErr(tt.args.err)
// 		})
// 	}
// }
