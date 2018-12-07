package util

import (
	"fmt"
	"io/ioutil"
	"testing"
)

var pubKeyFile = "/home/sun/.config/rsa_pem/rsa_pub.pem"
var priKeyFile = "/home/sun/.config/rsa_pem/rsa_pri.pem"

var dstFile = "/tmp/hyper.dat"

func TestHyperEncrypt_Encrypt(t *testing.T) {
	pubKey, err := ioutil.ReadFile(pubKeyFile)
	if err != nil {
		panic(err)
	}
	priKey, err := ioutil.ReadFile(priKeyFile)
	if err != nil {
		panic(err)
	}
	hyperEncrypt := NewHyperEncryptB(pubKey, priKey)
	originData := "hellosun"
	encryptData, err := hyperEncrypt.Encrypt([]byte(originData))
	if err != nil {
		panic(err)
	}
	fmt.Println(string(encryptData))
	ioutil.WriteFile(dstFile, encryptData, 0766)
}

func TestHyperEncrypt_Decrypt(t *testing.T) {
	pubKey, err := ioutil.ReadFile(pubKeyFile)
	if err != nil {
		panic(err)
	}
	priKey, err := ioutil.ReadFile(priKeyFile)
	if err != nil {
		panic(err)
	}
	hyperEncrypt := NewHyperEncryptB(pubKey, priKey)

	cyperData, err := ioutil.ReadFile(dstFile)
	if err != nil {
		panic(err)
	}
	plainData, err := hyperEncrypt.Decrypt(cyperData)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(plainData))
}
