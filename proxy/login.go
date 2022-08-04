package proxy

import (
	// "fmt"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/input"
	// "goscraper/local"

)

func Login(login WebsiteLogin) (Connection, error) {

	// Launch a new browser with default options, and connect to it.
	browser := rod.New().MustConnect()

	// Create a new page
	page := browser.MustPage(login.Url)

	//login to the page
	page.MustElement("input[name='email']").MustInput(login.Username).MustType(input.Enter)
	page.MustElement("input[name='password']").MustInput(login.Password).MustType(input.Enter)

	return Connection{Browser: browser, Page: page}, nil
}
