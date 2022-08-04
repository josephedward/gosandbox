package test

import (
	"fmt"
	"goscraper/acloud"
	"goscraper/local"
	// "goscraper/proxy"
	"testing"

)


// type ACloudProvider struct {
// 	loginInfo proxy.WebsiteLogin
// 	acgConnect proxy.Connection
// 	sandboxElems rod.Elements
// 	policyObjs []proxy.Policy
// 	downloadKey string 
// 	creds SandboxCredentials
// 	local local.LocalCreds
// }



func TestProvider(t *testing.T) {

	acgProvider := acloud.ACloudProvider{}
	
	//login
	login, err := local.SetEnv()
	local.PanicIfErr(err)
	fmt.Println("login : ", login)
	t.Log("login : ", login)

	//use acloud provider to login
	connect, err := acgProvider.Login(login)
	local.PanicIfErr(err)
	fmt.Println("login : ", connect)
	t.Log("login : ", connect)

	//create policies 
	policies, err := acgProvider.Policies([]string{""}, []string{""})
	local.PanicIfErr(err)
	fmt.Println("policies : ", policies)
	t.Log("policies : ", policies)

	//document download
	err = acgProvider.DocumentDownload("", policies)
	local.PanicIfErr(err)
	fmt.Println("Document Downloaded")
	t.Log("Document Downloaded")

}
