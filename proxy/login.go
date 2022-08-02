package proxy

import (
	"fmt"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/input"
	"github.com/go-rod/rod/lib/launcher"
	"os"
	"io"
	"reflect"
)

type AccessInstance struct {
	browser *rod.Browser
	page *rod.Page
}


func Login(username, password, webpage string) error {
	
	// Launch a new browser with default options, and connect to it.
	browser := rod.New().MustConnect()
	
	// Even you forget to close, rod will close it after main process ends.
	defer browser.MustClose()
	
	// Create a new page
	page := browser.MustPage(webpage)
	fmt.Println(reflect.TypeOf(webpage))
	fmt.Println(page)

	page.MustElement("input[name='email']").MustInput(username).MustType(input.Enter)
	page.MustElement("input[name='password']").MustInput(password).MustType(input.Enter)
	page.MustElementR("button", "Start AWS Sandbox").MustClick()

	// div[attr^="elem"]
	keys := page.MustElements("div[class^='CopyableInstanceField__Value'] label")
	vals := page.MustElements("div[class^='CopyableInstanceField__Value']")	
	
	Policies(keys,vals)
	DocumentDownload("download", page)

	return nil 
}