package mmpfile

import (
	"github.com/thesunnysky/godo/config"
	"os"
	"syscall"
)

type File struct {
	File *os.File
}

func (File) AppendNewLine(f *os.File, data []byte) {
	b := byte('\n')
	data = append(data, b)

	if _, err := f.Write(data); err != nil {
		panic(err)
	}
}

//todo there is bub!!!
func (File) RewriteFile(f *os.File, data []string) {
	if err := f.Truncate(0); err != nil {
		panic(err)
	}
	for _, line := range data {
		if _, err := f.WriteString(line); err != nil {
			panic(err)
		}
	}
}

func (File) ReadDataFile(f *os.File) (fileData []string) {
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

	fileData = make([]string, config.DEFAULT_LINE_CACHE)
	for begin, end := 0, 0; end < len(mMap); end++ {
		if mMap[end] == byte(config.LINE_SEPARATOR) {
			fileData = append(fileData, string(mMap[begin:end+1]))
			begin = end + 2
		}
	}
	return
}
