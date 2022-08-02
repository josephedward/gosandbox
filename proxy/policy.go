package proxy

import (
	"fmt"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/input"
	"github.com/go-rod/rod/lib/launcher"
	"os"
	"io"
	"reflect"
)

type Policy struct {
	CarrierID    string
	PolicyNumber string
}


func Policies(keys, vals rod.Elements) ([]Policy, error) {
	//create policy objects with keys and vals
	policies := make([]Policy, len(keys))

	for i := 0; i < 3; i++ {
		policies[i] = Policy{
			CarrierID:    keys[i].MustText(),
			PolicyNumber: vals[i].MustText(),
		}
	}
	
	return policies, nil
}
