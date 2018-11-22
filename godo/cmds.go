package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"regexp"
	"runtime"
	"strconv"
	"strings"
)

var dataFile string
var lineSeparator byte

func init() {
	initLineSeparator()

	initDataFile()
}

func initLineSeparator() {
	lineSeparator = byte(UNIX_LINE_SEPARATOR)
}

func initDataFile() {
	homeDir := os.Getenv("HOME")
	configFile := homeDir + "/" + CONFIG_FILE
	if !pathExist(configFile) {
		fmt.Printf("config file:$HOME/%s do not exist\n", CONFIG_FILE)
		os.Exit(CONFIG_FILE_DO_NOT_EXIST)
	}
	f, err := os.Open(configFile)
	if err != nil {
		panic(nil)
	}
	defer f.Close()

	dataFile, err = bufio.NewReader(f).ReadString(lineSeparator)
	if err != nil {
		panic(err)
	}

	dataFile = strings.TrimSpace(dataFile)
}

var r, _ = regexp.Compile("[[:alnum:]]")

func addCmdImpl(args []string) {
	var buf bytes.Buffer
	for _, str := range args {
		buf.WriteString(str)
		buf.WriteByte(' ')
	}

	f, err := os.OpenFile(dataFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, FILE_MAKS);
	if err != nil {
		panic(err)
	}

	appendNewLine(f, buf.Bytes())
	defer f.Close()

	fmt.Println("task add successfully")
}

func delCmdImpl(args []string) {
	num := make([]int, len(args))
	for _, str := range args {
		i, err := strconv.Atoi(str)
		if err != nil {
			fmt.Printf("invalid parameter value:%s\n", str)
			os.Exit(INVALID_PARAMETER_VALUE)
		} else {
			num = append(num, i)
		}
	}

	f, err := os.OpenFile(dataFile, os.O_CREATE|os.O_RDWR, FILE_MAKS)
	if err != nil {
		panic(err)
	}

	fileData := readDataFile(f)
	for _, index := range num {
		idx := index - 1;
		if (idx < 0) || (idx > len(fileData)-1) {
			continue
		}
		fileData[idx] = string('\n')
	}

	rewriteFile(f, fileData)
	defer f.Close()

	fmt.Println("delete task successfully")
}

func listCmdImpl(args []string) {
	f, err := os.OpenFile(dataFile, os.O_CREATE|os.O_RDONLY, FILE_MAKS)
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

func cleanCmdImpl(args []string) {
	f, err := os.OpenFile(dataFile, os.O_RDWR, 0666)
	if err != nil {
		panic(err)
	}

	//read and filter task file
	br := bufio.NewReader(f)
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

	//rewrite task file
	rewriteFile(f, fileData)
	defer f.Close()

	fmt.Println("tidy task file successfully")
}

func isBlankLine(str string) bool {
	return !r.MatchString(str)
}

func readDataFile(f *os.File) []string {
	br := bufio.NewReader(f)
	fileData := make([]string, 0)
	for {
		str, err := br.ReadString(lineSeparator)
		if err == io.EOF {
			break
		}
		fileData = append(fileData, str)
	}
	return fileData
}

func appendNewLine(f *os.File, data []byte) {
	b := byte('\n')
	data = append(data, b)

	if _, err := f.Write(data); err != nil {
		panic(err)
	}
}

func rewriteFile(f *os.File, data []string) {
	if err := f.Truncate(0); err != nil {
		panic(err)
	}
	for _, line := range data {
		if _, err := f.WriteString(line); err != nil {
			panic(err)
		}
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
