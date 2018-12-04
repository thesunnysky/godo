package util

import (
	"fmt"
	"testing"
)

func Test0(t *testing.T) {
	var data []byte
	var err error
	data, err = RsaEncrypt([]byte("fyxichen"))
	if err != nil {
		panic(err)
	}
	origData, err := RsaDecrypt(data)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(origData))
}
