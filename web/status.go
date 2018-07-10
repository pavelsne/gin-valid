package web

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
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

	fp := filepath.Join(srvconfig.Dir.Result, user, repo, "latest", "results.svg")
	content, err := ioutil.ReadFile(fp)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[Error] unable to serve '%s/%s' badge: %s\n", user, repo, err.Error())
		http.ServeContent(w, r, "unavailable.svg", time.Now(), bytes.NewReader([]byte(resources.BidsUnavailable)))
		return
	}
	http.ServeContent(w, r, "results.svg", time.Now(), bytes.NewReader(content))
}
