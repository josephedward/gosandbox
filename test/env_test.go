package proxy

import(
	"goscraper/local"
	// "goscraper/acloud"
	"goscraper/proxy"
	"testing"
	"fmt"
)


func TestLogin(t *testing.T){
	
	login := local.LoadEnv()
	fmt.Println(login)
	t.Log("Login Test Passed")

	//call the login function
	proxy.Login(login.username, login.password, login.url)

}

