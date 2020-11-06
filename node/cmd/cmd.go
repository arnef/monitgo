package cmd

import (
	"bytes"
	"fmt"
	"os/exec"
)

func Exec(cmd string, args ...string) ([]byte, error) {
	path, err := exec.LookPath(cmd)
	if err != nil {
		return nil, err
	}

	var outb, errb bytes.Buffer
	exe := exec.Cmd{
		Path:   path,
		Stderr: &errb,
		Stdout: &outb,
		Args:   append([]string{path}, args...),
	}
	err = exe.Run()
	if err != nil {
		return nil, fmt.Errorf(errb.String())
	}

	return outb.Bytes(), nil
}
