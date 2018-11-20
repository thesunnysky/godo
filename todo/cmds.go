package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
)

const dataFile = DATA_FILE

type dataEntity struct {
	id   int
	data string
}

var r, _ = regexp.Compile("[[:alnum:]]")

func addCmdImpl(args []string) {
	f, err := os.OpenFile(dataFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}

	var buf bytes.Buffer
	for _, str := range args {
		buf.WriteString(str)
		buf.WriteByte(' ')
	}
	buf.WriteByte('\n')

	if _, err = f.Write(buf.Bytes()); err != nil {
		panic(err)
	}
	fmt.Println("Task add successfully")
	defer f.Close()
}

func delCmdImpl(args []string) {
	num := make([]int, len(args))
	for _, str := range args {
		i, err := strconv.Atoi(str)
		if err != nil {
			fmt.Printf("Invalid Parameter Value:%s\n", str)
			os.Exit(INVALID_PARAMETER_VALUE)
		} else {
			num = append(num, i)
		}
	}

	f, err := os.OpenFile(dataFile, os.O_RDWR, 0666)
	if err != nil {
		panic(err)
	}

	br := bufio.NewReader(f)
	fileData := make([]string, 0)
	for {
		str, err := br.ReadString('\n')
		if err == io.EOF {
			break
		}
		fileData = append(fileData, str)
	}

	if err := f.Truncate(0); err != nil {
		panic(err)
	}

	for _, index := range num {
		idx := index - 1;
		if (idx < 0) || (idx > len(fileData)-1) {
			continue
		}
		fileData[idx] = string('\n')
	}

	var buf bytes.Buffer
	for _, str := range fileData {
		buf.WriteString(str)
	}
	if _, err := f.Write(buf.Bytes()); err != nil {
		panic(err)
	}
	defer f.Close()

	fmt.Println("Delete task successfully")
}

func listCmdImpl(args []string) {
	f, err := os.Open(dataFile)
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
		if !isEmptyLine(str) {
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

	/*
	 * read task
	 */
	br := bufio.NewReader(f)
	fileData := make([]string, 0)
	for {
		str, err := br.ReadString('\n')
		if err == io.EOF {
			break
		}
		if !isEmptyLine(str) {
			//remove empty line
			fileData = append(fileData, str)
		}
	}

	if err := f.Truncate(0); err != nil {
		panic(err)
	}

	//rewrite task file
	var buf bytes.Buffer
	for _, str := range fileData {
		buf.WriteString(str)
	}
	if _, err := f.Write(buf.Bytes()); err != nil {
		panic(err)
	}
	defer f.Close()

	fmt.Println("Tidy task file successfully")
}

func isEmptyLine(str string) bool {
	return !r.MatchString(str)
}
