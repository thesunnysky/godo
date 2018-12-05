package util

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
)

type Aes struct {
	Key   string
	Nonce string
}

func (conf *Aes) GcmEncrypt(plaintext []byte) ([]byte, error) {
	keyBytes, _ := hex.DecodeString(conf.Key)

	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		return nil, err
	}

	nonceBytes, _ := hex.DecodeString(conf.Nonce)

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	return aesgcm.Seal(nil, nonceBytes, plaintext, nil), nil
}

func (conf *Aes) GcmDecrypt(ciphertext []byte) ([]byte, error) {
	keyBytes, _ := hex.DecodeString(conf.Key)
	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		panic(err.Error())
	}

	nonceBytes, _ := hex.DecodeString(conf.Nonce)

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}

	return aesgcm.Open(nil, nonceBytes, ciphertext, nil)
}
