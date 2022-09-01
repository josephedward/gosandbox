package gh

import (
	"encoding/base64"
	"errors"
	"flag"
	"log"
	sodium "github.com/GoKillers/libsodium-go/cryptobox"
	"github.com/google/go-github/v47/github"
	"golang.org/x/oauth2"
	"os"
	"github.com/joho/godotenv"
	"context"
	"fmt"
	"gosandbox/core"
)

var (
	repo  = flag.String("repo", "", "The repo that the secret should be added to, ex. go-github")
	owner = flag.String("owner", "", "The owner of there repo this should be added to, ex. google")
)



func GetRepositories() {
	err := godotenv.Load("./gh/.env")
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("TOKEN")},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	core.Success("client: ",client)	

	// list all repositories for the authenticated user
	repos, _, err := client.Repositories.List(ctx, "josephedward", nil)

	core.Success("repos: ",repos)

	if err != nil {
		fmt.Println(err)
	}

}



func SecretEnv()(){
	flag.Parse()
	err := godotenv.Load("../gh/.env")
	core.PrintIfErr(err)
	token := os.Getenv("TOKEN")
	if token == "" {
		log.Fatal("please provide a GitHub API token via env variable GITHUB_AUTH_TOKEN")
	}

	if *repo == "" {
		log.Fatal("please provide required flag --repo to specify GitHub repository ")
	}

	if *owner == "" {
		log.Fatal("please provide required flag --owner to specify GitHub user/org owner")
	}

	secretName, err := getSecretName()
	if err != nil {
		log.Fatal(err)
	}

	secretValue, err := getSecretValue(secretName)
	if err != nil {
		log.Fatal(err)
	}

	ctx, client, err := githubAuth(token)
	if err != nil {
		log.Fatalf("unable to authorize using env GITHUB_AUTH_TOKEN: %v", err)
	}

	if err := addRepoSecret(ctx, client, *owner, *repo, secretName, secretValue); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Added secret %q to the repo %v/%v\n", secretName, *owner, *repo)
}

func getSecretName() (string, error) {
	secretName := flag.Arg(0)
	fmt.Println("secretName: ",secretName)
	if secretName == "" {
		return "", fmt.Errorf("missing argument secret name")
	}
	return secretName, nil
}

func getSecretValue(secretName string) (string, error) {
	// secretValue := os.Getenv(secretName)
	secretValue := flag.Arg(1)
	if secretValue == "" {
		return "", fmt.Errorf("missing argument secret value")
	}
	return secretValue, nil
}

// githubAuth returns a GitHub client and context.
func githubAuth(token string) (context.Context, *github.Client, error) {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)
	return ctx, client, nil
}


func addRepoSecret(ctx context.Context, client *github.Client, owner string, repo, secretName string, secretValue string) error {
	publicKey, _, err := client.Actions.GetRepoPublicKey(ctx, owner, repo)
	if err != nil {
		return err
	}

	encryptedSecret, err := encryptSecretWithPublicKey(publicKey, secretName, secretValue)
	if err != nil {
		return err
	}

	if _, err := client.Actions.CreateOrUpdateRepoSecret(ctx, owner, repo, encryptedSecret); err != nil {
		return fmt.Errorf("Actions.CreateOrUpdateRepoSecret returned error: %v", err)
	}

	return nil
}

func encryptSecretWithPublicKey(publicKey *github.PublicKey, secretName string, secretValue string) (*github.EncryptedSecret, error) {
	decodedPublicKey, err := base64.StdEncoding.DecodeString(publicKey.GetKey())
	if err != nil {
		return nil, fmt.Errorf("base64.StdEncoding.DecodeString was unable to decode public key: %v", err)
	}

	secretBytes := []byte(secretValue)
	encryptedBytes, exit := sodium.CryptoBoxSeal(secretBytes, decodedPublicKey)
	if exit != 0 {
		return nil, errors.New("sodium.CryptoBoxSeal exited with non zero exit code")
	}

	encryptedString := base64.StdEncoding.EncodeToString(encryptedBytes)
	keyID := publicKey.GetKeyID()
	encryptedSecret := &github.EncryptedSecret{
		Name:           secretName,
		KeyID:          keyID,
		EncryptedValue: encryptedString,
	}
	return encryptedSecret, nil
}