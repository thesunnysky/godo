package godo

import "C"
import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/thesunnysky/godo/consts"
	normalfile "github.com/thesunnysky/godo/util"
	"io"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
)

type File struct {
	File *os.File
}

var dataFile, privateKeyFile, publicKeyFile string

var dataFile2 = "/tmp/godo.dat2"

func init() {
	dataFile = ClientConfig.DataFile
	privateKeyFile = ClientConfig.PrivateKeyFile
	publicKeyFile = ClientConfig.PublicKeyFile
}

var r, _ = regexp.Compile("[[:alnum:]]")

func AddCmdImpl(args []string) {
	var buf bytes.Buffer
	for _, str := range args {
		buf.WriteString(str)
		buf.WriteByte(' ')
	}

	f, err := os.OpenFile(dataFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, consts.FILE_MAKS)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	file := normalfile.File{File: f}
	file.AppendNewLine(buf.Bytes())

	fmt.Println("task add successfully")
}

func DelCmdImpl(args []string) {
	num := make([]int, len(args))
	for _, str := range args {
		i, err := strconv.Atoi(str)
		if err != nil {
			fmt.Printf("invalid parameter value:%s\n", str)
			os.Exit(consts.INVALID_PARAMETER_VALUE)
		} else {
			num = append(num, i)
		}
	}

	f, err := os.OpenFile(dataFile, os.O_CREATE|os.O_RDWR, consts.FILE_MAKS)
	defer f.Close()
	file := normalfile.File{File: f}
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
	f, err := os.OpenFile(dataFile, os.O_CREATE|os.O_RDONLY, consts.FILE_MAKS)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	br := bufio.NewReader(f)
	var index int
	for {
		str, err := br.ReadString(consts.LINE_SEPARATOR)
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
	file := normalfile.File{File: f}
	defer file.Close()

	//read and filter task file
	br := bufio.NewReader(file.File)
	var fileData []string
	for {
		str, err := br.ReadString(consts.LINE_SEPARATOR)
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

func PushCmd(args []string) {
	data, err := ioutil.ReadFile(dataFile)
	if err != nil {
		panic(err)
	}
	encryptedData, err := RsaEncrypt(data)
	if err != nil {
		panic(err)
	}

	//todo create a new file and then rename to target file
	if err := ioutil.WriteFile(dataFile2, encryptedData, consts.FILE_MAKS);
		err != nil {
		panic(err)
	}

	fmt.Println("push task file successfully")
}
func PullCmd(args []string) {
	data, err := ioutil.ReadFile(dataFile2)
	if err != nil {
		panic(err)
	}
	encryptedData, err := RsaDecrypt(data)
	if err != nil {
		panic(err)
	}

	if err := ioutil.WriteFile(dataFile, encryptedData, consts.FILE_MAKS);
		err != nil {
		panic(err)
	}

	fmt.Println("pull task file successfully")
}
