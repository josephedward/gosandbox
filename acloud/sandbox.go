package acloud

import (
	"fmt"
	"goscraper/proxy"
	"github.com/go-rod/rod"
	"golang.design/x/clipboard"
	"time"
)

type SandboxCredentials struct {
	User      string
	Password  string
	URL       string
	KeyID     string
	AccessKey string
}

func Sandbox(connect proxy.Connection, url string) (rod.Elements, error) {
	//start a sandbox
	connect.Page.MustElementR("button", "Start AWS Sandbox").MustClick()

	//wait for the page to load (I know it is not best practice, but it works)
	time.Sleep(6 * time.Second)
	connect.Page.MustWaitLoad().MustScreenshot("creds.png")

	//find the right elements with traversal pattern div[attr^="elem"]
	vals := connect.Page.MustWaitLoad().MustElements("div[class^='CopyableInstanceField__Value']")

	return vals, nil
}

func Copy(vals rod.Elements) (SandboxCredentials, error) {
	//initialize cliboard package
	err := clipboard.Init()
	if err != nil {
		panic(err)
	}

	//have to copy to clipboard to get whole string
	vals[0].MustElement("svg[aria-label='copy icon']").MustClick()
	// write/read text format data of the clipboard, and
	// the byte buffer regarding the text are UTF8 encoded.
	un := clipboard.Read(clipboard.FmtText)
	//zero out the clipboard just in case
	clipboard.Write(clipboard.FmtText, nil)
	fmt.Println(string(un))

	vals[1].MustElement("svg[aria-label='copy icon']").MustClick()
	pw := clipboard.Read(clipboard.FmtText)
	clipboard.Write(clipboard.FmtText, nil)

	vals[2].MustElement("svg[aria-label='copy icon']").MustClick()
	url := clipboard.Read(clipboard.FmtText)
	clipboard.Write(clipboard.FmtText, nil)

	vals[3].MustElement("svg[aria-label='copy icon']").MustClick()
	keyid := clipboard.Read(clipboard.FmtText)
	clipboard.Write(clipboard.FmtText, nil)

	vals[4].MustElement("svg[aria-label='copy icon']").MustClick()
	accesskey := clipboard.Read(clipboard.FmtText)
	clipboard.Write(clipboard.FmtText, nil)

	return SandboxCredentials{
		User:      string(un),
		Password:  string(pw),
		URL:       string(url),
		KeyID:     string(keyid),
		AccessKey: string(accesskey),
	}, nil
}
