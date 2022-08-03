package proxy

import (
	"goscraper/local"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/input"
)

type Connection struct {
	Browser *rod.Browser
	Page *rod.Page
}


func Login(login local.WebsiteLogin) (Connection, error ){
	// Launch a new browser with default options, and connect to it.
	browser:= rod.New().MustConnect()
	
	// Create a new page
	page := browser.MustPage(login.Url)
	// fmt.Println(page)
	
	//login to the page
	page.MustElement("input[name='email']").MustInput(login.Username).MustType(input.Enter)
	page.MustElement("input[name='password']").MustInput(login.Password).MustType(input.Enter)
	
	return Connection{Browser: browser, Page: page}, nil
}
