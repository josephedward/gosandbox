package main

import(
	"goscraper/local"
	"goscraper/acloud"
	"goscraper/proxy"
	"testing"
	"fmt"
	// "encoding/json"
)


func TestLogin(t *testing.T){
	
	//load env credentials from .env file
	login, err := local.LoadEnv()
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	fmt.Println("login : ",login)

	//connect to website
	connect, err := proxy.Login(login)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	fmt.Println("connect : ",connect)

	//scrape credentials
	vals, err := acloud.Sandbox(connect, login.Url)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	fmt.Println("vals : ",vals)

	//copy credentials to clipboard
	creds , err := acloud.Copy(vals)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}	
	fmt.Println("creds : ",creds.User)

	keys := []string{"username","password","url","keyid","accesskey"}
	keyVals := []string{string(creds.User),
		string(creds.Password),
		string(creds.URL),
		string(creds.KeyID),
		string(creds.AccessKey)}


	//create policies with map
	policies, err := proxy.Policies(keys, keyVals)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}	
	fmt.Println("policies : ",policies)
	
	//download text file of policies
	err = proxy.DocumentDownload("creds",policies)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	fmt.Println("Document Downloaded")


	//create LocalCreds from creds
	localCreds, err := local.CreateLocalCreds(creds.User, creds.KeyID, creds.AccessKey)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	fmt.Println("localCreds : ",localCreds)

	// //append aws creds to .aws/credentials file
	// err = local.appendAwsCredentials(localCreds)
	// if err != nil {
	// 	fmt.Println(err)
	// 	panic(err)
	// }

}

	

