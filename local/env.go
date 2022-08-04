package local

import (
	"fmt"
	"goscraper/proxy"
	"os"

	"github.com/joho/godotenv"
)



func LoadEnv() error{
	//load env variables
	err := godotenv.Load("../.env")
	PanicIfErr(err)
	return err
}


func SetEnv()(login proxy.WebsiteLogin, err error) {
	err = LoadEnv()
	PanicIfErr(err)
	username := os.Getenv("USERNAME")
	password := os.Getenv("PASSWORD")
	url := os.Getenv("URL")
	fmt.Println(username)
	fmt.Println(password)
	fmt.Println(url)
	return proxy.WebsiteLogin{Url:url, Username:username, Password:password}, err
}