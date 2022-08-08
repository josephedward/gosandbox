package main

import (
	"fmt"
	"goscraper/acloud"
	"goscraper/local"
	"goscraper/proxy"
)

func main() {

	cliEnv, err := local.Execute()
	local.PanicIfErr(err)
	fmt.Println("cliEnv : ", cliEnv)

	//connect to website
	connect, err := proxy.Login(proxy.WebsiteLogin{Url: cliEnv.Url, Username: cliEnv.Username, Password: cliEnv.Password})
	local.PanicIfErr(err)
	fmt.Println("connect : ", connect)

	//scrape credentials
	elems, err := acloud.Sandbox(connect)
	local.PanicIfErr(err)

	//copy credentials to clipboard
	creds, err := acloud.Copy(elems)
	local.PanicIfErr(err)
	fmt.Println("creds : ", creds.User)

	keys, vals := acloud.KeyVals(creds)

	//create policies with map
	policies, err := proxy.Policies(keys, vals)
	local.PanicIfErr(err)
	fmt.Println("policies : ", policies)

	//download text file of policies
	err = proxy.DocumentDownload("creds", policies)
	local.PanicIfErr(err)
	fmt.Println("Document Downloaded")

//create LocalCreds from creds
	//append aws creds to .aws/credentials file
	err = local.AppendAwsCredentials(local.LocalCreds{
		Path:      cliEnv.Aws_path,
		User:      creds.User,
		KeyID:     creds.KeyID,
		AccessKey: creds.AccessKey,
	})
	local.PanicIfErr(err)
	fmt.Println("aws credentials appended")

}
