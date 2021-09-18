package api

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"

	"github.com/arnef/monitgo/interal/node-exporter/exec"
	log "github.com/sirupsen/logrus"
)

func Start(host string, port int, dockerSocket string, allowedCommands []string) error {
	httpc := http.Client{
		Transport: &http.Transport{
			DialContext: func(_ context.Context, _, _ string) (net.Conn, error) {
				return net.Dial("unix", dockerSocket)
			},
		},
	}
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost && strings.HasPrefix(r.RequestURI, "/xc") {
			body, err := io.ReadAll(r.Body)
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
			out, err := exec.Run(string(body), allowedCommands)
			if err != nil {
				log.Error(out, err)
				if err == exec.ErrCommandNotAllowed {
					w.WriteHeader(http.StatusMethodNotAllowed)
				} else {
					w.WriteHeader(http.StatusInternalServerError)
				}
				fmt.Fprintf(w, "%v\n%s", err, out)
				return
			}
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, string(out))
			return
		}

		if r.Method == http.MethodGet && r.RequestURI != "/xc" {
			requestUri := r.RequestURI
			log.Debugf("%s://%s", dockerSocket, requestUri)
			resp, err := httpc.Get("http://unix" + requestUri)
			if err != nil {
				log.Error(err)
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprint(w, err)
				return
			}
			w.WriteHeader(resp.StatusCode)
			defer resp.Body.Close()
			io.Copy(w, resp.Body)
			return
		}
		w.WriteHeader(http.StatusNotFound)
	})
	log.Infof("running node on %s:%d", host, port)
	return http.ListenAndServe(fmt.Sprintf("%s:%d", host, port), nil)
}
