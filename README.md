# GOSANDBOX

This project helps provide AWS credential information for user's of A Cloud Guru's sandbox [feature](https://help.acloud.guru/hc/en-us/articles/360001389256-AWS-Cloud-Sandbox). It obtains credentials for a Cloud Sandbox using environment variables ([.env](https://github.com/joho/godotenv)) from a file located in the root directory, outputs credentials and snapshots to the /test folder, and appends AWS-CLI credentials to the assumed default location in macos and most linux-based filesystems.

## A Cloud Guru Sandboxes

Sandboxes are temporary AWS environments that last for 4 hours and contain emphemeral resources that do not accrue costs in your AWS acccount ([FAQ](https://help.acloud.guru/hc/en-us/articles/360001477955-Cloud-Playground-FAQ#h_91f2897f-8feb-40a6-aeac-374c51c927c5)). 

## Interface

[Interfaces](https://golangdocs.com/interfaces-in-golang) are contracts that have a set of functions to be implemented to fulfill that contract

## Rod

[Rod](go-rod.github.io) is a high-level Chrome DevTools Protocol driver used for web automation and scraping. 

## Taskfile Commands

This project uses a [taskfile](https://taskfile.dev/) to run tasks:

- `cli`:  
uses a visible browser for ease of debugging, and runs a script with a CLI prompt for entering credentials:
    - with an .env file in a custom location 
    - manually, entering variables through golang string [prompt](https://github.com/manifoldco/promptui)

- `provider`:  
this project was initially build as a demo implementation of the `PolicyProvider` [interface](./proxy/README.md) - this commands runs that provider

- `get-creds-visible`: 
obtains credentials for a Cloud Sandbox using environment variables, except it uses a visible browser for debugging

- `test`:  
runs all integration and unit tests


## Folders 

### acloud

- **provider.go** : interface-specific implementation methods
- **sandbox.go** : contains methods specifically pertaining to crawling the acloudguru website for credential information 

### core

- **connection.go** : generalized structs and methods for logging into a website
- **env.go** :  loads a .env file with acloud-specific usage information into the go runtime
- **export.go**: code for writing credentials to a text file, appending them to a CLI config, and screenshotting the browser window
- **interactive.go**: code for obtaining CLI inputs for .env file 

### proxy

- **carrierproxy.go**: original interface for GloveBox demo, Policy struct and Policies method that returns an object of type Policy with keys and vals as CarrierIDs and PolicyNumbers

### scripts

- **getawscreds.go**: contains the initial POC script that was used to scrape the site at first 
- **interactive.go**: script that initiates the interactive CLI 

### Tests

Integration
- **methods_test.go** : integration tests of separate methods 
- **provider_test.go** : integration test for provider(and base script for obtaining credentials)

Unit
- **acloud_test.go** : unit tests for the /acloud folder code
- **core_test.go** : unit tests for the /core folder code
- **proxy_test.go** : unit tests for the /proxy folder code


### Creating Tests 

*alias gotests=~/go/bin/gotests*



## Brew Package
Install (macos and linux): 
```
brew tap josephedward/homebrew-aptnative
brew install gosandbox
```