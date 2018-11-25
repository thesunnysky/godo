package cmdImplNorm

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/thesunnysky/godo/config"
	"github.com/thesunnysky/godo/mmpfile"
	"github.com/thesunnysky/godo/normalfile"
	"io"
	"os"
	"regexp"
	"runtime"
	"strconv"
	"strings"
)

type File struct {
	File *os.File
}

var dataFile string
var lineSeparator byte

func init() {
	initLineSeparator()

	initDataFile()
}

func initLineSeparator() {
	lineSeparator = byte(config.LINE_SEPARATOR)
}

func initDataFile() {
	homeDir := os.Getenv("HOME")
	configFile := homeDir + "/" + config.CONFIG_FILE
	if !pathExist(configFile) {
		fmt.Printf("config myfile:$HOME/%s do not exist\n", config.CONFIG_FILE)
		os.Exit(config.CONFIG_FILE_DO_NOT_EXIST)
	}
	f, err := os.Open(configFile)
	if err != nil {
		panic(nil)
	}
	defer f.Close()

	dataFile, err = bufio.NewReader(f).ReadString(config.LINE_SEPARATOR)
	if err != nil {
		panic(err)
	}

	dataFile = strings.TrimSpace(dataFile)
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

	appendNewLine(f, buf.Bytes())
	defer f.Close()

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
	file := mmpfile.File{File: f}
	if err != nil {
		panic(err)
	}

	fileData := file.ReadDataFile(f)
	for _, index := range num {
		idx := index - 1
		if (idx < 0) || (idx > len(fileData)-1) {
			continue
		}
		fileData[idx] = string('\n')
	}

	file.RewriteFile(f, fileData)

	fmt.Println("delete task successfully")
}

func ListCmdImpl(args []string) {
	f, err := os.OpenFile(dataFile, os.O_CREATE|os.O_RDONLY, config.FILE_MAKS)
	if err != nil {
		panic(err)
	}
	br := bufio.NewReader(f)
	var index int
	for {
		str, err := br.ReadString('\n')
		if err == io.EOF {
			break
		}
		index++
		if !isBlankLine(str) {
			fmt.Printf("%d. %s", index, str)
		}
	}
	defer f.Close()
}

func CleanCmdImpl(args []string) {
	f, err := os.OpenFile(dataFile, os.O_RDWR, 0666)
	if err != nil {
		panic(err)
	}
	file := normalfile.File{File: f}

	//read and filter task myfile
	br := bufio.NewReader(file.File)
	fileData := make([]string, 0)
	for {
		str, err := br.ReadString('\n')
		if err == io.EOF {
			break
		}
		if !isBlankLine(str) {
			//remove empty line
			fileData = append(fileData, str)
		}
	}

	//rewrite task myfile
	file.RewriteFile(f, fileData)
	defer f.Close()

	fmt.Println("tidy task myfile successfully")
}

func isBlankLine(str string) bool {
	return !r.MatchString(str)
}

func appendNewLine(f *os.File, data []byte) {
	b := byte('\n')
	data = append(data, b)

	if _, err := f.Write(data); err != nil {
		panic(err)
	}
}

func pathExist(path string) bool {
	_, err := os.Stat(path)
	if err != nil && os.IsNotExist(err) {
		return false
	}
	return true
}

func osType() string {
	return runtime.GOOS
}
