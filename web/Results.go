package web

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mpsonntag/gin-valid/log"
)

// Results returns the results of a previously run BIDS validation.
func Results(w http.ResponseWriter, r *http.Request) {
	user := mux.Vars(r)["user"]
	repo := mux.Vars(r)["repo"]

	log.ShowWrite("[Info] delivering validation results for '%s/%s'", user, repo)
}
