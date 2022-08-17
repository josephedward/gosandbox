package core

import (
	"errors"
	"fmt"
	"github.com/manifoldco/promptui"
	"os"
	"strings"
)

type promptContent struct {
	Label string
	// Items    []string
	errorMsg string
}

type promptOptions struct {
	Label string
	Key   int64
}

var Reset = "\033[0m"
var Red = "\033[31m"
var Green = "\033[32m"
var Yellow = "\033[33m"
var Blue = "\033[34m"
var Purple = "\033[35m"
var Cyan = "\033[36m"
var Gray = "\033[37m"
var White = "\033[97m"

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() (ACloudEnv, error) {

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
		Label:     "Welcome to GOSANDBOX. Please choose your .env file options: ",
		Items:     options,
		Templates: templates,
		// Size:      4,
		Searcher: searcher,
	}

	i, _, err := prompt.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return ACloudEnv{}, err
	}

	fmt.Printf("You choose number %d: %s\n", i+1, options[i].Label)

	switch options[i].Key {
	case 1:
		return getEnv(".env")
	case 2:
		return promptEnvFile()
	case 3:
		return promptManual()
	}
	return getEnv(".env")
}

func getEnv(env_path string) (ACloudEnv, error) {

	env, err := LoadEnvPath(env_path)
	if err != nil {
		fmt.Println("Could not load .env file - Err: ", err)
		promptEnvFile()
	}
	return env, nil

}

func promptEnvFile() (ACloudEnv, error) {
	//load env variables
	env_path := promptGetInput(promptContent{
		Label: "Please enter the path to the .env file from this directory",
	})

	env, err := getEnv(env_path)

	return env, err
}

func promptManual() (ACloudEnv, error) {

	tempEnv := ACloudEnv{}

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



func PrintIfErr(err error) {
	if err != nil {
		fmt.Println(Red + err.Error() + Reset)
	}
}

func Success(message ...interface{}) {
	for _, msg := range message {
		s, ok := msg.(string) // the "ok" boolean will flag success.
		if ok {
			fmt.Print(Green + string(s) + Reset)
		} else {
			fmt.Println(msg)
		}
	}
}


