package aws

import (
	"fmt"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/input"
	"gosandbox/core"
)

func Login(login core.WebsiteLogin) (connect core.Connection, err error) {

	// Launch a new browser with default options, and connect to it.
	browser := rod.New().MustConnect()

	// Create a new page
	page := browser.MustPage(login.Url)

	connect = core.Connection{Browser: browser, Page: page}

	//need a race condition as this method is otherwise the same as the previous 
	page.MustElement("input[name='username']").MustInput(login.Username).MustType(input.Enter)
	page.MustElement("input[name='password']").MustInput(login.Password).MustType(input.Enter)

	//need to check for error conditions (might be another race with div.error?) before returning
	return connect, nil
}
