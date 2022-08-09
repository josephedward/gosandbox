package core

import (
	"fmt"
	"os"
	"goscraper/proxy"
)

type LocalCreds struct {
	Path      string
	User      string
	KeyID     string
	AccessKey string
}

func AppendAwsCredentials(creds LocalCreds) error {
	newLine := fmt.Sprintf("\n\n[%s]\n", creds.User)
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

func DocumentDownload(downloadKey string, policies []proxy.Policy) error {
	//create a file with list of policies
	file, err := os.Create(downloadKey + ".txt")
	if err != nil {
		return err
	}
	defer file.Close()

	//write the policies to the file
	for _, policy := range policies {
		_, err := fmt.Fprintln(file, policy.CarrierID, policy.PolicyNumber)
		if err != nil {
			return err
		}
	}

	return nil

}
