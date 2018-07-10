package web

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"time"

	"github.com/gorilla/mux"
	"github.com/mpsonntag/gin-valid/config"
	"github.com/mpsonntag/gin-valid/log"
	"github.com/mpsonntag/gin-valid/resources"
)

// Status returns the status of the latest BIDS validation for
// a provided gin user repository.
func Status(w http.ResponseWriter, r *http.Request) {
	srvconfig := config.Read()

	user := mux.Vars(r)["user"]
	repo := mux.Vars(r)["repo"]

	fp := filepath.Join(srvconfig.Dir.Result, user, repo, srvconfig.Label.ResultsFolder, srvconfig.Label.ResultsBadge)
	content, err := ioutil.ReadFile(fp)
	if err != nil {
		log.Write("[Error] serving '%s/%s' status: %s\n", user, repo, err.Error())
		http.ServeContent(w, r, "unavailable.svg", time.Now(), bytes.NewReader([]byte(resources.BidsUnavailable)))
		return
	}
	http.ServeContent(w, r, srvconfig.Label.ResultsBadge, time.Now(), bytes.NewReader(content))
}
