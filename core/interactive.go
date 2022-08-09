package core

import (
	"errors"
	"fmt"
	"os"
	"github.com/manifoldco/promptui"
)

type promptContent struct {
	Label    string
	Items    []string
	errorMsg string
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() (ACloudEnv, error) {

	choice := promptGetInput(promptContent{
		Label: "Welcome to the ACloud Sandbox Procurer. Would you like to use an .env file or enter your credentials manually (env/manual) ",
	})
	tempEnv := ACloudEnv{}
	if choice == "env" {
		//get env vars via path from prompt
		return promptEnvFile(tempEnv)
	} else if choice == "manual" {
		// get env vars via cli prompt
		return promptManual(tempEnv)
	} else {
		fmt.Println("Invalid choice. Please try again.")
		Execute()
	}
	return ACloudEnv{}, errors.New("Error")
}

func promptEnvFile(tempEnv ACloudEnv) (ACloudEnv, error) {
	//load env variables
	env_path := promptGetInput(promptContent{
		Label: "Please enter the path to the .env file from this directory",
	})
	// if env_path == "" {
	// 	env_path = ".env"
	// }
	tempEnv, err := LoadEnvPath(env_path)
	if err != nil {
		fmt.Println("Could not load .env file - Err: ", err)
		promptEnvFile(tempEnv)
	}
	return tempEnv, nil
}

func promptManual(tempEnv ACloudEnv) (ACloudEnv, error) {
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
		promptManual(tempEnv)
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


func PrintIfErr(err error) {
	if err != nil {
		fmt.Println(err.Error())
	}
}