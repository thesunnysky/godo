package cmdImplMmp

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/thesunnysky/godo/config"
	"github.com/thesunnysky/godo/normalfile"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"
	"syscall"
)

type File struct {
	File *os.File
}

var dataFile string

func init() {
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
	if err != nil {
		panic(err)
	}

	fileInfo, err := f.Stat()
	if err != nil {
		os.Exit(1)
	}
	size := int(fileInfo.Size())

	mMap, err := syscall.Mmap(int(f.Fd()), 0, size, syscall.PROT_READ|syscall.PROT_WRITE, syscall.MAP_SHARED)
	defer syscall.Munmap(mMap)
	if err != nil {
		os.Exit(1)
	}

	lineNum := 0
	for begin, end := 0, 0; end < len(mMap); end++ {
		if mMap[end] == byte(config.LINE_SEPARATOR) {
			lineNum++
			if intContains(num, lineNum) {
				//write null to deleted line
				writeNull(mMap, begin, end)
			}
			begin = end + 1
		}
	}

	fmt.Println("delete task successfully")
}

func ListCmdImpl(args []string) {
	f, err := os.OpenFile(dataFile, os.O_CREATE|os.O_RDONLY, config.FILE_MAKS)
	defer f.Close()
	if err != nil {
		panic(err)
	}

	fileInfo, err := f.Stat()
	if err != nil {
		os.Exit(1)
	}
	size := int(fileInfo.Size())

	mMap, err := syscall.Mmap(int(f.Fd()), 0, size, syscall.PROT_READ, syscall.MAP_SHARED)
	defer syscall.Munmap(mMap)
	if err != nil {
		os.Exit(1)
	}

	var index int
	for begin, end := 0, 0; end < len(mMap); end++ {
		if mMap[end] == byte(config.LINE_SEPARATOR) {
			index++
			str := string(mMap[begin:end])
			if !isBlankLine(str) {
				fmt.Printf("%d. %s\n", index, str)
			}
			begin = end + 1
		}
	}
}

func CleanCmdImpl(args []string) {
	f, err := os.OpenFile(dataFile, os.O_RDWR, 0666)
	defer f.Close()
	if err != nil {
		panic(err)
	}
	file := normalfile.File{File: f}

	//read and filter task myfile
	br := bufio.NewReader(file.File)
	fileData := make([]string, 0)
	for {
		str, err := br.ReadString(config.LINE_SEPARATOR)
		fmt.Printf("str:%s, len:%d", str, len(str))
		if err == io.EOF {
			break
		}
		if !isBlankLine(str) {
			//remove empty line
			fileData = append(fileData, str)
		}
	}

	//rewrite task file
	file.RewriteFile(fileData)
	_ = file.File.Sync()

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

func writeNull(data []byte, begin, end int) {
	for i := begin; i < end; i++ {
		data[i] = byte(0)
	}
}

func intContains(intArray []int, target int) bool {
	for _, value := range intArray {
		if value == target {
			return true
		}
	}
	return false

}
