package cmd

import (
	"github.com/thesunnysky/godo/cmd/cmdImplMmp"
	"testing"
)

func TestDelCmds(t *testing.T) {
	args := []string{"2"}
	cmdImplMmp.DelCmdImpl(args)
}

func TestListCmds(t *testing.T) {
	args := []string{}
	cmdImplMmp.ListCmdImpl(args)
}

func TestTidyCmds(t *testing.T) {
	args := []string{}
	cmdImplMmp.CleanCmdImpl(args)
}

/*func TestMmapReadFile(t *testing.T) {
	f, err := os.OpenFile(dataFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, config.FILE_MAKS)
	defer f.Close()
	if err != nil {
		panic(err)
	}
}*/
