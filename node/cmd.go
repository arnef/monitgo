package node

import (
	"encoding/json"
	"fmt"
	"net/http"

	"git.arnef.de/monitgo/docker"
	"github.com/urfave/cli/v2"
)

// Cmd start node exporter
func Cmd(ctx *cli.Context) error {
	port := ctx.Uint("port")
	host := ctx.String("host")
	http.HandleFunc("/stats", stats)
	fmt.Printf("🚀️ running at %s:%d\n", host, port)
	return http.ListenAndServe(fmt.Sprintf("%s:%d", host, port), nil)
}

func stats(w http.ResponseWriter, r *http.Request) {
	stats, err := docker.GetStats()

	w.Header().Set("content-type", "application/json")
	if err != nil {
		json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
	} else {
		json.NewEncoder(w).Encode(map[string][]docker.Stats{"data": stats})
	}
}
