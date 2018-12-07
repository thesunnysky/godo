package util

import (
	"bufio"
	"fmt"
	"github.com/CodisLabs/codis/pkg/utils/errors"
	"github.com/thesunnysky/godo/consts"
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
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

func ExtractFileName(filepath string) string {
	index := strings.LastIndex(filepath, "/")
	if index == -1 {
		index = 0
	}
	return filepath[index+1:]
}

func ExtractFileDir(filepath string) (string, error) {
	index := strings.LastIndex(filepath, "/")
	if index == -1 {
		return "", errors.New("file path is not start with root(/)")
	} else if index == 0 {
		return "/", nil
	} else {
		return filepath[:index], nil
	}
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

func RewriteFile(path string, data []byte) error {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, consts.FILE_MAKS)
	if err != nil {
		fmt.Printf("open file %s error\n", path)
		return err
	}

	if err := f.Truncate(0); err != nil {
		fmt.Printf("truncate file %s error\n", path)
		return err
	}

	if _, err := f.Seek(0, 0); err != nil {
		fmt.Printf("seek file %s error\n", path)
		return err
	}

	if _, err := f.Write(data); err != nil {
		fmt.Printf("write file %s error\n", path)
		return err
	}
	return nil
}

func BackupFile(filepath string) error {
	filename := ExtractFileName(filepath)
	fileDir, err := ExtractFileDir(filepath)
	if err != nil {
		return err
	}

	fileInfos, err := ioutil.ReadDir(fileDir)
	if err != nil {
		return err
	}

	var maxSuffix int
	for _, fileInfo := range fileInfos {
		backupFileName := fileInfo.Name()
		if !fileInfo.Mode().IsRegular() || !strings.HasPrefix(backupFileName, filename) {
			continue
		}
		suffixPos := strings.LastIndex(backupFileName, ".")
		suffixI, err := strconv.Atoi(backupFileName[suffixPos+1:])
		if err == nil && suffixI > maxSuffix {
			maxSuffix = suffixI
		}
	}

	backFilepath := fileDir + "/" + filename + "." + strconv.Itoa(int(maxSuffix+1))
	if err := CopyFile(backFilepath, filepath); err != nil {
		return err
	}
	return nil
}

func CopyFile(dst, src string) error {
	from, err := os.Open(src)
	if err != nil {
		return err
	}
	defer from.Close()

	to, err := os.OpenFile(dst, os.O_CREATE|os.O_WRONLY, consts.FILE_MAKS)
	if err != nil {
		return err
	}
	defer to.Close()

	_, err = io.Copy(to, from)
	if err != nil {
		return err
	}
	return nil
}
