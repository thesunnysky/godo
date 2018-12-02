package cmd

import (
	"github.com/thesunnysky/godo/cmd/cmdImplMmp"
	"github.com/thesunnysky/godo/cmd/cmdImplNorm"
	"github.com/thesunnysky/godo/normalfile"
	"os"
	"testing"
)

/*func TestDelCmds(t *testing.T) {
	args := []string{"2"}
	cmdImplMmp.DelCmdImpl(args)
}*/

func TestListCmds(t *testing.T) {
	args := []string{}
	cmdImplMmp.ListCmdImpl(args)
}

func TestTidyCmds(t *testing.T) {
	args := []string{}
	cmdImplMmp.CleanCmdImpl(args)
}

func TestRewriteFile(t *testing.T) {
	path := "/tmp/aaa.txt"
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0666)

	if err != nil {
		panic(err)
	}

	file := normalfile.File{File: f}
	defer file.File.Close()

	data := file.ReadFile()

	file.RewriteFile(data)
}

func TestTidyNormCmds(t *testing.T) {
	args := []string{}
	cmdImplNorm.TidyCmdImpl(args)
}

func TestPushNormCmds(t *testing.T) {
	args := []string{}
	cmdImplNorm.PushCmd(args)
}

func TestPullNormCmds(t *testing.T) {
	args := []string{}
	cmdImplNorm.PullCmd(args)
}

/*func TestMmapReadFile(t *testing.T) {
	f, err := os.OpenFile(dataFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, config.FILE_MAKS)
	defer f.Close()
	if err != nil {
		panic(err)
	}
}*/
