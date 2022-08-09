package core

import (
	"github.com/joho/godotenv"
	"os"
)

type ACloudEnv struct {
	Url          string
	Username     string
	Password     string
	Aws_path     string
	Download_key string
}

func LoadEnv() (login ACloudEnv, err error) {

	//load env variables
	err = godotenv.Load("../.env")
	
	//set all needed vendor credentials
	username := os.Getenv("USERNAME")
	password := os.Getenv("PASSWORD")
	url := os.Getenv("URL")
	aws_path := os.Getenv("AWS_RELATIVE_PATH")
	download_key := os.Getenv("DOWNLOAD_KEY")
	return ACloudEnv{
		Url:          url,
		Username:     username,
		Password:     password,
		Aws_path:     aws_path,
		Download_key: download_key,
	}, err
}

func LoadEnvPath(path string) (login ACloudEnv, err error) {

	//load env variables
	err = godotenv.Load(path)

	//set all needed vendor credentials
	username := os.Getenv("USERNAME")
	password := os.Getenv("PASSWORD")
	url := os.Getenv("URL")
	aws_path := os.Getenv("AWS_RELATIVE_PATH")
	download_key := os.Getenv("DOWNLOAD_KEY")
	return ACloudEnv{
		Url:          url,
		Username:     username,
		Password:     password,
		Aws_path:     aws_path,
		Download_key: download_key,
	}, err
}
