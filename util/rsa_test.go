package util

import (
	"fmt"
	"testing"
)

func Test0(t *testing.T) {
	var data []byte
	var err error
	r := &Rsa{}
	data, err = r.RsaEncrypt([]byte("fyxichen"))
	if err != nil {
		panic(err)
	}
	origData, err := r.RsaDecrypt(data)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(origData))
}
