package main

import (
	"fmt"
	"gosandbox/acloud"
	"gosandbox/proxy"
	"gosandbox/core"
	// "github.com/manifoldco/promptui"
	// "strings"
)


var Reset  = "\033[0m"
var Red    = "\033[31m"
var Green  = "\033[32m"
var Yellow = "\033[33m"
var Blue   = "\033[34m"
var Purple = "\033[35m"
var Cyan   = "\033[36m"
var Gray   = "\033[37m"
var White  = "\033[97m"

func main() {

	cliEnv, err := core.Execute()
	core.PrintIfErr(err)
	fmt.Println(Green+"cliEnv : "+Reset,cliEnv)

	//connect to website
	connect, err := core.Login(core.WebsiteLogin{Url: cliEnv.Url, Username: cliEnv.Username, Password: cliEnv.Password})
	core.PrintIfErr(err)
	fmt.Println("connect : ", connect)

	//scrape credentials
	elems, err := acloud.Sandbox(connect, cliEnv.Download_key)
	core.PrintIfErr(err)

	//copy credentials to clipboard
	creds, err := acloud.Copy(elems)
	core.PrintIfErr(err)
	fmt.Println("credentials for : ", creds.User)

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
