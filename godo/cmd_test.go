package main

import "testing"

func TestDelCmds(t *testing.T) {
	args := []string{"4"}
	delCmdImpl(args)
}

func TestListCmds(t *testing.T) {
	args := []string{}
	listCmdImpl(args)
}
