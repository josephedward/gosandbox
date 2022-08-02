package proxy

import(
	"goscraper/local"
	"goscraper/acloud"
	"goscraper/proxy"
	"testing"
	"fmt"
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
	creds , err := acloud.Copy(vals)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	fmt.Println("creds : ",creds)
}

	

