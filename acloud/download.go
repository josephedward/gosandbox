package acloud

import (
	"fmt"
	"os"
	"goscraper/local"
	"goscraper/proxy"
	// "github.com/joho/godotenv"
)

func (p *ACloudProvider) DocumentDownload(downloadKey string, policies []proxy.Policy) error {
		//if downloadKey is not provided, use env variables
		if downloadKey == "" {
			local.LoadEnv()
			downloadKey = os.Getenv("DOWNLOAD_KEY")
		}
		p.downloadKey = downloadKey

		//if policies are not provided, use policies from provider
		if len(policies) == 0 {
			policies = p.policyObjs
		}
		//download text file of policies
		err := proxy.DocumentDownload(downloadKey, policies)
		local.PanicIfErr(err)
		fmt.Println("Document Downloaded")
	
		//create LocalCreds from creds
		localCreds, err := local.CreateLocalCreds(p.creds.User, p.creds.KeyID, p.creds.AccessKey)
		p.local = localCreds
		local.PanicIfErr(err)
		fmt.Println("localCreds : ", localCreds)
	
		//append aws creds to .aws/credentials file
		err = local.AppendAwsCredentials(localCreds)
		local.PanicIfErr(err)
		fmt.Println("aws credentials appended")
	
		return err
}
