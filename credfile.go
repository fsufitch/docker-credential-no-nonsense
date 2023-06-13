package dkr_nn

import (
	"bytes"
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

	loadedFile := CredentialFile{Credentials: map[string]CredentialPair{}}

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
	stampedData := cf
	stampedData.Timestamp = time.Now()
	data, err := json.Marshal(stampedData)
	if err != nil {
		return fmt.Errorf("could not serialize cred file: %w", err)
	}
	prettyData := bytes.Buffer{}
	json.Indent(&prettyData, data, "", "\t")
	fmt.Fprintln(&prettyData)

	err = os.WriteFile(path, prettyData.Bytes(), 0644)
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

func (cf CredentialFile) SetCredentials(enc Encryption, key string, username string, password string) error {
	encryptedPassword, err := enc.Encrypt([]byte(password))
	if err != nil {
		return fmt.Errorf("set credentials failed (key=%v username=%v): %w", key, username, err)
	}
	credPair := CredentialPair{
		Username:          username,
		EncryptedPassword: B64Encoding.EncodeToString(encryptedPassword),
	}
	cf.Credentials[key] = credPair
	return nil
}

func (cf CredentialFile) DelCredentials(key string) {
	delete(cf.Credentials, key)
}
