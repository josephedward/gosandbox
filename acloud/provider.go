package acloud

import (
	"fmt"
	"goscraper/local"
	"goscraper/proxy"
)

type ACloudProvider struct {
	ACloudEnv *local.ACloudEnv
	proxy.Connection
	SandboxCredentials
}

func (p *ACloudProvider) Login(username, password string) (err error) {

	//load env credentials from .env file
	login, err := local.LoadEnv()
	local.PanicIfErr(err)
	fmt.Println("login : ", login)

	//connect to website
	connect, err := proxy.Login(proxy.WebsiteLogin{Url: login.Url, Username: username, Password: password})
	local.PanicIfErr(err)
	fmt.Println("connect : ", connect)

	//set the provider's connection
	p.Connection = connect
	return err
}

func (p *ACloudProvider) Policies() (policies []proxy.Policy, err error) {

	//scrape credentials
	elems, err := Sandbox(p.Connection)
	local.PanicIfErr(err)

	//copy credentials to clipboard
	creds, err := Copy(elems)
	local.PanicIfErr(err)
	fmt.Println("creds : ", creds.User)

	//set the provider's credentials
	p.SandboxCredentials = creds

	//create string arrays of credentials
	keys, vals := KeyVals(creds)

	//create policies with map
	policies, err = proxy.Policies(keys, vals)
	local.PanicIfErr(err)
	fmt.Println("policies : ", policies)

	return policies, err
}

func (p ACloudProvider) DocumentDownload(downloadKey string, policies []proxy.Policy) (err error) {

	//create LocalCreds from creds
	localCreds, err := local.CreateLocalCreds(p.SandboxCredentials.User, p.SandboxCredentials.KeyID, p.SandboxCredentials.AccessKey)
	local.PanicIfErr(err)
	fmt.Println("localCreds : ", localCreds)

	//append aws creds to .aws/credentials file
	err = local.AppendAwsCredentials(localCreds)
	local.PanicIfErr(err)
	fmt.Println("aws credentials appended")

	return err
}
