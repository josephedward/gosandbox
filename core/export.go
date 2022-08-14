package core

import (
	"fmt"
	"gosandbox/proxy"
	"os"
)

// LocalCreds is a struct that holds the credentials for a user
type LocalCreds struct {
	Path      string
	User      string
	KeyID     string
	AccessKey string
}

// Adds CLI credentials to the config file
func AppendAwsCredentials(creds LocalCreds) error {
	newLine := fmt.Sprintf("\n\n[%s]\n", creds.User)
	newLine += fmt.Sprintf("aws_access_key_id = %s\n", creds.KeyID)
	newLine += fmt.Sprintf("aws_secret_access_key = %s\n", creds.AccessKey)
	err := AppendLine(newLine, creds.Path)

	return err
}

// Appends a line to a specified text file
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


//Creates a text file with exported 'policies'
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

//screenshots the current browser window of the connection passed to it. 
func ScreenShot(filename string, connect Connection) {
	connect.Page.MustWaitLoad().MustScreenshot(filename + ".png")
}
