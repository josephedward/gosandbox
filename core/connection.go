package core

import (
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/input"
)

type WebsiteLogin struct {
	Url      string
	Username string
	Password string
}

type Connection struct {
	Browser *rod.Browser
	Page    *rod.Page
}

func Login(login WebsiteLogin) (Connection, error) {
	// Launch a new browser with default options, and connect to it.
	browser := rod.New().MustConnect()

	// Create a new page
	page := browser.MustPage(login.Url)

	//login to the page
	page.MustElement("input[name='email']").MustInput(login.Username).MustType(input.Enter)
	page.MustElement("input[name='password']").MustInput(login.Password).MustType(input.Enter)

	//create connection object to return
	return Connection{Browser: browser, Page: page}, nil
}