package cli

import (
	// "fmt"
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

func Env() (login ACloudEnv) {
	username := os.Getenv("USERNAME")
	//set env variable
	os.Setenv("USERNAME", username)
	password := os.Getenv("PASSWORD")
	os.Setenv("PASSWORD", password)
	url := os.Getenv("URL")
	os.Setenv("URL", url)
	aws_path := os.Getenv("AWS_RELATIVE_PATH")
	os.Setenv("AWS_RELATIVE_PATH", aws_path)
	download_key := os.Getenv("DOWNLOAD_KEY")
	os.Setenv("DOWNLOAD_KEY", download_key)

	// fmt.Print("env: ", username, password, url, aws_path, download_key)
	return ACloudEnv{
		Url:          url,
		Username:     username,
		Password:     password,
		Aws_path:     aws_path,
		Download_key: download_key,
	}
}

func LoadEnv() (login ACloudEnv, err error) {
	//load env variables
	err = godotenv.Load("./.env")
	env := Env()
	// fmt.Print("env: ", env)
	return env, err
}

func LoadEnvPath(path string) (login ACloudEnv, err error) {
	//load env variables
	err = godotenv.Load(path)
	env := Env()
	return env, err
}

func ArgEnv() (login ACloudEnv, err error) {
	//set all needed vendor credentials from cli arg
	username := os.Args[1]
	password := os.Args[2]
	url := os.Args[3]
	download_key := os.Args[4]
	//may not need this yet
	// aws_path := os.Args[3]

	return ACloudEnv{
		Url:          url,
		Username:     username,
		Password:     password,
		Aws_path:     "",
		Download_key: download_key,
	}, err
}
