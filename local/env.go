package local

import(
	"github.com/joho/godotenv"
	"log"
	"os"
	"fmt"
)


type WebsiteLogin struct {
	Url string
	Username string
	Password string
}

type ACloudEnv struct {
	Url string
	Username string
	Password string
	Aws_path string
	Download_key string
}


func LoadEnv() (login ACloudEnv, err error) {
	//load env variables
	err = godotenv.Load("../.env")
	if err != nil {
		log.Fatalf("Could not load .env file - Err: %s", err)
	}

	username := os.Getenv("USERNAME")
	password := os.Getenv("PASSWORD")
	url := os.Getenv("URL")
	aws_path := os.Getenv("AWS_RELATIVE_PATH")
	download_key := os.Getenv("DOWNLOAD_KEY")
	fmt.Println(username)
	fmt.Println(password)
	fmt.Println(url)
	return ACloudEnv{
		Url: url,
		Username: username,
		Password: password,
		Aws_path: aws_path,
		Download_key: download_key,
		}, err
}

/*
package local

import (
	// "fmt"
	// "goscraper/proxy"
	
	"os"
	"github.com/joho/godotenv"
)







func LoadEnv() error{
	//load env variables
	err := godotenv.Load("../.env")
	PanicIfErr(err)
	return err
}


func SetACloudEnv()(env ACloudEnv, err error) {
	
	err = LoadEnv()
	PanicIfErr(err)
	env.un = os.Getenv("USERNAME")
	env.pw = os.Getenv("PASSWORD")
	env.url = os.Getenv("URL")
	env.aws_path = os.Getenv("AWS_RELATIVE_PATH")
	env.download_key = os.Getenv("DOWNLOAD_KEY")
	return env, err
}
*/