package exec

import (
	"bytes"
	"errors"
	"fmt"
	"os/exec"
	"strings"

	log "github.com/sirupsen/logrus"
)

var ErrCommandNotAllowed error = errors.New("command not allowed")

func Run(command string, allowedCommands []string) ([]byte, error) {
	log.Debugln(command)
	in := strings.Split(command, " ")
	if len(in) == 0 {
		return nil, fmt.Errorf("no command given")
	}

	for i := range allowedCommands {
		if in[0] == allowedCommands[i] {
			var outb, errb bytes.Buffer
			cmd := exec.Command(in[0], in[1:]...)
			cmd.Stdout = &outb
			cmd.Stderr = &errb

			err := cmd.Run()

			if err != nil && outb.Len() == 0 {
				return nil, fmt.Errorf(errb.String())
			}

			return outb.Bytes(), nil
		}
	}
	return nil, ErrCommandNotAllowed
}
