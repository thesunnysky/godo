package util

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
	"time"
)

type Git struct {
	Repo   string
	Author string
}

func (g *Git) AddAll() error {
	args := []string{"add -A"}
	stdout, stderr, err := g.exec(args...)
	return checkResult(stdout, stderr, err)
}

func (g *Git) Commit() error {
	args := []string{"commit -m \"\""}
	stdout, stderr, err := g.exec(args...)
	return checkResult(stdout, stderr, err)
}

func (g *Git) Push() error {
	args := []string{"push"}
	stdout, stderr, err := g.exec(args...)
	return checkResult(stdout, stderr, err)
}

func checkResult(stdout, stderr *bytes.Buffer, err error) error {
	if err != nil {
		return err
	}
	if stderr.Len() > 0 {
		log.Print(stderr.Bytes())
	}
	if stdout.Len() > 0 {
		log.Print(stdout.Bytes())
	}
	return nil

}

func (g *Git) exec(args ...string) (stdout, stderr *bytes.Buffer, err error) {
	cmd := exec.Command("git", args...)

	stdout = new(bytes.Buffer)
	stderr = new(bytes.Buffer)

	cmd.Stdout = stdout
	cmd.Stderr = stderr
	cmd.Dir = g.Repo

	if err := cmd.Start(); err != nil {
		return
	}

	done := make(chan error)
	go func() {
		done <- cmd.Wait()
	}()

	timeout := 30 * time.Second
	select {
	case <-time.After(timeout):
		if cmd.Process != nil && cmd.ProcessState != nil && !cmd.ProcessState.Exited() {
			if err := cmd.Process.Kill(); err != nil {
				return
			}
		}
		<-done
		err = fmt.Errorf("command execute timeout")
		return

	case err = <-done:
		return
	}
}
