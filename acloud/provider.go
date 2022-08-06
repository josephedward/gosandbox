package acloud

import (
	"fmt"

	"goscraper/proxy"
// 	// "github.com/go-rod/rod"
// 	// "golang.design/x/clipboard"
// 	// "time"
	"goscraper/local"
// 	// "github.com/go-rod/rod"
// 	// "github.com/go-rod/rod/lib/input"

)

// type ACloudEnv struct {
// 	un string
// 	pw string
// 	url string
// 	aws_path string
// 	download_key string
// }




type ACloudProvider struct {
	ACloudEnv *local.ACloudEnv
	proxy.Connection
}


func (p *ACloudProvider)Login(username, password string) (err error){

		//load env credentials from .env file
	login, err := local.LoadEnv()
	local.PanicIfErr(err)
	fmt.Println("login : ", login)
	
	
	//connect to website
	connect, err := proxy.Login(local.WebsiteLogin{Url: login.Url, Username: username, Password: password})
	local.PanicIfErr(err)
	fmt.Println("connect : ", connect)
	
	//set the provider's connection
	p.Connection = connect
	// fmt.Println("p.Connection : ", p.Connection)
	// fmt.Println("p.Connection.Session : ", p.Connection.Browser)
	//just log p 
	// fmt.Println("p : ", p)

	return err
}

func (p ACloudProvider)Policies() (policies []proxy.Policy, err error){

		//scrape credentials
		elems, err := Sandbox(p.Connection)
		local.PanicIfErr(err)
	
		//copy credentials to clipboard
		creds, err := Copy(elems)
		local.PanicIfErr(err)
		fmt.Println("creds : ", creds.User)
		
	
		keys, vals := KeyVals(creds)
	
		//create policies with map
		policies, err = proxy.Policies(keys, vals)
		local.PanicIfErr(err)
		fmt.Println("policies : ", policies)
		

	return policies, err
	// //create policies 
	// policies, err = proxy.Policies(p.Connection)
	// local.PanicIfErr(err)
	// fmt.Println("policies : ", policies)
	// return policies, err
}

