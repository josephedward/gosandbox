package gh

import (
	"context"
	crypto_rand "crypto/rand"
	"encoding/base64"
	"flag"
	"fmt"
	"os"
	"github.com/google/go-github/v47/github"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/nacl/box"
	"golang.org/x/oauth2"
	"gosandbox/cli"
)

func GetSecretName() (string, error) {
	secretName := flag.Arg(0)
	if secretName == "" {
		return "", fmt.Errorf("missing argument secret name")
	}
	return secretName, nil
}

func GetSecretValue(secretName string) (string, error) {
	secretValue := os.Getenv(secretName)
	if secretValue == "" {
		return "", fmt.Errorf("secret value not found under env variable %q", secretName)
	}
	return secretValue, nil
}

// GithubAuth returns a GitHub client and context.
func GithubAuth(token string) (context.Context, *github.Client, error) {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)
	return ctx, client, nil
}

func AddRepoSecret(ctx context.Context, client *github.Client, owner string, repo, secretName string, secretValue string) error {
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

	var boxKey [32]byte
	copy(boxKey[:], decodedPublicKey)
	secretBytes := []byte(secretValue)
	encryptedBytes, err := box.SealAnonymous([]byte{}, secretBytes, &boxKey, crypto_rand.Reader)
	if err != nil {
		return nil, fmt.Errorf("box.SealAnonymous failed with error %w", err)
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

var (
	repo  = flag.String("repo", "", "The repo that the secret should be added to, ex. go-github")
	owner = flag.String("owner", "", "The owner of there repo this should be added to, ex. google")
)

func GetToken() (token string, err error) {
	err = godotenv.Load("./gh/.env")
	cli.PrintIfErr(err)
	token = os.Getenv("TOKEN")
	if token == "" {
		cli.Error("please provide a GitHub API token via env variable TOKEN")
	}

	return token, err
}

func SecretEnv() {
	flag.Parse()

	secretName, err := GetSecretName()
	secretValue, err := GetSecretValue(flag.Arg(1))

	token, err := GetToken()
	if err != nil {
		cli.Error(err)
	}

	ctx, client, err := GithubAuth(token)
	if err != nil {
		cli.Error("unable to authorize using env TOKEN: %v", err)
	}

	if *repo == "" {
		cli.Error("please provide required flag --repo to specify GitHub repository ")
	}

	if *owner == "" {
		cli.Error("please provide required flag --owner to specify GitHub user/org owner")
	}

	if err := AddRepoSecret(ctx, client, *owner, *repo, secretName, secretValue); err != nil {
		cli.Error(err)
	}

	fmt.Printf("Added secret %q to the repo %v/%v\n", secretName, *owner, *repo)
}
