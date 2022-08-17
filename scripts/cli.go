package main

import (
	"fmt"
	"gosandbox/acloud"
	"gosandbox/proxy"
	"gosandbox/core"
	// "github.com/manifoldco/promptui"
	// "strings"
)

func main() {

	cliEnv, err := core.Execute()
	core.PrintIfErr(err)
	fmt.Println("cliEnv : ", cliEnv)

	//connect to website
	connect, err := core.Login(core.WebsiteLogin{Url: cliEnv.Url, Username: cliEnv.Username, Password: cliEnv.Password})
	core.PrintIfErr(err)
	fmt.Println("connect : ", connect)

	//scrape credentials
	elems, err := acloud.Sandbox(connect, cliEnv.Download_key)
	core.PrintIfErr(err)

	//copy credentials to clipboard
	creds, err := acloud.Copy(elems)
	core.PrintIfErr(err)
	fmt.Println("creds : ", creds.User)

	//create string arrays of credentials
	keys, vals := acloud.KeyVals(creds)

	//create policies with map
	policies, err := proxy.Policies(keys, vals)
	core.PrintIfErr(err)
	fmt.Println("policies : ", policies)

	//download text file of policies
	err = core.DocumentDownload("creds", policies)
	core.PrintIfErr(err)
	fmt.Println("Document Downloaded")

	//create LocalCreds from creds
	//append aws creds to .aws/credentials file
	err = core.AppendAwsCredentials(core.LocalCreds{
		Path:      cliEnv.Aws_path,
		User:      creds.User,
		KeyID:     creds.KeyID,
		AccessKey: creds.AccessKey,
	})
	core.PrintIfErr(err)
	fmt.Println("aws credentials appended")

}

// type pepper struct {
// 	Name     string
// 	HeatUnit int
// 	Peppers  int
// }

// func main() {
// 	peppers := []pepper{
// 		{Name: "Bell Pepper", HeatUnit: 0, Peppers: 0},
// 		{Name: "Banana Pepper", HeatUnit: 100, Peppers: 1},
// 		{Name: "Poblano", HeatUnit: 1000, Peppers: 2},
// 		{Name: "Jalapeño", HeatUnit: 3500, Peppers: 3},
// 		{Name: "Aleppo", HeatUnit: 10000, Peppers: 4},
// 		{Name: "Tabasco", HeatUnit: 30000, Peppers: 5},
// 		{Name: "Malagueta", HeatUnit: 50000, Peppers: 6},
// 		{Name: "Habanero", HeatUnit: 100000, Peppers: 7},
// 		{Name: "Red Savina Habanero", HeatUnit: 350000, Peppers: 8},
// 		{Name: "Dragon’s Breath", HeatUnit: 855000, Peppers: 9},
// 	}

// 	templates := &promptui.SelectTemplates{
// 		Label:    "{{ . }}?",
// 		Active:   "\U0001F336 {{ .Name | cyan }} ({{ .HeatUnit | red }})",
// 		Inactive: "  {{ .Name | cyan }} ({{ .HeatUnit | red }})",
// 		Selected: "\U0001F336 {{ .Name | red | cyan }}",
// 		Details: `
// --------- Pepper ----------
// {{ "Name:" | faint }}	{{ .Name }}
// {{ "Heat Unit:" | faint }}	{{ .HeatUnit }}
// {{ "Peppers:" | faint }}	{{ .Peppers }}`,
// 	}

// 	searcher := func(input string, index int) bool {
// 		pepper := peppers[index]
// 		name := strings.Replace(strings.ToLower(pepper.Name), " ", "", -1)
// 		input = strings.Replace(strings.ToLower(input), " ", "", -1)

// 		return strings.Contains(name, input)
// 	}

// 	prompt := promptui.Select{
// 		Label:     "Spicy Level",
// 		Items:     peppers,
// 		Templates: templates,
// 		Size:      4,
// 		Searcher:  searcher,
// 	}

// 	i, _, err := prompt.Run()

// 	if err != nil {
// 		fmt.Printf("Prompt failed %v\n", err)
// 		return
// 	}

// 	fmt.Printf("You choose number %d: %s\n", i+1, peppers[i].Name)
// }
