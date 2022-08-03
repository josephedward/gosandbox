package proxy

import (
	"github.com/go-rod/rod"
	"io"
)

type PolicyProvider interface {
	Login(username, password string) error
	Policies(rod.Elements, rod.Elements) ([]Policy, error)
	DocumentDownload(downloadKey string, page *rod.Page)  (io.ReadCloser, error)
}



