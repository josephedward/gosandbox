package acloud

import (
	"fmt"
	"goscraper/core"
	"goscraper/proxy"
)

type ACloudProvider struct {
	core.ACloudEnv
	core.Connection
	SandboxCredentials
}

func (p *ACloudProvider) Login(username, password string) (err error) {

	//load env credentials from .env file
	login, err := core.LoadEnv()
	core.PrintIfErr(err)
	fmt.Println("login : ", login)

	//connect to website
	connect, err := core.Login(core.WebsiteLogin{Url: login.Url, Username: username, Password: password})
	core.PrintIfErr(err)
	fmt.Println("connect : ", connect)

	//set the provider's connection
	p.Connection = connect
	return err
}

func (p *ACloudProvider) Policies() (policies []proxy.Policy, err error) {

	//scrape credentials
	elems, err := Sandbox(p.Connection)
	core.PrintIfErr(err)

	//copy credentials to clipboard
	creds, err := Copy(elems)
	core.PrintIfErr(err)
	fmt.Println("creds : ", creds.User)

	//set the provider's credentials
	p.SandboxCredentials = creds

	//create string arrays of credentials
	keys, vals := KeyVals(creds)

	//create policies with map
	policies, err = proxy.Policies(keys, vals)
	core.PrintIfErr(err)
	fmt.Println("policies : ", policies)

	return policies, err
}

func (p *ACloudProvider) DocumentDownload(downloadKey string, policies []proxy.Policy) (err error) {

	//download text file of policies
	err = core.DocumentDownload(downloadKey, policies)
	core.PrintIfErr(err)
	fmt.Println("download text file of policies : ", downloadKey)

	//create LocalCreds from creds
	//append aws creds to .aws/credentials file
	fmt.Println("p  :",p)
	err = core.AppendAwsCredentials(core.LocalCreds{
		Path:      p.ACloudEnv.Aws_path,
		User:      p.SandboxCredentials.User,
		KeyID:     p.SandboxCredentials.KeyID,
		AccessKey: p.SandboxCredentials.AccessKey,
	})
	core.PrintIfErr(err)
	fmt.Println("appended aws creds to .aws/credentials file @ ", p.ACloudEnv.Aws_path)
	return err
}
