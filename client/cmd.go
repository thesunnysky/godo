package godo

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/thesunnysky/godo/config"
	"github.com/thesunnysky/godo/server"
	"github.com/thesunnysky/godo/util"
	"io"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type File struct {
	File *os.File
}

var dataFile, privateKeyFile, publicKeyFile string

var dataFile2 = "/tmp/godo.dat2"

func init() {
	dataFile = config.ClientConfig.DataFile
	privateKeyFile = config.ClientConfig.PrivateKeyFile
	publicKeyFile = config.ClientConfig.PublicKeyFile
}

var r, _ = regexp.Compile("[[:alnum:]]")

func AddCmdImpl(args []string) {
	var buf bytes.Buffer
	for _, str := range args {
		buf.WriteString(str)
		buf.WriteByte(' ')
	}

	f, err := os.OpenFile(dataFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, config.FILE_MAKS)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	file := util.File{File: f}
	file.AppendNewLine(buf.Bytes())

	fmt.Println("task add successfully")
}

func DelCmdImpl(args []string) {
	num := make([]int, len(args))
	for _, str := range args {
		i, err := strconv.Atoi(str)
		if err != nil {
			fmt.Printf("invalid parameter value:%s\n", str)
			os.Exit(config.INVALID_PARAMETER_VALUE)
		} else {
			num = append(num, i)
		}
	}

	f, err := os.OpenFile(dataFile, os.O_CREATE|os.O_RDWR, config.FILE_MAKS)
	defer f.Close()
	file := util.File{File: f}
	if err != nil {
		panic(err)
	}

	fileData := file.ReadFile()
	for _, index := range num {
		idx := index - 1
		if (idx < 0) || (idx > len(fileData)-1) {
			continue
		}
		fileData[idx] = string('\n')
	}

	file.RewriteFile(fileData)

	fmt.Println("delete task successfully")
}

func ListCmdImpl(args []string) {
	f, err := os.OpenFile(dataFile, os.O_CREATE|os.O_RDONLY, config.FILE_MAKS)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	br := bufio.NewReader(f)
	var index int
	for {
		str, err := br.ReadString(config.LINE_SEPARATOR)
		if err == io.EOF {
			break
		}
		index++
		if !isBlankLine(str) {
			fmt.Printf("%d. %s", index, str)
		}
	}
}

func TidyCmdImpl(args []string) {
	f, err := os.OpenFile(dataFile, os.O_RDWR, 0666)
	if err != nil {
		panic(err)
	}
	file := util.File{File: f}
	defer file.Close()

	//read and filter task file
	br := bufio.NewReader(file.File)
	var fileData []string
	for {
		str, err := br.ReadString(config.LINE_SEPARATOR)
		if err == io.EOF {
			break
		}
		if !isBlankLine(str) {
			//remove empty line
			fileData = append(fileData, str)
		}
	}

	//rewrite task myfile
	file.RewriteFile(fileData)

	fmt.Println("tidy task file successfully")
}

func isBlankLine(str string) bool {
	return !r.MatchString(str)
}

func PullCmd(args []string) {
	filename := config.ClientConfig.DataFile
	index := strings.LastIndex(filename, "/")
	if index == -1 {
		index = 0
	}
	fileName := filename[index:]

	apiClient := server.ApiClient{Url: config.ClientConfig.GodoServerUrl}
	reader, err := apiClient.DownloadFile(fileName)
	if err != nil {
		log.Printf("download file:%s successfullly\n", filename)
		os.Exit(-1)
	}
	data, err := ioutil.ReadAll(reader)
	if err != nil {
		log.Printf("read data from response error:%s\n", err)
	}
	decryptedData, err := util.RsaDecrypt(data)
	if err != nil {
		log.Printf("decrypt data error:%s\n", err)
	}

	//write download data to tempFile
	tempFile := "temp" + config.ClientConfig.DataFile
	targetFile := config.ClientConfig.DataFile
	f, err := os.OpenFile(tempFile, os.O_WRONLY|os.O_TRUNC, config.FILE_MAKS)
	if err != nil {
		log.Printf("open file error:%s\n", err)
	}
	defer f.Close()
	if _, err := f.Write(decryptedData); err != nil {
		log.Printf("open file error:%s\n", err)
		os.Exit(-1)
	}

	//remove current task file
	if err := os.Remove(targetFile); err != nil {
		log.Printf("remove old task file error%s\n", err)
		os.Exit(-1)
	}
	// rename tempFile to targetFile
	if err := os.Rename(tempFile, targetFile); err != nil {
		log.Printf("remove old task file error%s\n", err)
		os.Exit(-1)
	}
	fmt.Println("pull task file successfully")
}

func PushCmd(args []string) {

	apiClient := server.ApiClient{Url: config.ClientConfig.GodoServerUrl}
	if err := apiClient.PostFile(config.GODO_DATA_FILE, config.ClientConfig.DataFile);
		err != nil {
		fmt.Printf("push task file to server error:%s\n", err)
		os.Exit(-1)
	}

	fmt.Println("pull task file successfully")
}
