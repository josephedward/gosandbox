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

func LoadEnv() (login WebsiteLogin, err error) {
	//load env variables
	err = godotenv.Load("../.env")
	if err != nil {
		log.Fatalf("Could not load .env file - Err: %s", err)
	}

	username := os.Getenv("USERNAME")
	password := os.Getenv("PASSWORD")
	url := os.Getenv("URL")
	fmt.Println(username)
	fmt.Println(password)
	fmt.Println(url)
	return WebsiteLogin{Url:url, Username:username, Password:password}, err
}