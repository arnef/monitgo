package api

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os/exec"
	"strings"

	log "github.com/sirupsen/logrus"
)

func (a *Api) HandleExec(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	r.Body.Close()
	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, err)
		return
	}
	if len(body) == 0 {
		log.Error("empty command")
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "empty command")
		return
	}
	out, err := run(string(body), a.allowedCommands)
	if err != nil {
		log.Error(out, err)
		if err == ErrCommandNotAllowed {
			w.WriteHeader(http.StatusMethodNotAllowed)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		fmt.Fprintf(w, "%v\n%s", err, out)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, string(out))
}

var ErrCommandNotAllowed error = errors.New("command not allowed")

func run(command string, allowedCommands []string) ([]byte, error) {
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
