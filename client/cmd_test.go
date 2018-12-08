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
	ListLocalTasks(args)
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
	PullServerCmd(args)
}

func TestGitCmd(t *testing.T) {
}

func TestListBackupFiles(t *testing.T) {
	args := []string{}
	ListBackupFiles(args)
}

func TestListBackupTasks(t *testing.T) {
	ListBackupTasks(1)
}

/*func TestMmapReadFile(t *testing.T) {
	f, err := os.OpenFile(dataFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, consts.FILE_MAKS)
	defer f.Close()
	if err != nil {
		panic(err)
	}
}*/

func TestDelBackupTaskFile(t *testing.T) {
	type args struct {
		args []string
	}
	tests := []struct {
		name string
		args args
	}{
		{"hello", args{args: []string{"all"}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			DelBackupTaskFile(tt.args.args)
		})
	}
}
