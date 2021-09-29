package api

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"

	log "github.com/sirupsen/logrus"
)

func (a *Api) HandleDocker(w http.ResponseWriter, r *http.Request) {
	httpc := http.Client{
		Transport: &http.Transport{
			DialContext: func(_ context.Context, _, _ string) (net.Conn, error) {
				return net.Dial("unix", a.dockerSocket)
			},
		},
	}
	defer httpc.CloseIdleConnections()
	requestUri := r.RequestURI
	log.Debugf("%s://%s", a.dockerSocket, requestUri)
	resp, err := httpc.Get("http://unix" + requestUri)
	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, err)
		return
	}
	defer resp.Body.Close()
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}
