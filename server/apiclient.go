package server

import "C"
import (
	"bytes"
	"encoding/base64"
	"github.com/thesunnysky/godo/util"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"strings"
)

type ApiClient struct {
	Url string
}

func NewApiClient(url string) *ApiClient {
	client := &ApiClient{}
	client.Url = url
	return client
}

func (client *ApiClient) PostFile(fieldname, filename string) error {
	index := strings.LastIndex(filename, "/")
	if index == -1 {
		index = 0
	}
	fileName := filename[index:]

	// 创建表单文件
	// CreateFormFile 用来创建表单，第一个参数是字段名，第二个参数是文件名
	buf := new(bytes.Buffer)
	writer := multipart.NewWriter(buf)
	formFile, err := writer.CreateFormFile(fieldname, fileName)

	if err != nil {
		log.Fatalf("Create form file failed: %s\n", err)
		return err
	}

	// 从文件读取数据，写入表单
	srcFileData, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatalf("Open source file failed: %s\n", err)
		return err
	}

	base64.NewEncoder()
	encryptData, err := util.RsaEncrypt(srcFileData)
	if err != nil {
		log.Fatalf("Encrypt data error:%s\n", err)
		return err
	}
	if _, err := formFile.Write(encryptData); err != nil {
		log.Fatalf("Write to form file falied: %s\n", err)
		return err
	}

	// 发送表单
	contentType := writer.FormDataContentType()
	// 发送之前必须调用Close()以写入结尾行
	if err := writer.Close(); err != nil {
		log.Fatalf("Write to form file falied: %s\n", err)
		return err
	}

	_, err = http.Post(client.Url+"/upload", contentType, buf)
	if err != nil {
		log.Fatalf("Post failed: %s\n", err)
		return err
	}
	return nil
}

func (client *ApiClient) DownloadFile(filename string) (io.Reader, error) {
	r, err := http.Get(client.Url + "/download/" + filename)
	if err != nil {
		log.Printf("failed to download file\n:%s", filename)
		return nil, err
	}
	return r.Body, nil
}
