package api

import (
	"fmt"
	"net/http"
	"strings"

	log "github.com/sirupsen/logrus"
)

type Api struct {
	dockerSocket    string
	host            string
	port            int
	allowedCommands []string
}

func New(host string, port int, dockerSocket string, allowedCommands []string) *Api {

	return &Api{
		dockerSocket:    dockerSocket,
		host:            host,
		port:            port,
		allowedCommands: allowedCommands,
	}
}

func (a *Api) Start() error {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost && strings.HasPrefix(r.RequestURI, "/xc") {
			a.HandleExec(w, r)
		} else if r.Method == http.MethodGet && r.RequestURI != "/xc" {
			a.HandleDocker(w, r)
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	})
	log.Infof("running node on %s:%d", a.host, a.port)
	return http.ListenAndServe(fmt.Sprintf("%s:%d", a.host, a.port), nil)
}
