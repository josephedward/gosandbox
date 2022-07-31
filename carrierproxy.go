package main

import (
	"io"
	
	// "github.com/chromedp/chromedp"
	// "context"
	// "errors"
	"fmt"
	"github.com/joho/godotenv"
	// "math/rand"
	// "net/http"
	// "path/filepath"
	// "sync"
	"time"
	// "github.com/chromedp/chromedp"
	"log"
	"os"
	"github.com/go-rod/rod"
	// "github.com/go-rod/rod/lib/cdp"
	"github.com/go-rod/rod/lib/input"
	// "github.com/go-rod/rod/lib/launcher"
	// "github.com/go-rod/rod/lib/proto"
	// "github.com/go-rod/rod/lib/utils"
	// "github.com/ysmood/gson"
	// "github.com/otiai10/gosseract/v2"
)

type PolicyProvider interface {
	Login(username, password string) error
	Policies() ([]Policy, error)
	DocumentDownload(downloadKey string) (io.ReadCloser, error)
}

func main() {
	// Launch a new browser with default options, and connect to it.
	browser := rod.New().MustConnect()

	// Even you forget to close, rod will close it after main process ends.
	defer browser.MustClose()

	// Create a new page
	page := browser.MustPage("https://learn.acloud.guru/cloud-playground/cloud-sandboxes")

	// fmt.Println(page)
	log.Print(page)

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Some error occured. Err: %s", err)
	}

	//login to the page
	username := os.Getenv("USERNAME")
	password := os.Getenv("PASSWORD")	
	page.MustElement("input[name='email']").MustInput(username).MustType(input.Enter)
	page.MustElement("input[name='password']").MustInput(password).MustType(input.Enter)

	//start a sandbox
	page.MustElementR("button","Start AWS Sandbox").MustClick()
	page.MustWaitLoad().MustScreenshot("creds.png")
		
	//MustElements immediately returns empty list if no element found, so we wait for it. 
	time.Sleep(6 * time.Second)

	// div[attr^="elem"]
	vals :=	page.MustElements("div[class^='CopyableInstanceField__Value']")
	
	for _, v := range vals {
		fmt.Println(v.MustText())
	}
		
}
