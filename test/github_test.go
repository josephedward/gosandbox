package test

import (
	"context"
	"reflect"
	"testing"

	"github.com/google/go-github/v47/github"
)

func TestGetRepositories(t *testing.T) {
	tests := []struct {
		name string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			GetRepositories()
		})
	}
}

func TestGetToken(t *testing.T) {
	tests := []struct {
		name      string
		wantToken string
		wantErr   bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotToken, err := GetToken()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotToken != tt.wantToken {
				t.Errorf("GetToken() = %v, want %v", gotToken, tt.wantToken)
			}
		})
	}
}

func TestSecretEnv(t *testing.T) {
	tests := []struct {
		name string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SecretEnv()
		})
	}
}

func TestGetSecretName(t *testing.T) {
	type args struct {
		secretName string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetSecretName(tt.args.secretName)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetSecretName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetSecretName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetSecretValue(t *testing.T) {
	type args struct {
		secretValue string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetSecretValue(tt.args.secretValue)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetSecretValue() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetSecretValue() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGithubAuth(t *testing.T) {
	type args struct {
		token string
	}
	tests := []struct {
		name    string
		args    args
		want    context.Context
		want1   *github.Client
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := GithubAuth(tt.args.token)
			if (err != nil) != tt.wantErr {
				t.Errorf("GithubAuth() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GithubAuth() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("GithubAuth() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestAddRepoSecret(t *testing.T) {
	type args struct {
		ctx         context.Context
		client      *github.Client
		owner       string
		repo        string
		secretName  string
		secretValue string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := AddRepoSecret(tt.args.ctx, tt.args.client, tt.args.owner, tt.args.repo, tt.args.secretName, tt.args.secretValue); (err != nil) != tt.wantErr {
				t.Errorf("AddRepoSecret() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_encryptSecretWithPublicKey(t *testing.T) {
	type args struct {
		publicKey   *github.PublicKey
		secretName  string
		secretValue string
	}
	tests := []struct {
		name    string
		args    args
		want    *github.EncryptedSecret
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := encryptSecretWithPublicKey(tt.args.publicKey, tt.args.secretName, tt.args.secretValue)
			if (err != nil) != tt.wantErr {
				t.Errorf("encryptSecretWithPublicKey() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("encryptSecretWithPublicKey() = %v, want %v", got, tt.want)
			}
		})
	}
}
