package main

import (
	"fmt"
	"gosandbox/acloud"
	"gosandbox/cli"
	"gosandbox/core"
	"gosandbox/gh"
	"gosandbox/proxy"
	// "log"
	"os"
	"strings"

	"github.com/manifoldco/promptui"
)

func main() {
	var p acloud.ACloudProvider
	Execute(p)
}


func GetTemplates() *promptui.SelectTemplates {
	templates := &promptui.SelectTemplates{
		Label:    "{{ . }}",
		Active:   "\U0001F336 {{ .Label | yellow }} ",
		Inactive: "  {{ .Label | cyan }} ",
		Selected: "\U0001FAD1 {{ .Label | green | cyan }}",
	}
	return templates
}

func GetSearcher(options []cli.PromptOptions) func(input string, index int) bool {
	searcher := func(input string, index int) bool {
		option := options[index]
		name := strings.Replace(strings.ToLower(option.Label), " ", "", -1)
		input = strings.Replace(strings.ToLower(input), " ", "", -1)

		return strings.Contains(name, input)
	}
	return searcher
}

func Select(promptTitle string, options []cli.PromptOptions) *promptui.Select {
	prompt := promptui.Select{
		Label:     promptTitle,
		Items:     options,
		Templates: GetTemplates(),
		Size:      4,
		Searcher:  GetSearcher(options),
	}
	return &prompt
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute(p acloud.ACloudProvider) {
	
	//if p doesnt have ACloudEnv set, automatically load it
	fmt.Println("p.ACloudEnv : ",p.ACloudEnv)
	if len(p.ACloudEnv.Url) == 0 {
		env, err := cli.GetEnv(".env")
		cli.PrintIfErr(err)
		p.ACloudEnv = env
		//get sandbox creds
		p, err = GetSandboxCreds(p.ACloudEnv, &p)
	}

	options := []cli.PromptOptions{
		{
			Label: "Exit CLI",
			Key:   0,
		}, {
			Label: "Scrape Sandbox Credentials",
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
		{
			Label: "Open AWS Console for Sandbox",
			Key:   6,
		},
	}

	prompt := Select("Welcome to GOSANDBOX - Please select an option: ", options)

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
		cli.PrintIfErr(err)
		cli.Success("environment : ", p.ACloudEnv)
		p, err = GetSandboxCreds(p.ACloudEnv, &p)
		cli.PrintIfErr(err)
		cli.Success("credentials : ", p.SandboxCredentials)
	case 2:
		// download text file of policies
		DownloadTextFile(p.SandboxCredentials)
	case 3:
		// append aws creds to .aws/credentials file
		AppendCreds(p.SandboxCredentials)
	case 4:
		//DISPLAY WITH COLORS PROMINENTLY TO THE USER
		acloud.DisplayCreds(p.SandboxCredentials)
	case 5:
		//set sandbox creds in github secret
		SandboxToGithub(p.SandboxCredentials)
	case 6:
		//open aws console for sandbox
		OpenAWSConsole(p.SandboxCredentials)
	}
	Execute(p)
}

func OpenAWSConsole(creds acloud.SandboxCredentials) {
	//if credentials are empty, return error
	if len(creds.AccessKey) == 0 || len(creds.KeyID) == 0 || len(creds.User) == 0 {
		cli.Error("Warning: credentials are empty")
		return
	}

	//login to console with credentials url, username, and password
	browser, err := core.Login(core.WebsiteLogin{creds.URL, creds.User, creds.Password})
	cli.PrintIfErr(err)
	cli.Success("browser : ",browser)
}

func SandboxToGithub(creds acloud.SandboxCredentials) {
	//github PAT
	token, err := gh.GetToken()
	cli.PrintIfErr(err)

	//if token is empty, return error
	if len(token) == 0 {
		cli.Error("Github token is empty")
		return
	}

	//if credentials are empty, return error
	if len(creds.AccessKey) == 0 || len(creds.KeyID) == 0 || len(creds.User) == 0 {
		cli.Error("credentials are empty")
		return
	}

	// authorize using env TOKEN
	ctx, client, err := gh.GithubAuth(token)
	cli.PrintIfErr(err)

	// get repo owner
	owner, err := cli.PromptRepoOwner()

	// get repo name
	repo, err := cli.PromptRepoName()

	//create string arrays of credentials
	keys, vals := acloud.KeyVals(creds)

	cli.Success("writing credentials to github secrets....")
	//loop over keys and vals
	for i, key := range keys {
		//create secret in github
		// err := gh.CreateSecret(key, vals[i])
		if err := gh.AddRepoSecret(ctx, client, owner, repo, key, vals[i]); err != nil {
			cli.PrintIfErr(err)
		}
		cli.Success("secret : "+key)
		fmt.Println("value : "+vals[i])
	}
	cli.Success("credentials written to "+owner+"/"+repo)
}

func DownloadTextFile(creds acloud.SandboxCredentials) {
	//if credentials are empty, return error
	if len(creds.AccessKey) == 0 || len(creds.KeyID) == 0 || len(creds.User) == 0 {
		cli.Error("Warning: credentials are empty")
		return
	}

	//create string arrays of credentials
	keys, vals := acloud.KeyVals(creds)
	//create policies with map
	policies, err := proxy.Policies(keys, vals)
	cli.PrintIfErr(err)
	cli.Success("policies : ", policies)
	// ask if they want to download a text file with the credentials
	if cli.PromptDownload() == true {
		// download text file of policies
		filename := cli.PromptFileName()
		err = core.DocumentDownload(filename, policies)
		cli.PrintIfErr(err)
	}
}

func AppendCreds(creds acloud.SandboxCredentials) {
	//ask if they want the credentials to be added to their aws config
	path := cli.PromptGetInput(cli.PromptContent{
		Label: "Where would you like your sandbox credentials appended?",
	})

	//if credentials are empty, return error
	if len(creds.AccessKey) == 0 || len(creds.KeyID) == 0 || len(creds.User) == 0 {
		cli.Error("Warning: credentials are empty")
		return
	}

	if cli.PromptConfig() == true {
		//ask for path to aws config
		// path := core.PromptFilePath()
		// append aws creds to .aws/credentials file
		err := core.AppendAwsCredentials(core.LocalCreds{
			Path:      path,
			User:      creds.User,
			KeyID:     creds.KeyID,
			AccessKey: creds.AccessKey,
		})

		// if error, ask for path to aws config again
		if err != nil {
			cli.PrintIfErr(err)
			AppendCreds(creds)
		}
		cli.Success("aws credentials appended @ :", path)
	}
}

func GetSandboxCreds(cliEnv core.ACloudEnv, p *acloud.ACloudProvider) (acloud.ACloudProvider, error) {

	//connect to website
	connect, err := core.Login(core.WebsiteLogin{Url: cliEnv.Url, Username: cliEnv.Username, Password: cliEnv.Password})
	cli.PrintIfErr(err)
	// fmt.Println("connect : ", connect)
	cli.Success("connection : ", connect)
	p.Connection = connect

	//scrape credentials
	elems, err := acloud.Sandbox(p.Connection, cliEnv.Download_key)
	cli.PrintIfErr(err)
	// cli.Success("rod html elements : ", elems)

	//copy credentials to clipboard
	creds, err := acloud.Copy(elems)
	cli.PrintIfErr(err)
	// cli.Success("credentials : ", creds)
	p.SandboxCredentials = creds

	//DISPLAY WITH COLORS PROMINENTLY TO THE USER
	acloud.DisplayCreds(creds)

	return *p, err
}

func EnvLocation() (cliEnv core.ACloudEnv, err error) {
	options := []cli.PromptOptions{
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

	prompt := Select("How would you like to input .env details: ", options)

	i, _, err := prompt.Run()

	if err != nil {
		cli.PrintIfErr(err)
		return core.ACloudEnv{}, err
	}
	cli.Success("You choose number %d: %s\n", i+1, options[i].Label)

	switch options[i].Key {
	case 1:
		cliEnv, err = cli.GetEnv(".env")
	case 2:
		cliEnv, err = cli.PromptEnvFile()

	case 3:
		cliEnv, err = cli.PromptManual()
	}
	cli.PrintIfErr(err)
	if err != nil {
		return core.ACloudEnv{}, err
	}

	return cliEnv, nil
}
