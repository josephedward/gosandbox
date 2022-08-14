package main

import (
	"fmt"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/input"
	"github.com/joho/godotenv"
	"golang.design/x/clipboard"
	"log"
	"os"
	"time"
)

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
		log.Fatalf("Could not load .env file - Err: %s", err)
	}

	//login to the page
	username := os.Getenv("USERNAME")
	password := os.Getenv("PASSWORD")
	page.MustElement("input[name='email']").MustInput(username).MustType(input.Enter)
	page.MustElement("input[name='password']").MustInput(password).MustType(input.Enter)

	//start a sandbox
	page.MustElementR("button", "Start AWS Sandbox").MustClick()
	page.MustWaitLoad().MustScreenshot("creds.png")

	//MustElements immediately returns empty list if no element found, so we wait for it.
	time.Sleep(6 * time.Second)

	// div[attr^="elem"]
	vals := page.MustElements("div[class^='CopyableInstanceField__Value']")

	//have to copy to clipboard to get whole string
	vals[3].MustElement("svg[aria-label='copy icon']").MustClick()

	//initialize cliboard package
	err = clipboard.Init()
	if err != nil {
		panic(err)
	}

	//new entry for aws sandbox creds
	appendLine("\n\n")
	appendLine("[cloud_user]")

	// write/read text format data of the clipboard, and
	// the byte buffer regarding the text are UTF8 encoded.
	keyid := clipboard.Read(clipboard.FmtText)

	//zero out the clipboard just in case
	clipboard.Write(clipboard.FmtText, nil)

	fmt.Println(keyid)
	//append key id to the file
	appendLine("aws_access_key_id = " + string(keyid))

	vals[4].MustElement("svg[aria-label='copy icon']").MustClick()
	accessKey := clipboard.Read(clipboard.FmtText)
	clipboard.Write(clipboard.FmtText, nil)
	fmt.Println(accessKey)

	//append access key to the file
	appendLine("aws_secret_access_key = " + string(accessKey))

}

func appendLine(newLine string) {
	f, err := os.OpenFile("../../../.aws/credentials", os.O_APPEND|os.O_WRONLY, 0644)
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
}
