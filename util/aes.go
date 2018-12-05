package util

import (
	"crypto/aes"
	"crypto/cipher"
)

type Aes struct {
	Key   []byte
	Nonce []byte
}

func (a *Aes) GcmEncrypt(plaintext []byte) ([]byte, error) {
	block, err := aes.NewCipher(a.Key)
	if err != nil {
		return nil, err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	return aesgcm.Seal(nil, a.Nonce, plaintext, nil), nil
}

func (a *Aes) GcmDecrypt(ciphertext []byte) ([]byte, error) {
	block, err := aes.NewCipher(a.Key)
	if err != nil {
		panic(err.Error())
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}

	return aesgcm.Open(nil, a.Nonce, ciphertext, nil)
}
