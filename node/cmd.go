package node

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"git.arnef.de/monitgo/docker"
	"github.com/urfave/cli/v2"
)

// Cmd start node exporter
func Cmd(ctx *cli.Context) error {
	port := ctx.Uint("port")
	host := ctx.String("host")
	if ctx.Bool("dry-run") {
		dryRun()
		return nil
	}
	http.HandleFunc("/stats", stats)
	fmt.Printf("🚀️ running at %s:%d\n", host, port)
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
	start := time.Now()
	fmt.Print("⏳ get stats ")
	stats, err := docker.GetStats()
	duration := time.Since(start)
	fmt.Printf("took %s\n", duration)

	encoder := json.NewEncoder(w)
	if pretty {
		encoder.SetIndent("", "  ")
	}
	if err != nil {
		encoder.Encode(map[string]string{
			"error": err.Error(),
		})
	} else {
		encoder.Encode(map[string][]docker.Stats{"data": stats})
	}
}
