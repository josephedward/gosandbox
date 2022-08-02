package proxy

import (
	"fmt"
	"github.com/go-rod/rod"
	// "github.com/go-rod/rod/lib/input"
	// "github.com/go-rod/rod/lib/launcher"
	"os"
	"io"
	// "reflect"
)


type PolicyProvider interface {
	Login(username, password string) error
	Policies(rod.Elements, rod.Elements) ([]Policy, error)
	DocumentDownload(downloadKey string, page *rod.Page)  (io.ReadCloser, error)
}





func DocumentDownload(downloadKey string, page *rod.Page) (io.ReadCloser, error){
	//start a sandbox


	// page.MustWaitLoad().MustScreenshot(downloadKey)

	f, err := os.OpenFile(downloadKey, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println(err)
		return nil , err
	}

	_, err = fmt.Fprintln(f, downloadKey)
	if err != nil {
		fmt.Println(err)
		f.Close()
		return nil , err
	}
	err = f.Close()
	if err != nil {
		fmt.Println(err)
		return nil , err
	}
	fmt.Println("file appended successfully")

	return nil, nil
}




