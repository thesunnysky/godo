package util

import "testing"

func TestBackupFile(t *testing.T) {
	filepath := "/home/sun/data/godo/godo.dat"
	if err := BackupFile(filepath); err != nil {
		panic(err)
	}
}
