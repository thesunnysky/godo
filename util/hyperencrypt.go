package util

import (
	"bytes"
	"crypto/rand"
	"encoding/binary"
	"github.com/CodisLabs/codis/pkg/utils/log"
	"io"
	"io/ioutil"
)

type HyperEncrypt struct {
	rsa Rsa
	aes Aes
}

func NewHyperEncryptB(rsaPublicKey, rsaPrivateKey []byte) *HyperEncrypt {
	rasConf := Rsa{PublicKey: rsaPublicKey, PrivateKey: rsaPrivateKey}
	aesKey, aesNonce := genAesKeyAndNonce()
	aesConf := Aes{Key: aesKey, Nonce: aesNonce}
	hyperEncryptConf := &HyperEncrypt{rsa: rasConf, aes: aesConf}
	return hyperEncryptConf
}

func NewHyperEncryptF(rsaPublicKeyFile, rsaPrivateKeyFile string) (*HyperEncrypt, error) {
	rsaPublicKey, err := ioutil.ReadFile(rsaPublicKeyFile)
	if err != nil {
		return nil, err
	}

	rsaPrivateKey, err := ioutil.ReadFile(rsaPrivateKeyFile)
	if err != nil {
		return nil, err
	}
	return NewHyperEncryptB(rsaPublicKey, rsaPrivateKey), nil
}

func (hyper *HyperEncrypt) Encrypt(plaintext []byte) ([]byte, error) {
	cipherText, err := hyper.aes.GcmEncrypt(plaintext)
	if err != nil {
		log.Errorf("aes encrypt data error:%s\n", err)
		return nil, err
	}

	cipherAesKey, err := hyper.rsa.RsaEncrypt(hyper.aes.Key)
	if err != nil {
		log.Errorf("rsa encrypt data error:%s\n", err)
		return nil, err
	}

	cipherAesNonce, err := hyper.rsa.RsaEncrypt(hyper.aes.Nonce)
	if err != nil {
		log.Errorf("rsa encrypt data error:%s\n", err)
		return nil, err
	}

	var buffer bytes.Buffer
	keySize := len(cipherAesKey)
	keybs := make([]byte, 4)
	binary.BigEndian.PutUint32(keybs, uint32(keySize))
	buffer.Write(keybs)
	buffer.Write(cipherAesKey)

	nonceSize := len(cipherAesNonce)
	noncebs := make([]byte, 4)
	binary.BigEndian.PutUint32(noncebs, uint32(nonceSize))
	buffer.Write(keybs)
	buffer.Write(cipherAesNonce)

	dataSize := len(cipherText)
	databs := make([]byte, 4)
	binary.BigEndian.PutUint32(databs, uint32(dataSize))
	buffer.Write(databs)
	buffer.Write(cipherText)
	return buffer.Bytes(), nil
}

func (hyper *HyperEncrypt) Decrypt(cipherData []byte) ([]byte, error) {
	br := bytes.NewReader(cipherData)

	//read aes key size
	cipherAesKeySizeSlice := make([]byte, 4)
	if _, err := br.Read(cipherAesKeySizeSlice); err != nil {
		return nil, err
	}
	//read aes key
	cipherAesKeySize := binary.BigEndian.Uint32(cipherAesKeySizeSlice)
	cipherAesKey := make([]byte, cipherAesKeySize)
	if _, err := br.Read(cipherAesKey); err != nil {
		return nil, err
	}

	//read aes nonce size
	cipherAesNonceSizeSlice := make([]byte, 4)
	if _, err := br.Read(cipherAesNonceSizeSlice); err != nil {
		return nil, err
	}
	//read aes nonce
	cipherAesNonceSize := binary.BigEndian.Uint32(cipherAesNonceSizeSlice)
	cipherAesNonce := make([]byte, cipherAesNonceSize)
	if _, err := br.Read(cipherAesNonce); err != nil {
		return nil, err
	}

	//read aes nonce size
	cipherAesTextSizeSlice := make([]byte, 4)
	if _, err := br.Read(cipherAesTextSizeSlice); err != nil {
		return nil, err
	}
	//read aes nonce
	cipherAesTextSize := binary.BigEndian.Uint32(cipherAesTextSizeSlice)
	cipherAesText := make([]byte, cipherAesTextSize)
	if _, err := br.Read(cipherAesText); err != nil {
		return nil, err
	}

	aesKey, err := hyper.rsa.RsaDecrypt(cipherAesKey)
	if err != nil {
		return nil, err
	}

	aesNonce, err := hyper.rsa.RsaDecrypt(cipherAesNonce)
	if err != nil {
		return nil, err
	}

	hyper.aes.Key = aesKey
	hyper.aes.Nonce = aesNonce
	plainText, err := hyper.aes.GcmDecrypt(cipherAesText)
	if err != nil {
		return nil, err
	}
	return plainText, nil
}

/*func (hyper *HyperEncrypt) Decrypt2(cipherData []byte) ([]byte, error) {
	aesKeyLen := binary.BigEndian.Uint32(cipherData[0:4])
	aesKeyEndPos := 4 + aesKeyLen
	cipherAesKey := cipherData[4:aesKeyEndPos]

	aesNonceLen := binary.BigEndian.Uint32(cipherData[aesKeyEndPos : aesKeyEndPos+4])
	aesNonceStartPos := aesKeyEndPos + 4
	aesNonceEndPos := aesNonceStartPos + aesNonceLen
	cipherAesNonce := cipherData[aesNonceStartPos:aesNonceEndPos]

	cipherTextStartPos := aesNonceEndPos + 4
	cipherText := cipherData[cipherTextStartPos:]

	aesKey, err := hyper.rsa.RsaDecrypt(cipherAesKey)
	if err != nil {
		return nil, err
	}

	aesNonce, err := hyper.rsa.RsaDecrypt(cipherAesNonce)
	if err != nil {
		return nil, err
	}

	hyper.aes.Key = aesKey
	hyper.aes.Nonce = aesNonce
	plainText, err := hyper.aes.GcmDecrypt(cipherText)
	if err != nil {
		return nil, err
	}
	return plainText, nil
}*/

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
