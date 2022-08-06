package test

import (
	"fmt"
	"goscraper/acloud"
	"goscraper/local"
	"testing"
)



func TestProvider(t *testing.T) {

	//load env credentials from .env file
	login, err := local.LoadEnv()
	local.PanicIfErr(err)


	var p acloud.ACloudProvider
	
	//use acloud provider to login
	err = p.Login(login.Username, login.Password)
	local.PanicIfErr(err)
	//print p ACloudEnv
	fmt.Println("p.ACloudEnv : ", p.ACloudEnv)
	//just print p
	fmt.Println("p : ", p)
	fmt.Println("p.Connection : ", p.Connection)

	//create policies 
	policies, err := p.Policies()
	local.PanicIfErr(err)
	fmt.Println("policies : ", policies)
	t.Log("policies : ", policies)

	// //document download
	err = p.DocumentDownload("creds", policies)
	local.PanicIfErr(err)
	fmt.Println("Document Downloaded")
	t.Log("Document Downloaded")

}
