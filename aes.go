package dkr_nn

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"fmt"
)

// Salt to apply to AES cipher generation
// From: dd if=/dev/urandom bs=1 count=64 | base64
const aesSalt = "dcOS3SD7EHiT4gZKh/OYDcWGzHbVD/TCaP4z21SBOR7iExE3enJ0MTa/ZIlvW0mgKMeeFH8wvBWAuvz3QOZkww"

const saltSplitter = "++salt:data++"

type Encryption struct {
	block cipher.Block
}

func NewEncryption(encryptionKey string) (*Encryption, error) {
	hash := sha256.New()
	fmt.Fprint(hash, aesSalt)
	fmt.Fprint(hash, encryptionKey)
	hashSum := hash.Sum(nil)

	// 32 bytes selects AES-256
	if len(hashSum) != 32 {
		return nil, errors.New("sha256 did not produce a 32-byte sum; this should be impossible")
	}

	ciph, err := aes.NewCipher(hashSum)
	if err != nil {
		return nil, fmt.Errorf("cipher creation failed: %w", err)
	}
	return &Encryption{ciph}, nil
}

func (enc Encryption) Encrypt(data []byte) []byte {
	salt := make([]byte, 16)
	rand.Read(salt)

	inputBuf := bytes.Buffer{}
	inputBuf.Write(salt)
	inputBuf.Write([]byte(saltSplitter))
	inputBuf.Write(data)

	output := []byte{}
	enc.block.Encrypt(output, inputBuf.Bytes())
	return output
}

func (enc Encryption) Decrypt(data []byte) ([]byte, error) {
	saltedOutput := []byte{}
	enc.block.Decrypt(saltedOutput, data)
	_, decryptedOutput, ok := bytes.Cut(saltedOutput, []byte(saltSplitter))
	if !ok {
		return nil, errors.New("input data did did not include a salt")
	}
	return decryptedOutput, nil
}

func (enc Encryption) Reencrypt(data []byte) ([]byte, error) {
	decrypted, err := enc.Decrypt(data)
	if err != nil {
		return nil, err
	}
	return enc.Encrypt(decrypted), nil
}
