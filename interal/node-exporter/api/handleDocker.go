package api

import (
	"context"
	"fmt"
	"io/ioutil"
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
	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.WriteHeader(resp.StatusCode)
		w.Write(body)
	}
}
