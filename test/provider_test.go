package test

import (
	"goscraper/acloud"
	"goscraper/core"
	"testing"
)

func TestProvider(t *testing.T) {
	//create provicer 
	var p acloud.ACloudProvider

	//declare empty error 
	var err error

	//load env credentials from .env file
	p.ACloudEnv, err = core.LoadEnv()
	core.PrintIfErr(err)
	//print p ACloudEnv
	t.Log("p.ACloudEnv : ", p.ACloudEnv)

	//use acloud provider to login
	err = p.Login(p.ACloudEnv.Username, p.ACloudEnv.Password)
	core.PrintIfErr(err)
	t.Log("p.Connection : ", p.Connection)

	//create policies
	policies, err := p.Policies()
	core.PrintIfErr(err)
	t.Log("policies : ", policies)

	// //document download
	err = p.DocumentDownload(p.ACloudEnv.Download_key, policies)
	core.PrintIfErr(err)
	t.Log("Document Downloaded")

}
