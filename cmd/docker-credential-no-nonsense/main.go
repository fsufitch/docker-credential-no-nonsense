package main

import (
	"errors"
	"fmt"
	"os"
	"path"

	"github.com/adrg/xdg"
	"github.com/docker/docker-credential-helpers/credentials"
	dkr_nn "github.com/fsufitch/docker-credential-no-nonsense"
	"golang.org/x/exp/slices"
)

// This file implements the docker credentials protocol

const NO_NONSENSE_CREDFILE = "NO_NONSENSE_CREDFILE"
const NO_NONSENSE_ENC_KEY = "NO_NONSENSE_ENC_KEY"

func getDefaultCredFilePath() (string, error) {
	return xdg.DataFile("dkr-no-nonsense-credfile.json")
}

func getEncryptionKey() (string, error) {
	encKey := os.Getenv(NO_NONSENSE_ENC_KEY)
	if encKey == "" {
		return "", errors.New("NO_NONSENSE_ENC_KEY not set; it is required for interacting with the credfile")
	}
	return encKey, nil
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

var executable string

func init() {
	executable, _ = os.Executable()
	executable = path.Base(executable)
}

func main() {

	path, err := getCredFilePath()
	if err != nil {
		fmt.Fprintln(os.Stderr, fmt.Errorf("%s: %w", executable, err))
		os.Exit(1)
	}

	if slices.Contains(os.Args, "--where") {
		fmt.Println(path)
		return
	}

	enc, err := initEncryption()
	if err != nil {
		fmt.Fprintln(os.Stderr, fmt.Errorf("%s: %w", executable, err))
		os.Exit(1)
	}

	helper := dkr_nn.Helper{
		Encryption:    *enc,
		CredsFilePath: path,
	}

	credentials.Serve(helper)
}
