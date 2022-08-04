package acloud

import (
	"fmt"
	"goscraper/local"
	"goscraper/proxy"
	"os"
)

func (p ACloudProvider) Login(login proxy.WebsiteLogin) (proxy.Connection, error) {
	//if login is not provided, use env variables
	if login.Url == "" {
		local.LoadEnv()
		login.Url = os.Getenv("URL")}
	if login.Username == "" {
		local.LoadEnv()
		login.Username = os.Getenv("USERNAME")}
	if login.Password == "" {
		local.LoadEnv()
		login.Password = os.Getenv("PASSWORD")}
	p.loginInfo = login
	fmt.Println("login : ", login)
			 
	//connect to website
	connect, err := proxy.Login(login)
	p.acgConnect = connect
	local.PanicIfErr(err)
	fmt.Println("connect : ", connect)
	return connect,err	
}
