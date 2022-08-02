package local

import(
	// "github.com/josephedward/glovebox-go-code-challenge/acloud"
	// "github.com/josephedward/glovebox-go-code-challenge/local"
	// "testing"
	// "github.com/go-rod/rod"
	// "github.com/go-rod/rod/lib/input"
	"github.com/joho/godotenv"
	// "golang.design/x/clipboard"
	// "io"
	"log"
	"os"
	// "reflect"
	// "time"
	"fmt"
)


type AcloudLogin struct {
	Url string
	User string
	Password string
}

func LoadEnv() AcloudLogin{
	//load env variables
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatalf("Could not load .env file - Err: %s", err)
	}

	username := os.Getenv("USERNAME")
	password := os.Getenv("PASSWORD")
	url := os.Getenv("URL")
	fmt.Println(username)
	fmt.Println(password)
	fmt.Println(url)
	return AcloudLogin{Url: url, User: username, Password: password}
}