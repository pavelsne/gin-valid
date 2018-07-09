package web

import (
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/gorilla/mux"
	"github.com/mpsonntag/gin-valid/config"
	"github.com/mpsonntag/gin-valid/valutils"
)

// Status returns the status of the latest BIDS validation for
// a provided gin user repository.
func Status(w http.ResponseWriter, r *http.Request) {
	srvconfig := config.Read()

	user := mux.Vars(r)["user"]
	repo := mux.Vars(r)["repo"]

	// Check whether this repo has ever been built
	latestPath := filepath.Join(srvconfig.Dir.Result, user, repo, "latest")
	if !valutils.ValidDirectory(latestPath) {
		fmt.Fprintln(w, "Never built")
		return
	}
	fmt.Fprintln(w, "Going on")
}
