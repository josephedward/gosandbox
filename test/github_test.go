package test

import (
	// "gosandbox/acloud"
	// "gosandbox/core"
	"gosandbox/gh"
	// "gosandbox/aws"
	"testing"
)

func TestGithub(t *testing.T) {
	//	go run main.go -owner google -repo go-github SECRET_VARIABLE
	gh.SecretEnv()

}
