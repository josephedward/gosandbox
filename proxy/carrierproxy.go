package proxy

import (
	"github.com/go-rod/rod"
)

//original
// type PolicyProvider interface {
// 	Login(username, password string) error
// 	Policies() ([]Policy, error)
// 	DocumentDownload(downloadKey string) (io.ReadCloser, error)
// }


//My Implementation Signatures
type PolicyProvider interface {
	Login(login WebsiteLogin) (Connection, error)
	Policies(keys, vals []string) ([]Policy, error)
	DocumentDownload(downloadKey string, policies []Policy)  error
}

type WebsiteLogin struct {
	Url string
	Username string
	Password string
}

type Connection struct {
	Browser *rod.Browser
	Page *rod.Page
}

type Policy struct {
	CarrierID    string
	PolicyNumber string
}



