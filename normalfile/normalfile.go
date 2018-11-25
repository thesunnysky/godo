package normalfile

import (
	"bufio"
	"fmt"
	"github.com/thesunnysky/godo/config"
	"io"
	"os"
)

type File struct {
	File *os.File
}

func (File) ReadDataFile(f *os.File) []string {
	br := bufio.NewReader(f)
	fileData := make([]string, 0)
	for {
		str, err := br.ReadString(config.LINE_SEPARATOR)
		if err == io.EOF {
			break
		}
		fileData = append(fileData, str)
	}
	return fileData
}

func (File) AppendNewLine(f *os.File, data []byte) {
	b := byte('\n')
	data = append(data, b)

	if _, err := f.Write(data); err != nil {
		panic(err)
	}
}

func (File) RewriteFile(f *os.File, data []string) {
	fmt.Println(data)
	if err := f.Truncate(0); err != nil {
		panic(err)
	}
	f.Sync()
	for _, line := range data {
		fmt.Println(line)
		if _, err := f.WriteString(line); err != nil {
			panic(err)
		}
	}
}
