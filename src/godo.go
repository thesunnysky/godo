package godo

import "syscall"

func main() {
	cmd := "ls /home"
	syscall.Exec(cmd)

}
