package godo

import (
	"fmt"
	"github.com/thesunnysky/godo/config"
	"os"
	"testing"
)

func TestDelCmds(t *testing.T) {
	args := []string{"4"}
	delCmdImpl(args)
}

func TestListCmds(t *testing.T) {
	args := []string{}
	listCmdImpl(args)
}

func TestMmapReadFile(t *testing.T) {
	f, err := os.OpenFile(dataFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, config.FILE_MAKS)
	defer f.Close()
	if err != nil {
		panic(err)
	}
}

func TestSlice(t *testing.T) {
	slice := make([]int, 10)
	for i := 0; i < 10; i++ {
		slice[i] = i
	}

	b := slice[0:2]
	fmt.Println(b)
}
