package proxy

import (
	// "fmt"
	"github.com/go-rod/rod"
	// "github.com/go-rod/rod/lib/input"
	// "github.com/go-rod/rod/lib/launcher"
	// "os"
	"io"
	// "reflect"
)


type PolicyProvider interface {
	Login(username, password string) error
	Policies(rod.Elements, rod.Elements) ([]Policy, error)
	DocumentDownload(downloadKey string, page *rod.Page)  (io.ReadCloser, error)
}



