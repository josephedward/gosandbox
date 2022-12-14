package main

import (
	"gosandbox/cli"
	"gosandbox/core"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
)

func main() {
	//declare error
	//load login information into memory

	ACloudEnv, err := cli.LoadEnv()
	cli.Success("environment : ", ACloudEnv)
	cli.PrintIfErr(err)

	u := launcher.MustResolveURL("")
	browser := rod.New().ControlURL(u).MustConnect()

	// cli.Success("Connection before: ", Connection)
	Connection := core.Connect(browser, ACloudEnv.Url)
	cli.Success("Connection after: ", Connection)

	// //login to aws
	Connection, err = core.SimpleLogin(Connection, core.WebsiteLogin{ACloudEnv.Url, ACloudEnv.Username, ACloudEnv.Password})
	//wait for 2fa - this is a hack for now, need to remove
	cli.Success("Connection: ", Connection)
	cli.PrintIfErr(err)
}
