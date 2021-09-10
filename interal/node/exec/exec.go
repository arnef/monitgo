package exec

import (
	"errors"
	"fmt"
	"os/exec"
	"strings"

	log "github.com/sirupsen/logrus"
)

var ErrCommandNotAllowed error = errors.New("command not allowed")

func Run(command string, allowedCommands []string) (string, error) {
	log.Debugln(command)
	in := strings.Split(command, " ")
	if len(in) == 0 {
		return "", fmt.Errorf("no command given")
	}

	for i := range allowedCommands {
		if in[0] == allowedCommands[i] {
			cmd := exec.Command(in[0], in[1:]...)
			out, err := cmd.Output()

			return string(out), err
		}
	}
	return "", ErrCommandNotAllowed
}
