package acloud

import (
	"github.com/go-rod/rod"
	"golang.design/x/clipboard"
	"goscraper/proxy"
	"time"
)

type SandboxCredentials struct {
	User      string
	Password  string
	URL       string
	KeyID     string
	AccessKey string
}

func Sandbox(connect proxy.Connection) (rod.Elements, error) {
	//start a sandbox
	connect.Page.MustElementR("button", "Start AWS Sandbox").MustClick()

	//wait for the page to load (I know it is not best practice, but it works)
	time.Sleep(6 * time.Second)
	connect.Page.MustWaitLoad().MustScreenshot("creds.png")

	//find the right elements with traversal pattern div[attr^="elem"]
	elems := connect.Page.MustWaitLoad().MustElements("div[class^='CopyableInstanceField__Value']")

	return elems, nil
}

func Copy(elems rod.Elements) (SandboxCredentials, error) {
	//initialize cliboard package
	err := clipboard.Init()
	if err != nil {
		panic(err)
	}

	//have to copy to clipboard to get whole string
	elems[0].MustElement("svg[aria-label='copy icon']").MustClick()
	// write/read text format data of the clipboard, and
	// the byte buffer regarding the text are UTF8 encoded.
	un := clipboard.Read(clipboard.FmtText)
	//zero out the clipboard just in case
	clipboard.Write(clipboard.FmtText, nil)

	elems[1].MustElement("svg[aria-label='copy icon']").MustClick()
	pw := clipboard.Read(clipboard.FmtText)
	clipboard.Write(clipboard.FmtText, nil)

	elems[2].MustElement("svg[aria-label='copy icon']").MustClick()
	url := clipboard.Read(clipboard.FmtText)
	clipboard.Write(clipboard.FmtText, nil)

	elems[3].MustElement("svg[aria-label='copy icon']").MustClick()
	keyid := clipboard.Read(clipboard.FmtText)
	clipboard.Write(clipboard.FmtText, nil)

	elems[4].MustElement("svg[aria-label='copy icon']").MustClick()
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

func KeyVals(creds SandboxCredentials) ([]string, []string) {
	keys := []string{"username", "password", "url", "keyid", "accesskey"}
	vals := []string{string(creds.User),
		string(creds.Password),
		string(creds.URL),
		string(creds.KeyID),
		string(creds.AccessKey)}

	return keys, vals
}
