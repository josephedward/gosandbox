package proxy

import (
	"fmt"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/input"
	"goscraper/local"
	"io"
	"os"
)

//original
type PolicyProvider interface {
	Login(username, password string) error
	Policies() ([]Policy, error)
	DocumentDownload(downloadKey string) (io.ReadCloser, error)
}

//My Implementation Signatures
// type PolicyProvider interface {
// 	Login(login local.WebsiteLogin) (Connection, error)
// 	Policies(keys, vals []string) ([]Policy, error)
// 	DocumentDownload(downloadKey string, policies []Policy)  error
// }

type Policy struct {
	CarrierID    string
	PolicyNumber string
}

func Policies(keys, vals []string) ([]Policy, error) {
	//if keys and vals aren't the same length, return an error
	if len(keys) != len(vals) {
		return nil, fmt.Errorf("keys and vals must be the same length")
	}

	//create a slice of policies
	policies := make([]Policy, len(keys))

	//loop through the keys and vals and assign the values to the policies
	for i, key := range keys {
		policies[i].CarrierID = key
		policies[i].PolicyNumber = vals[i]
	}

	return policies, nil
}

type Connection struct {
	Browser *rod.Browser
	Page    *rod.Page
}

func Login(login local.WebsiteLogin) (Connection, error) {
	// Launch a new browser with default options, and connect to it.
	browser := rod.New().MustConnect()

	// Create a new page
	page := browser.MustPage(login.Url)
	// fmt.Println(page)

	//login to the page
	page.MustElement("input[name='email']").MustInput(login.Username).MustType(input.Enter)
	page.MustElement("input[name='password']").MustInput(login.Password).MustType(input.Enter)

	return Connection{Browser: browser, Page: page}, nil
}

func DocumentDownload(downloadKey string, policies []Policy) error {
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
