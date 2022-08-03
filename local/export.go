package local

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
)

type LocalCreds struct {
	Path      string
	User      string
	KeyID     string
	AccessKey string
}

func LoadAwsPath() (string, error) {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatalf("Could not load .env file - Err: %s", err)
	}
	return os.Getenv("AWS_RELATIVE_PATH"), err
}

func CreateLocalCreds(un, keyid, accesskey string) (LocalCreds, error) {
	path, err := LoadAwsPath()
	if err != nil {
		log.Fatalf("Could not load .env file - Err: %s", err)
	}

	return LocalCreds{
		Path:      path,
		User:      un,
		KeyID:     keyid,
		AccessKey: accesskey,
	}, err
}

func AppendAwsCredentials(creds LocalCreds) error {
	newLine := fmt.Sprintf("[%s]\n\n\n", creds.User)
	newLine += fmt.Sprintf("aws_access_key_id = %s\n", creds.KeyID)
	newLine += fmt.Sprintf("aws_secret_access_key = %s\n", creds.AccessKey)
	err := AppendLine(newLine, creds.Path)

	return err
}

func AppendLine(newLine string, path string) error {
	f, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println(err)
		return err
	}

	_, err = fmt.Fprintln(f, newLine)
	if err != nil {
		fmt.Println(err)
		f.Close()
		return err
	}
	err = f.Close()
	if err != nil {
		fmt.Println(err)
		return err
	}

	fmt.Println("file appended successfully")
	return err
}
