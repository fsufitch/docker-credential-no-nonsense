package dkr_nn

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/adrg/xdg"
)

var CREDFILE string

var B64Encoding = base64.StdEncoding

func init() {
	path, err := xdg.DataFile("dkr-no-nonsense-credfile.json")
	if err != nil {
		panic(err)
	}
	CREDFILE = path
}

type CredentialFile struct {
	Timestamp   time.Time                 `json:"timestamp"`
	Credentials map[string]CredentialPair `json:"credentials"`
}

type CredentialPair struct {
	Username          string `json:"username"`
	EncryptedPassword string `json:"enc_password"`
}

func Load(path string) (*CredentialFile, error) {
	if path == "" {
		path = CREDFILE
	}

	loadedFile := CredentialFile{}

	byts, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		// The file missing is OK, we will just use defaults
		return &loadedFile, nil
	}
	if err != nil {
		return nil, fmt.Errorf("could not read credential file '%v': %w", path, err)
	}

	err = json.Unmarshal(byts, &loadedFile)
	if err != nil {
		return nil, fmt.Errorf("failed to load credential file '%v': %w", path, err)
	}

	return &loadedFile, nil
}

func (cf CredentialFile) Save(path string) error {
	data, err := json.Marshal(cf)
	if err != nil {
		return fmt.Errorf("could not serialize cred file: %w", err)
	}

	err = os.WriteFile(path, data, 0)
	if err != nil {
		return fmt.Errorf("failed to write creds file '%s': %w", path, err)
	}
	return nil
}

func (cf CredentialFile) GetCredentials(enc Encryption, key string) (username string, password string, err error) {
	credPair, ok := cf.Credentials[key]
	if !ok {
		err = fmt.Errorf("key not in credfile: %s", key)
		return
	}

	username = credPair.Username
	encryptedPassword, err := B64Encoding.DecodeString(credPair.EncryptedPassword)
	if err != nil {
		return
	}
	decryptedPassword, err := enc.Decrypt(encryptedPassword)
	if err != nil {
		return
	}
	password = string(decryptedPassword)
	return
}

func (cf CredentialFile) SetCredentials(enc Encryption, key string, username string, password string) {
	encryptedPassword := enc.Encrypt([]byte(password))
	credPair := CredentialPair{
		Username:          username,
		EncryptedPassword: B64Encoding.EncodeToString(encryptedPassword),
	}
	cf.Credentials[key] = credPair
}

func (cf CredentialFile) DelCredentials(key string) {
	delete(cf.Credentials, key)
}
