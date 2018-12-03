package server

import "C"
import (
	"bytes"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
)

type ApiClient struct {
	Url string
}

func NewApiClient(url string) *ApiClient {
	client := &ApiClient{}
	client.Url = url
	return client
}

func (client *ApiClient) PostFile(filename, fieldname string) error {
	// 创建表单文件
	// CreateFormFile 用来创建表单，第一个参数是字段名，第二个参数是文件名
	buf := new(bytes.Buffer)
	writer := multipart.NewWriter(buf)
	formFile, err := writer.CreateFormFile(fieldname, filename)
	if err != nil {
		log.Fatalf("Create form file failed: %s\n", err)
		return err
	}

	// 从文件读取数据，写入表单
	srcFile, err := os.Open(filename)
	if err != nil {
		log.Fatalf("%Open source file failed: s\n", err)
		return err
	}
	defer srcFile.Close()
	_, err = io.Copy(formFile, srcFile)
	if err != nil {
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

func (client *ApiClient) GetFile(filename, fieldname string) error {
	return nil
}
