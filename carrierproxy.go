package main

import (
	"fmt"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/input"
	"golang.design/x/clipboard"
	"os"
	"time"
	// "github.com/joho/godotenv"
	"io"
	"reflect"
)


type PolicyProvider interface {
	Login(username, password, webpage string) error
	Policies(rod.Elements, rod.Elements) ([]Policy, error)
	DocumentDownload(downloadKey string, page *rod.Page)  (io.ReadCloser, error)
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


func Policies(keys, vals rod.Elements) ([]Policy, error) {
	//create policy objects with keys and vals
	policies := make([]Policy, len(keys))

	for i := 0; i < 3; i++ {
		policies[i] = Policy{
			CarrierID:    keys[i].MustText(),
			PolicyNumber: vals[i].MustText(),
		}
	}
	
	return policies, nil
}



func DocumentDownload(downloadKey string, page *rod.Page) (io.ReadCloser, error){
	//start a sandbox



	// page.MustWaitLoad().MustScreenshot(downloadKey)

	f, err := os.OpenFile(downloadKey, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println(err)
		return
	}

	_, err = fmt.Fprintln(f, newLine)
	if err != nil {
		fmt.Println(err)
		f.Close()
		return
	}
	err = f.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("file appended successfully")

	return os.CloseFile(), nil
}




