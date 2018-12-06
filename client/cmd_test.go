package godo

import (
	"testing"
)

/*func TestDelCmds(t *testing.T) {
	args := []string{"2"}
	cmdImplMmp.DelCmdImpl(args)
}*/

func TestListCmds(t *testing.T) {
	args := []string{}
	ListCmdImpl(args)
}

func TestTidyCmds(t *testing.T) {
	args := []string{}
	TidyCmdImpl(args)
}

func TestTidyNormCmds(t *testing.T) {
	args := []string{}
	TidyCmdImpl(args)
}

func TestPushNormCmds(t *testing.T) {
	args := []string{}
	PushServerCmd(args)
}

func TestPullNormCmds(t *testing.T) {
	args := []string{}
	PullCmd(args)
}

func TestGitCmd(t *testing.T) {
	args := []string{"status"}
	_ = GitCmd(args)
}

/*func TestMmapReadFile(t *testing.T) {
	f, err := os.OpenFile(dataFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, consts.FILE_MAKS)
	defer f.Close()
	if err != nil {
		panic(err)
	}
}*/
