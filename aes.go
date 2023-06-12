package dkr_nn

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
)

// Salt to apply to AES cipher generation
// From: dd if=/dev/urandom bs=1 count=64 | base64
const aesSalt = "dcOS3SD7EHiT4gZKh/OYDcWGzHbVD/TCaP4z21SBOR7iExE3enJ0MTa/ZIlvW0mgKMeeFH8wvBWAuvz3QOZkww"

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

func (enc Encryption) Encrypt(data []byte) ([]byte, error) {
	iv := enc.makeIV()

	src := bytes.NewReader(data)
	dest := bytes.Buffer{}
	dest.Write(iv)

	writer := cipher.StreamWriter{W: &dest, S: cipher.NewOFB(enc.block, iv)}
	if _, err := io.Copy(writer, src); err != nil {
		return nil, fmt.Errorf("encryption failed: %w", err)
	}
	return dest.Bytes(), nil
}

func (enc Encryption) Decrypt(data []byte) ([]byte, error) {
	iv := data[:enc.block.BlockSize()]
	data = data[enc.block.BlockSize():]

	src := bytes.NewReader(data)
	reader := cipher.StreamReader{R: src, S: cipher.NewOFB(enc.block, iv)}

	decryptedData, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("decryption failed: %w", err)
	}
	return decryptedData, nil
}

func (enc Encryption) Reencrypt(data []byte) ([]byte, error) {
	decrypted, err := enc.Decrypt(data)
	if err != nil {
		return nil, err
	}
	return enc.Encrypt(decrypted)
}

func (enc Encryption) makeIV() []byte {
	iv := make([]byte, enc.block.BlockSize())
	rand.Read(iv)
	return iv
}
