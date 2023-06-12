package dkr_nn

import (
	"fmt"

	dkrcred "github.com/docker/docker-credential-helpers/credentials"
)

type Helper struct {
	Encryption    Encryption
	CredsFilePath string
}

func (h Helper) Add(c *dkrcred.Credentials) error {
	errFailed := func(inner error) error { return fmt.Errorf("failed to add credential '%s': %w", c.ServerURL, inner) }
	credsFile, err := Load(h.CredsFilePath)
	if err != nil {
		return errFailed(err)
	}

	credsFile.SetCredentials(h.Encryption, c.ServerURL, c.Username, c.Secret)

	err = credsFile.Save(h.CredsFilePath)
	if err != nil {
		return errFailed(err)
	}
	return nil
}

func (h Helper) Delete(serverURL string) error {
	errFailed := func(inner error) error { return fmt.Errorf("failed to delete credential '%s': %w", serverURL, inner) }
	credsFile, err := Load(h.CredsFilePath)
	if err != nil {
		return errFailed(err)
	}

	credsFile.DelCredentials(serverURL)

	err = credsFile.Save(h.CredsFilePath)
	if err != nil {
		return errFailed(err)
	}
	return nil
}

func (h Helper) Get(serverURL string) (username string, password string, err error) {
	errFailed := func(inner error) error { return fmt.Errorf("failed to get credential '%s': %w", serverURL, inner) }
	credsFile, err := Load(h.CredsFilePath)
	if err != nil {
		return "", "", errFailed(err)
	}

	username, password, err = credsFile.GetCredentials(h.Encryption, serverURL)
	if err != nil {
		return "", "", errFailed(err)
	}
	return username, password, nil
}

func (h Helper) List() (map[string]string, error) {
	errFailed := func(inner error) error { return fmt.Errorf("failed to list credentials: %w", inner) }
	credsFile, err := Load(h.CredsFilePath)
	if err != nil {
		return nil, errFailed(err)
	}

	result := map[string]string{}
	for url := range credsFile.Credentials {
		result[url] = credsFile.Credentials[url].Username
	}
	return result, nil
}
