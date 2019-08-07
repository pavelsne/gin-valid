package web

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"time"

	"github.com/G-Node/gin-valid/internal/config"
	"github.com/G-Node/gin-valid/internal/helpers"
	"github.com/G-Node/gin-valid/internal/log"
	"github.com/G-Node/gin-valid/internal/resources"
	"github.com/gorilla/mux"
)

// Status returns the status of the latest BIDS validation for
// a provided gin user repository.
func Status(w http.ResponseWriter, r *http.Request) {
	validator := mux.Vars(r)["validator"]
	if !helpers.SupportedValidator(validator) {
		log.Write("[Error] unsupported validator '%s'\n", validator)
		http.ServeContent(w, r, "unavailable", time.Now(), bytes.NewReader([]byte("404 Nothing to see here...")))
		return
	}
	user := mux.Vars(r)["user"]
	repo := mux.Vars(r)["repo"]
	log.Write("[Info] '%s' status for repo '%s/%s'\n", validator, user, repo)

	srvcfg := config.Read()

	fp := filepath.Join(srvcfg.Dir.Result, "bids", user, repo, srvcfg.Label.ResultsFolder, srvcfg.Label.ResultsBadge)
	content, err := ioutil.ReadFile(fp)
	if err != nil {
		log.Write("[Error] serving '%s/%s' status: %s\n", user, repo, err.Error())
		http.ServeContent(w, r, "unavailable.svg", time.Now(), bytes.NewReader([]byte(resources.BidsUnavailable)))
		return
	}
	http.ServeContent(w, r, srvcfg.Label.ResultsBadge, time.Now(), bytes.NewReader(content))
}
