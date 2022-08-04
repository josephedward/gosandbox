package acloud

import (
	// "fmt"
	"goscraper/local"
	"goscraper/proxy"
	"github.com/go-rod/rod"
	// "github.com/go-rod/rod"
	// "golang.design/x/clipboard"
	// "time"
)

type ACloudProvider struct {
	loginInfo proxy.WebsiteLogin
	acgConnect proxy.Connection
	sandboxElems rod.Elements
	policyObjs []proxy.Policy
	downloadKey string 
	creds SandboxCredentials
	local local.LocalCreds
}

