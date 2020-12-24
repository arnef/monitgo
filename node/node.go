package node

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"git.arnef.de/monitgo/node/docker"
	"git.arnef.de/monitgo/node/host"
	"github.com/urfave/cli/v2"

	"git.arnef.de/monitgo/log"
)

type JsonStats struct {
	Container map[string]docker.Stats
	Host      host.Stats
	Error     *string
}

// Cmd start node exporter
func Cmd(ctx *cli.Context) error {
	port := ctx.Uint("port")
	host := ctx.String("host")
	if ctx.Bool("dry-run") {
		dryRun()
		return nil
	}
	http.HandleFunc("/stats", stats)
	log.Infof("🚀️ running at %s:%d\n", host, port)
	return http.ListenAndServe(fmt.Sprintf("%s:%d", host, port), nil)

}

func dryRun() {
	writeStats(os.Stdout, true)
}

func stats(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	writeStats(w, false)
}

func writeStats(w io.Writer, pretty bool) {

	container, containerError := docker.GetStats()
	host, hostError := host.GetStats()

	encoder := json.NewEncoder(w)
	if pretty {
		encoder.SetIndent("", "  ")
	}
	if containerError != nil {
		encoder.Encode(map[string]string{
			"Error": containerError.Error(),
		})
	} else if hostError != nil {
		encoder.Encode(map[string]string{
			"Error": hostError.Error(),
		})
	} else {
		encoder.Encode(map[string]interface{}{
			"Container": container,
			"Host":      host,
		})
	}
}
