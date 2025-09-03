package core

import (
	"errors"
	"fmt"
	"time"

	"github.com/go-rod/rod"
)

// for holding information about a login
type WebsiteLogin struct {
	Url      string
	Username string
	Password string
}

// for holding information about a rod connection
type Connection struct {
	Browser *rod.Browser
	Page    *rod.Page
}

func Login(login WebsiteLogin) (Connection, error) {
	// Launch a new browser with default options, and connect to it.
	browser := rod.New()
	err := browser.Connect()
	if err != nil {
		return Connection{}, fmt.Errorf("failed to connect to browser: %w", err)
	}

	// Create a new page
	page := browser.MustPage(login.Url)

	page.WaitIdle(15 * time.Second) // Wait for the page to become idle

	// Print HTML content for debugging
	html, err := page.HTML()
	if err != nil {
		return Connection{}, err
	}
	fmt.Println("--- Page HTML (initial load) ---")
	fmt.Println(html)
	fmt.Println("--- End Page HTML ---")

	//if browser is nil, page is nil
	if browser == nil || page == nil {
		return Connection{}, errors.New("browser or page is nil")
	}

	//create connection object to return
	return Connection{Browser: browser, Page: page}, nil
}

// should be able to work with most websites that use a login form
func SimpleLogin(connect Connection, login WebsiteLogin) (Connection, error) {
	//navigate to login page
	page := connect.Page.MustNavigate(login.Url)

	//Race Condition: It will keep polling until one selector has found a match
	// Pluralsight login
	page.MustElement("#Username").MustInput(login.Username)
	page.MustElement("#Password").MustInput(login.Password)
	page.MustElement("#login-button").MustClick()
	page.MustWaitLoad()

	return Connection{Browser: connect.Browser, Page: page}, nil
}

func Connect(browser *rod.Browser, url string) Connection {
	page := browser.MustPage(url)
	return Connection{Browser: browser, Page: page}
}
