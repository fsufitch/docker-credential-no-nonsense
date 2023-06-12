package main

import (
	"fmt"
	"os"

	"github.com/adrg/xdg"
	dkrcred "github.com/docker/docker-credential-helpers/credentials"
	dkr_nn "github.com/fsufitch/docker-credential-no-nonsense"
	"github.com/manifoldco/promptui"
)

// This file implements the docker credentials protocol

const NO_NONSENSE_CREDFILE = "NO_NONSENSE_CREDFILE"
const NO_NONSENSE_ENC_KEY = "NO_NONSENSE_ENC_KEY"

func getDefaultCredFilePath() (string, error) {
	return xdg.DataFile("dkr-no-nonsense-credfile.json")
}

func getEncryptionKey() (string, error) {
	encKey := os.Getenv(NO_NONSENSE_ENC_KEY)
	if encKey != "" {
		return encKey, nil
	}
	fmt.Fprintln(os.Stderr, "No password in NO_NONSENSE_ENC_KEY variable")
	prompt := promptui.Prompt{
		Label:       "Encryption Key: ",
		Mask:        '.',
		HideEntered: true,
	}
	return prompt.Run()
}

func getCredFilePath() (string, error) {
	path := os.Getenv(NO_NONSENSE_CREDFILE)
	if path != "" {
		return path, nil
	}
	path, err := getDefaultCredFilePath()
	if err != nil {
		return "", fmt.Errorf("could not resolve cred file path: %w", err)
	}
	return path, nil
}

func initEncryption() (*dkr_nn.Encryption, error) {
	errFailed := func(inner error) error { return fmt.Errorf("failed to initialize encryption: %w", inner) }
	encKey, err := getEncryptionKey()
	if err != nil {
		return nil, errFailed(err)
	}
	enc, err := dkr_nn.NewEncryption(encKey)
	if err != nil {
		return nil, errFailed(err)
	}
	return enc, nil
}

func main() {
	enc, err := initEncryption()
	if err != nil {
		panic(err)
	}

	path, err := getCredFilePath()
	if err != nil {
		panic(err)
	}

	helper := dkr_nn.Helper{
		Encryption:    *enc,
		CredsFilePath: path,
	}

	dkrcred.Serve(helper)

}
