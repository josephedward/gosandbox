package proxy

import (
	"fmt"
	"os"
)

func DocumentDownload(downloadKey string, policies []Policy) error {
	//create a file with list of policies
	file, err := os.Create(downloadKey + ".txt")
	if err != nil {
		return err
	}
	defer file.Close()

	//write the policies to the file
	for _, policy := range policies {
		_, err := fmt.Fprintln(file, policy.CarrierID, policy.PolicyNumber)
		if err != nil {
			return err
		}
	}

	return nil

}
