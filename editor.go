package main

import (
	"os"
	"os/exec"
)

func (u *User) OpenEditor(filename string) error {
	if u.TerminalEditor {
		return RunCommand(u.Editor, filename)
	} else {
		return StartCommand(u.Editor, filename)
	}
}

func StartCommand(name string, arg ...string) error {
	cmd := exec.Command(name, arg...)
	return cmd.Start()
}

func RunCommand(name string, arg ...string) error {
	cmd := exec.Command(name, arg...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
