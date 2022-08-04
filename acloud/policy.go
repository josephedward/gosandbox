package acloud

import (
	"fmt"
	"goscraper/local"
	"goscraper/proxy"

)

func (p *ACloudProvider) Policies(keys, vals []string) ([]proxy.Policy, error) {
	//if keys or vals are not provided, scrape them from sandbox
	fmt.Println("keys : ", len(keys))
	fmt.Println("vals : ", len(vals))
	if len(keys) <= 1 || len(vals) <= 1 {
		elems, err := Sandbox(p.acgConnect)
		p.sandboxElems = elems
		local.PanicIfErr(err)
		
		//copy credentials to clipboard
		creds, err := Copy(elems)
		p.creds = creds
		local.PanicIfErr(err)
		fmt.Println("creds : ", creds.User)

		keys, vals = KeyVals(creds)
	}
	
	//create policies with map
	policies, err := proxy.Policies(keys, vals)
	p.policyObjs = policies
	local.PanicIfErr(err)
	fmt.Println("policies : ", policies)
	return policies, err
}
