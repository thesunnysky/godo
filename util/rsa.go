package util

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
)

var privateKey, publicKey []byte
var publicKeyFile, privateKeyFile string

/*func init() {
	//publicKeyFile = config.ClientConfig.PublicKeyFile
	//privateKeyFile = config.ClientConfig.PrivateKeyFile

	var err error
	publicKey, err = ioutil.ReadFile(publicKeyFile)
	if err != nil {
		panic(err)
	}
	privateKey, err = ioutil.ReadFile(privateKeyFile)
	if err != nil {
		panic(err)
	}
}*/

// 加密
func RsaEncrypt(origData []byte) ([]byte, error) {
	block, _ := pem.Decode(publicKey)
	if block == nil {
		return nil, errors.New("public key error")
	}
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	pub := pubInterface.(*rsa.PublicKey)
	return rsa.EncryptPKCS1v15(rand.Reader, pub, origData)
}

// 解密
func RsaDecrypt(cipherText []byte) ([]byte, error) {
	block, _ := pem.Decode(privateKey)
	if block == nil {
		return nil, errors.New("private key error!")
	}
	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return rsa.DecryptPKCS1v15(rand.Reader, priv, cipherText)
}
