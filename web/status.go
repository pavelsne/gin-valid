package web

import (
	"bytes"
	"net/http"
	"path/filepath"
	"time"

	"github.com/gorilla/mux"
	"github.com/mpsonntag/gin-valid/config"
	"github.com/mpsonntag/gin-valid/resources"
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
		http.ServeContent(w, r, "unavailable.svg", time.Now(), bytes.NewReader([]byte(resources.BidsUnavailable)))
		return
	}
}
