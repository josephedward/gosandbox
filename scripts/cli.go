package main

import (
	"gosandbox/acloud"
	"gosandbox/core"
	"gosandbox/proxy"
	"github.com/manifoldco/promptui"
	"os"
	"strings"
	"errors"
	"fmt"
)

func main() {

	//create provider
	// var p acloud.ACloudProvider
	
	Execute()
	// core.Success("environment : ", p.ACloudEnv)
	// cliEnv := p.ACloudEnv

	// //connect to website
	// // connect, err := core.Login(core.WebsiteLogin{Url: cliEnv.Url, Username: cliEnv.Username, Password: cliEnv.Password})
	// // core.PrintIfErr(err)
	// // // fmt.Println("connect : ", connect)
	// // core.Success("connection : ", connect)

	// //scrape credentials
	// elems, err := acloud.Sandbox(connect, cliEnv.Download_key)
	// core.PrintIfErr(err)
	// core.Success("rod html elements : ", elems)

	// //copy credentials to clipboard
	// creds, err := acloud.Copy(elems)
	// core.PrintIfErr(err)
	// core.Success("credentials : ", creds)

	// //DISPLAY WITH COLORS PROMINENTLY TO THE USER
	// acloud.DisplayCreds(creds)

	// //create string arrays of credentials
	// keys, vals := acloud.KeyVals(creds)
	// //create policies with map
	// policies, err := proxy.Policies(keys, vals)
	// core.PrintIfErr(err)
	// core.Success("policies : ", policies)


	// //ask if they want the credentials to be added to their aws config
	// if PromptConfig() == true {
	// 	//ask for path to aws config
	// 	// path := core.PromptFilePath()
	// 	// append aws creds to .aws/credentials file
	// 	err = core.AppendAwsCredentials(core.LocalCreds{
	// 		Path:      cliEnv.Aws_path,
	// 		User:      creds.User,
	// 		KeyID:     creds.KeyID,
	// 		AccessKey: creds.AccessKey,
	// 	})
	// 	core.PrintIfErr(err)
	// 	core.Success("aws credentials appended @ :", cliEnv.Aws_path)
	// }
}

type promptContent struct {
	Label string
	// Items    []string
	errorMsg string
}

type promptOptions struct {
	Label string
	Key   int64
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute()  {

	var p acloud.ACloudProvider

	options := []promptOptions{
		{
			Label: "Exit CLI",
			Key:   0,
		},{
			Label: "Get new Sandbox credentials",
			Key:   1,
		},
		{
			Label: "Download Text File of Sandbox Credentials",
			Key:   2,
		},
		{
			Label: "Append Sandbox Credentials to AWS Config",
			Key:   3,
		},
		{
			Label: "Display Sandbox Credentials",
			Key:   4,
		},
		{
			Label: "Set Credentials in GitHub Secret",
			Key:   5,
		},
	}

	templates := &promptui.SelectTemplates{
		Label:    "{{ . }}?",
		Active:   "\U0001F336 {{ .Label | cyan }} ",
		Inactive: "  {{ .Label | cyan }} ",
		Selected: "\U0001F336 {{ .Label | red | cyan }}",
	}

	searcher := func(input string, index int) bool {
		option := options[index]
		name := strings.Replace(strings.ToLower(option.Label), " ", "", -1)
		input = strings.Replace(strings.ToLower(input), " ", "", -1)

		return strings.Contains(name, input)
	}

	prompt := promptui.Select{
		Label:     "Welcome to GOSANDBOX - Please select an option",
		Items:     options,
		Templates: templates,
		// Size:      4,
		Searcher: searcher,
	}

	i, _, err := prompt.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return //os.Exit(1)
	}

	fmt.Printf("You choose number %d: %s\n", i+1, options[i].Label)

	switch options[i].Key {
	case 0:
		os.Exit(0)
	case 1:
		p.ACloudEnv, err = EnvLocation()
		core.PrintIfErr(err)
		core.Success("environment : ", p.ACloudEnv)
		p, err = GetSandboxCreds(p.ACloudEnv, p)
		core.PrintIfErr(err)
		core.Success("credentials : ", p.SandboxCredentials)
	case 2:
		// DownloadTextFile(p.SandboxCredentials)
	case 3:
		// download text file of policies
	case 4:
		// append aws creds to .aws/credentials file
	case 5: 
		//set credentials in github repo secrets
	}
	// return getEnv(".env")
	main()
}



func DownloadTextFile(creds acloud.SandboxCredentials)(){
	//if credentials are empty, return error
	if len(creds.AccessKey) == 0 || len(creds.KeyID) == 0 || len(creds.User) == 0 {
		fmt.Println("credentials are empty")
		return
	}

	//create string arrays of credentials
	keys, vals := acloud.KeyVals(creds)
	//create policies with map
	policies, err := proxy.Policies(keys, vals)
	core.PrintIfErr(err)
	core.Success("policies : ", policies)
	// ask if they want to download a text file with the credentials
	if PromptDownload() == true {
		// download text file of policies
		filename := PromptFileName()
		err = core.DocumentDownload(filename, policies)
		core.PrintIfErr(err)
	}
}

func GetSandboxCreds(cliEnv core.ACloudEnv, p *acloud.ACloudProvider) (acloud.ACloudProvider, error){

	//connect to website
	connect, err := core.Login(core.WebsiteLogin{Url: cliEnv.Url, Username: cliEnv.Username, Password: cliEnv.Password})
	core.PrintIfErr(err)
	// fmt.Println("connect : ", connect)
	core.Success("connection : ", connect)
	p.Connection = connect

	//scrape credentials
	elems, err := acloud.Sandbox(p.Connection, cliEnv.Download_key)
	core.PrintIfErr(err)
	core.Success("rod html elements : ", elems)

	//copy credentials to clipboard
	creds, err := acloud.Copy(elems)
	core.PrintIfErr(err)
	core.Success("credentials : ", creds)
	p.SandboxCredentials = creds

	//DISPLAY WITH COLORS PROMINENTLY TO THE USER
	acloud.DisplayCreds(creds)

	return p, err
}


func EnvLocation()(cliEnv core.ACloudEnv, err error) {
	options := []promptOptions{
		{
			Label: "Get sandbox credentials with .env file located in your current directory",
			Key:   1,
		}, {

			Label: "Get sandbox credentials from .env file in a custom location",
			Key:   2,
		}, {
			Label: "Get sandbox credentials manually with env information entered via cli prompt",
			Key:   3,
		},
	}

	templates := &promptui.SelectTemplates{
		Label:    "{{ . }}?",
		Active:   "\U0001F336 {{ .Label | cyan }} ",
		Inactive: "  {{ .Label | cyan }} ",
		Selected: "\U0001F336 {{ .Label | red | cyan }}",
	}

	searcher := func(input string, index int) bool {
		option := options[index]
		name := strings.Replace(strings.ToLower(option.Label), " ", "", -1)
		input = strings.Replace(strings.ToLower(input), " ", "", -1)

		return strings.Contains(name, input)
	}

	prompt := promptui.Select{
		Label:     "How would you like to input .env details: ",
		Items:     options,
		Templates: templates,
		// Size:      4,
		Searcher: searcher,
	}

	i, _, err := prompt.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return core.ACloudEnv{}, err
	}

	fmt.Printf("You choose number %d: %s\n", i+1, options[i].Label)

	switch options[i].Key {
	case 1:
		cliEnv, err = getEnv(".env")
	case 2:
		cliEnv, err = promptEnvFile()

	case 3:
		cliEnv, err = promptManual()
	}
	core.PrintIfErr(err)
	if err != nil {
		return core.ACloudEnv{}, err
	}

	return cliEnv, nil
}


func getEnv(env_path string) (core.ACloudEnv, error) {

	env, err := core.LoadEnvPath(env_path)
	if err != nil {
		fmt.Println("Could not load .env file - Err: ", err)
		promptEnvFile()
	}
	return env, nil

}

func promptEnvFile() (core.ACloudEnv, error) {
	//load env variables
	env_path := promptGetInput(promptContent{
		Label: "Please enter the path to the .env file from this directory",
	})

	env, err := getEnv(env_path)

	return env, err
}

func promptManual() (core.ACloudEnv, error) {

	tempEnv := core.ACloudEnv{}

	// get env vars via cli prompt
	tempEnv.Url = promptGetInput(promptContent{
		Label: "Name of web property URL you would like to login to",
	})
	// if tempEnv.Url == "" {
	// 	tempEnv.Url = "https://learn.acloud.guru/cloud-playground/cloud-sandboxes"
	// }
	tempEnv.Username = promptGetInput(promptContent{
		Label: "What is your username",
	})
	tempEnv.Password = promptGetInput(promptContent{
		Label: "What is your password",
	})
	// get aws_path via cli prompt
	tempEnv.Aws_path = promptGetInput(promptContent{
		Label: "Where would you like your sandbox credentials appended",
	})
	// get download path via cli prompt
	tempEnv.Download_key = promptGetInput(promptContent{
		Label: "What would you like the name of your sandbox credentials file to be",
	})
	//if all env vars are set, return the env
	if tempEnv.Url != "" && tempEnv.Username != "" && tempEnv.Password != "" && tempEnv.Aws_path != "" && tempEnv.Download_key != "" {
		return tempEnv, nil
	} else {
		fmt.Println("Please fill out all fields")
		promptManual()
	}
	return tempEnv, nil
}

func promptGetInput(pc promptContent) string {
	validate := func(input string) error {
		if len(input) <= 0 {
			return errors.New(pc.errorMsg)
		}
		return nil
	}
	prompt := promptui.Prompt{
		Label: pc.Label,
		// Templates: templates,
		Validate: validate,
	}
	result, err := prompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Input: %s\n", result)
	return result
}


func PromptDownload() bool{
	willDownload := promptGetInput(
		promptContent{Label: "Would you like to download the sandbox credentials file in plaintext? (yes/no)"})
	if willDownload == "yes" {
		fmt.Println("Downloading Sandbox Credentials...")
		return true		
	} else if willDownload == "no" {
		fmt.Println("Not downloading...")
		return false
	} else {
		fmt.Println("Invalid Answer")
		PromptDownload()
	}
	return false
}

func PromptFileName() string {
	filename := promptGetInput(promptContent{Label: "What would you like to name the file?"})
	return filename
}

func PromptFilePath() string {
	filepath := promptGetInput(promptContent{Label: "Where would you like to save the file to?"})
	return filepath
}


func PromptConfig() bool{
	willAppend := promptGetInput(
		promptContent{Label: "Would you like to append the sandbox credentials file to your AWS config file? (yes/no)"})
	if willAppend == "yes" {
		fmt.Println("Appending Sandbox Credentials to AWS configs...")
		return true

	} else if willAppend == "no" {
		fmt.Println("Not Appending to AWS configs...")
		return false
	} else {
		fmt.Println("Invalid Answer")
		PromptConfig()
	}
	return false
}
