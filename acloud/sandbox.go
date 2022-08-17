package acloud

import (
	"gosandbox/core"
	"errors"
	"github.com/go-rod/rod"
	"golang.design/x/clipboard"
	"time"
	"fmt"
)

type SandboxCredentials struct {
	User      string
	Password  string
	URL       string
	KeyID     string
	AccessKey string
}

func Sandbox(connect core.Connection, downloadKey string) (rod.Elements, error) {

	elems := make(rod.Elements, 0)
	time.Sleep(3 * time.Second)
	// It will keep polling until one selector has found a match
	connect.Page.Race().ElementR("button", "Start AWS Sandbox").MustHandle(func(e *rod.Element) {
		e.MustClick()
		time.Sleep(3 * time.Second)
		core.ScreenShot(downloadKey, connect)
		elems = Scrape(connect)
	}).Element("div[class^='CopyableInstanceField__Value']").MustHandle(func(e *rod.Element) {
		time.Sleep(3 * time.Second)
		core.ScreenShot(downloadKey, connect)
		elems = Scrape(connect)
	}).MustDo()

	if len(elems) == 0 {
		return nil, errors.New("no elements found")
	}
	return elems, nil
}

func Scrape(connect core.Connection) rod.Elements {
	elems := connect.Page.MustWaitLoad().MustElements("div[class^='CopyableInstanceField__Value']")
	return elems
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


func DisplayCreds(creds SandboxCredentials) {
	fmt.Println("--------------------------------------------------------------------------------")
	fmt.Println("Sandbox Credentials:")
	fmt.Println("--------------------------------------------------------------------------------")
	fmt.Println("          "+core.Cyan+"Username: " +core.Yellow+ creds.User+core.Reset)
	fmt.Println("          "+core.Cyan+"Password: " +core.Yellow+ creds.Password+core.Reset)
	fmt.Println("          "+core.Cyan+"URL: " +core.Yellow+ creds.URL+core.Reset)
	fmt.Println("          "+core.Cyan+"KeyID: " +core.Yellow+ creds.KeyID+core.Reset)
	fmt.Println("          "+core.Cyan+"AccessKey: " +core.Yellow+ creds.AccessKey+core.Reset)
	fmt.Println("--------------------------------------------------------------------------------")
}