package api

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

func (a *Api) HandleDocker(w http.ResponseWriter, r *http.Request) {

	if a.dockerClient == nil {
		a.dockerClient = &http.Client{
			Transport: &http.Transport{
				DialContext: func(_ context.Context, _, _ string) (net.Conn, error) {
					return net.Dial("unix", a.dockerSocket)
				},
			},
			Timeout: time.Minute * 3,
		}
	}

	requestUri := r.RequestURI
	log.Debugf("%s://%s", a.dockerSocket, requestUri)
	resp, err := a.dockerClient.Get("http://unix" + requestUri)
	if err != nil {
		log.Error(err, ", ", a.host)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, err)
		a.dockerClient = nil
		return
	}
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
	resp.Body.Close()
}
