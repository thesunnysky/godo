package main

import (
	"fmt"
	"runtime"
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

func TestGOOS(t *testing.T) {
	fmt.Println(runtime.GOOS)
}
