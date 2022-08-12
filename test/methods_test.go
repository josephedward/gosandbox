package test

import (
	"fmt"
	"gosandbox/acloud"
	"gosandbox/core"
	"gosandbox/proxy"
	"testing"
)

func TestMethods(t *testing.T) {
	//load env credentials from .env file
	login, err := core.LoadEnv()
	core.PrintIfErr(err)
	fmt.Println("login : ", login)
	t.Log("login : ", login)

	//connect to website
	connect, err := core.Login(core.WebsiteLogin{Url: login.Url, Username: login.Username, Password: login.Password})
	core.PrintIfErr(err)
	fmt.Println("connect : ", connect)
	t.Log("connect : ", connect)

	//scrape credentials
	elems, err := acloud.Sandbox(connect)
	core.PrintIfErr(err)

	//copy credentials to clipboard
	creds, err := acloud.Copy(elems)
	core.PrintIfErr(err)
	fmt.Println("creds : ", creds.User)
	t.Log("creds : ", creds.User)

	keys, vals := acloud.KeyVals(creds)

	//create policies with map
	policies, err := proxy.Policies(keys, vals)
	core.PrintIfErr(err)
	fmt.Println("policies : ", policies)
	t.Log("policies : ", policies)

	//download text file of policies
	err = core.DocumentDownload("creds", policies)
	core.PrintIfErr(err)
	fmt.Println("Document Downloaded")
	t.Log("Document Downloaded")

	//create LocalCreds from creds
	//append aws creds to .aws/credentials file
	err = core.AppendAwsCredentials(core.LocalCreds{
		Path:      login.Aws_path,
		User:      creds.User,
		KeyID:     creds.KeyID,
		AccessKey: creds.AccessKey,
	})
	core.PrintIfErr(err)
	fmt.Println("aws credentials appended")
	t.Log("aws credentials appended")

}
