package util

import (
	"bytes"
	"crypto/rand"
	"errors"
	"github.com/CodisLabs/codis/pkg/utils/log"
	"io"
)

type HyperEncrypt struct {
	Rsa Rsa
	Aes Aes
}

func NewHyperEncrypt(rsaPublicKey, rsaPrivateKey []byte) *HyperEncrypt {
	rasConf := Rsa{PublicKey: rsaPublicKey, PrivateKey: rsaPrivateKey}
	aesKey, aesNonce := genAesKeyAndNonce()
	aesConf := Aes{Key: aesKey, Nonce: aesNonce}
	hyperEncryptConf := &HyperEncrypt{Rsa: rasConf, Aes: aesConf}
	return hyperEncryptConf
}

func (hyper *HyperEncrypt) Encrypt(plaintext []byte) ([]byte, error) {
	cipherText, err := hyper.Aes.GcmEncrypt(plaintext)
	if err != nil {
		log.Errorf("Aes encrypt data error:%s\n", err)
		return nil, err
	}

	cipherAesKey, err := hyper.Rsa.RsaEncrypt(hyper.Aes.Key)
	if err != nil {
		log.Errorf("Rsa encrypt data error:%s\n", err)
		return nil, err
	}

	cipherAesNonce, err := hyper.Rsa.RsaEncrypt(hyper.Aes.Nonce)
	if err != nil {
		log.Errorf("Rsa encrypt data error:%s\n", err)
		return nil, err
	}

	var buffer bytes.Buffer
	buffer.Write(cipherAesKey)
	buffer.WriteByte('\n')
	buffer.Write(cipherAesNonce)
	buffer.WriteByte('\n')
	buffer.Write(cipherText)
	return buffer.Bytes(), nil
}

func (hyper *HyperEncrypt) Decrypt(cipherData []byte) ([]byte, error) {
	splitData := bytes.Split(cipherData, []byte{'\n'})
	if len(splitData) != 3 {
		return nil, errors.New("the composition of cipher data is invalid")
	}
	cipherAesKey := splitData[0]
	cipherAesNonce := splitData[1]
	cipherText := splitData[2]
	aesKey, err := hyper.Rsa.RsaDecrypt(cipherAesKey)
	if err != nil {
		return nil, err
	}

	aesNonce, err := hyper.Rsa.RsaDecrypt(cipherAesNonce)
	if err != nil {
		return nil, err
	}

	hyper.Aes.Key = aesKey
	hyper.Aes.Nonce = aesNonce
	plainText, err := hyper.Aes.GcmDecrypt(cipherText)
	if err != nil {
		return nil, err
	}
	return plainText, nil
}

func genAesKeyAndNonce() (key, nonce []byte) {
	key = make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, key); err != nil {
		panic(err.Error())
	}

	nonce = make([]byte, 12)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		panic(err.Error())
	}
	return
}
