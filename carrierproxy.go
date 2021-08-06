package main

import (
	"io"
)

type PolicyProvider interface {
	Login(username, password string) error
	Policies() ([]Policy, error)
	DocumentDownload(downloadKey string) (io.ReadCloser, error)
}
