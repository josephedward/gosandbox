package test

import (
	"gosandbox/proxy"
	"gosandbox/core"
	"reflect"
	"testing"
)

func TestLogin(t *testing.T) {
	type args struct {
		login core.WebsiteLogin
	}
	tests := []struct {
		name    string
		args    args
		want    core.Connection
		wantErr bool
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := core.Login(tt.args.login)
			if (err != nil) != tt.wantErr {
				t.Errorf("Login() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Login() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLoadEnv(t *testing.T) {
	tests := []struct {
		name      string
		wantLogin core.ACloudEnv
		wantErr   bool
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotLogin, err := core.LoadEnv()
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadEnv() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotLogin, tt.wantLogin) {
				t.Errorf("LoadEnv() = %v, want %v", gotLogin, tt.wantLogin)
			}
		})
	}
}

func TestLoadEnvPath(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name      string
		args      args
		wantLogin core.ACloudEnv
		wantErr   bool
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotLogin, err := core.LoadEnvPath(tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadEnvPath() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotLogin, tt.wantLogin) {
				t.Errorf("LoadEnvPath() = %v, want %v", gotLogin, tt.wantLogin)
			}
		})
	}
}

func TestAppendAwsCredentials(t *testing.T) {
	type args struct {
		creds core.LocalCreds
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := core.AppendAwsCredentials(tt.args.creds); (err != nil) != tt.wantErr {
				t.Errorf("AppendAwsCredentials() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAppendLine(t *testing.T) {
	type args struct {
		newLine string
		path    string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := core.AppendLine(tt.args.newLine, tt.args.path); (err != nil) != tt.wantErr {
				t.Errorf("AppendLine() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDocumentDownload(t *testing.T) {
	type args struct {
		downloadKey string
		policies    []proxy.Policy
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := core.DocumentDownload(tt.args.downloadKey, tt.args.policies); (err != nil) != tt.wantErr {
				t.Errorf("DocumentDownload() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestScreenShot(t *testing.T) {
	type args struct {
		filename string
		connect  core.Connection
	}
	tests := []struct {
		name string
		args args
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			core.ScreenShot(tt.args.filename, tt.args.connect)
		})
	}
}

// func TestExecute(t *testing.T) {
// 	tests := []struct {
// 		name    string
// 		want    core.ACloudEnv
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
// 		tempEnv core.ACloudEnv
// 	}
// 	tests := []struct {
// 		name    string
// 		args    args
// 		want    core.ACloudEnv
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
// 		tempEnv core.ACloudEnv
// 	}
// 	tests := []struct {
// 		name    string
// 		args    args
// 		want   core.ACloudEnv
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
