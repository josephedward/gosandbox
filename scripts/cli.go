package main

import (
	"gosandbox/acloud"
	"gosandbox/core"
	"gosandbox/proxy"
)

func main() {

	cliEnv, err := core.Execute()
	core.PrintIfErr(err)
	core.Success("environment : ", cliEnv)

	//connect to website
	connect, err := core.Login(core.WebsiteLogin{Url: cliEnv.Url, Username: cliEnv.Username, Password: cliEnv.Password})
	core.PrintIfErr(err)
	// fmt.Println("connect : ", connect)
	core.Success("connection : ", connect)

	//scrape credentials
	elems, err := acloud.Sandbox(connect, cliEnv.Download_key)
	core.PrintIfErr(err)
	core.Success("rod html elements : ", elems)

	//copy credentials to clipboard
	creds, err := acloud.Copy(elems)
	core.PrintIfErr(err)
	core.Success("credentials : ", creds)

	//DISPLAY WITH COLORS PROMINENTLY TO THE USER
	acloud.DisplayCreds(creds)

	//create string arrays of credentials
	keys, vals := acloud.KeyVals(creds)
	//create policies with map
	policies, err := proxy.Policies(keys, vals)
	core.PrintIfErr(err)
	core.Success("policies : ", policies)

	//ask if they want to download a text file with the credentials
	if core.PromptDownload() == true {
		// download text file of policies
		filename := core.PromptFileName()
		err = core.DocumentDownload(filename, policies)
		core.PrintIfErr(err)

	}

	//ask if they want the credentials to be added to their aws config
	if core.PromptConfig() == true {
		//ask for path to aws config
		path := core.PromptFilePath()
		// append aws creds to .aws/credentials file
		err = core.AppendAwsCredentials(core.LocalCreds{
			Path:      path,
			User:      creds.User,
			KeyID:     creds.KeyID,
			AccessKey: creds.AccessKey,
		})
		core.PrintIfErr(err)
		core.Success("aws credentials appended @ :", cliEnv.Aws_path)
	}

}
