package local

import(
	// "github.com/josephedward/glovebox-go-code-challenge/acloud"
	// "github.com/josephedward/glovebox-go-code-challenge/local"
	// "testing"
	// "github.com/go-rod/rod"
	// "github.com/go-rod/rod/lib/input"
	// "github.com/joho/godotenv"
	// "golang.design/x/clipboard"
	// "io"
	// "log"
	"os"
	// "reflect"
	// "time"
	"fmt"
)

type LocalCreds struct {
	Path string
	User string
	KeyID string
	AccessKey string
}



func appendLine(newLine string) {
	f, err := os.OpenFile("../../../.aws/credentials", os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println(err)
		return
	}

	_, err = fmt.Fprintln(f, newLine)
	if err != nil {
		fmt.Println(err)
		f.Close()
		return
	}
	err = f.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("file appended successfully")
}
