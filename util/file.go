package util

import (
	"bufio"
	"github.com/thesunnysky/godo/consts"
	"io"
	"os"
)

type File struct {
	File *os.File
}

func (f File) ReadFile() []string {
	br := bufio.NewReader(f.File)
	var fileData []string
	for {
		str, err := br.ReadString(consts.LINE_SEPARATOR)
		if err != nil || err == io.EOF {
			break
		}
		fileData = append(fileData, str)
	}
	return fileData
}

func (f File) AppendNewLine(data []byte) {
	b := byte('\n')
	data = append(data, b)

	if _, err := f.File.Write(data); err != nil {
		panic(err)
	}
}

func (f File) RewriteFile(data []string) {
	if err := f.File.Truncate(0); err != nil {
		panic(err)
	}
	_, err := f.File.Seek(0, 0)
	if err != nil {
		panic(err)
	}
	for _, line := range data {
		if _, err := f.File.WriteString(line); err != nil {
			panic(err)
		}
	}
}

func (f File) Close() {
	f.File.Close()
}

func PathExist(path string) bool {
	_, err := os.Stat(path)
	if err != nil && os.IsNotExist(err) {
		return false
	}
	return true
}

func CreateDirIsNotExist(path string) error {
	if PathExist(path) {
		return nil
	}
	return os.MkdirAll(path, 0711)
}
