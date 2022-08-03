package proxy

import (
	"fmt"
)

type Policy struct {
	CarrierID    string
	PolicyNumber string
}


func Policies(keys, vals []string) ([]Policy, error) {
	//if keys and vals aren't the same length, return an error
	if len(keys) != len(vals) {
		return nil, fmt.Errorf("keys and vals must be the same length")
	}

	//create a slice of policies
	policies := make([]Policy, len(keys))

	//loop through the keys and vals and assign the values to the policies
	for i, key := range keys {
		policies[i].CarrierID = key
		policies[i].PolicyNumber = vals[i]
	}

	return policies, nil
}

