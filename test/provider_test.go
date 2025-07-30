package test

import (
	"gosandbox/acloud"
	"gosandbox/cli"
	"gosandbox/core"
	"testing"
	// "gosandbox/aws"
)

func TestProvider(t *testing.T) {
	//create provider
	var p acloud.ACloudProvider

	//load env credentials from .env file
	var err error
	p.ACloudEnv, err = cli.LoadEnv()
	if err != nil {
		t.Fatalf("Failed to load environment variables: %v", err)
	}
	t.Log("p.ACloudEnv : ", p.ACloudEnv)

	//use acloud provider to login
	err = p.Login(p.ACloudEnv.Username, p.ACloudEnv.Password)
	if err != nil {
		t.Fatalf("Failed to login: %v", err)
	}
	defer p.Connection.Browser.MustClose()
	t.Log("p.Connection : ", p.Connection)

	//create policies
	policies, err := p.Policies()
	if err != nil {
		t.Fatalf("Failed to get policies: %v", err)
	}
	if len(policies) == 0 {
		t.Fatal("Policies are empty")
	}
	t.Log("policies : ", policies)

	// //document download
	err = p.DocumentDownload(p.ACloudEnv.Download_key, policies)
	if err != nil {
		t.Fatalf("Failed to download document: %v", err)
	}
	t.Log("Document Downloaded")

	//login to AWS (for final verification of credentials)))
	if p.SandboxCredential.URL == "" {
		t.Log("Skipping AWS login verification, no URL found in credentials.")
		return
	}
	awsConnect, err := core.Login(core.WebsiteLogin{
		Url:      p.SandboxCredential.URL,
		Username: p.SandboxCredential.User,
		Password: p.SandboxCredential.Password,
	})
	if err != nil {
		t.Fatalf("Failed to login to AWS: %v", err)
	}
	defer awsConnect.Browser.MustClose()
	t.Log("awsConnect : ", awsConnect)

}
