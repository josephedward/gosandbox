package main

import (
	"fmt"
	"gosandbox/acloud"
	"gosandbox/proxy"
)

func main() {

	cliEnv, err := core.Execute()
	core.PrintIfErr(err)
	fmt.Println("cliEnv : ", cliEnv)

	//connect to website
	connect, err := core.Login(core.WebsiteLogin{Url: cliEnv.Url, Username: cliEnv.Username, Password: cliEnv.Password})
	core.PrintIfErr(err)
	fmt.Println("connect : ", connect)

	//scrape credentials
	elems, err := acloud.Sandbox(connect)
	core.PrintIfErr(err)

	//copy credentials to clipboard
	creds, err := acloud.Copy(elems)
	core.PrintIfErr(err)
	fmt.Println("creds : ", creds.User)

	//create string arrays of credentials
	keys, vals := acloud.KeyVals(creds)

	//create policies with map
	policies, err := proxy.Policies(keys, vals)
	core.PrintIfErr(err)
	fmt.Println("policies : ", policies)

	//download text file of policies
	err = core.DocumentDownload("creds", policies)
	core.PrintIfErr(err)
	fmt.Println("Document Downloaded")

	//create LocalCreds from creds
	//append aws creds to .aws/credentials file
	err = core.AppendAwsCredentials(core.LocalCreds{
		Path:      cliEnv.Aws_path,
		User:      creds.User,
		KeyID:     creds.KeyID,
		AccessKey: creds.AccessKey,
	})
	core.PrintIfErr(err)
	fmt.Println("aws credentials appended")

}
