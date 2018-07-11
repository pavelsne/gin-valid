package web

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"time"

	"github.com/gorilla/mux"
	"github.com/mpsonntag/gin-valid/config"
	"github.com/mpsonntag/gin-valid/log"
)

// BidsResultStruct is the struct to parse a full BIDS validation json.
type BidsResultStruct struct {
	Issues struct {
		Errors   []interface{} `json:"errors"`
		Warnings []struct {
			Key      string `json:"key"`
			Severity string `json:"severity"`
			Reason   string `json:"reason"`
			Files    []struct {
				Key  string `json:"key"`
				Code int    `json:"code"`
				File struct {
					Name         string `json:"name"`
					Path         string `json:"path"`
					RelativePath string `json:"relativePath"`
					Stats        struct {
						Dev         int       `json:"dev"`
						Mode        int       `json:"mode"`
						Nlink       int       `json:"nlink"`
						UID         int       `json:"uid"`
						Gid         int       `json:"gid"`
						Rdev        int       `json:"rdev"`
						Blksize     int       `json:"blksize"`
						Ino         int       `json:"ino"`
						Size        int       `json:"size"`
						Blocks      int       `json:"blocks"`
						AtimeMs     float64   `json:"atimeMs"`
						MtimeMs     float64   `json:"mtimeMs"`
						CtimeMs     float64   `json:"ctimeMs"`
						BirthtimeMs float64   `json:"birthtimeMs"`
						Atime       time.Time `json:"atime"`
						Mtime       time.Time `json:"mtime"`
						Ctime       time.Time `json:"ctime"`
						Birthtime   time.Time `json:"birthtime"`
					} `json:"stats"`
				} `json:"file"`
				Evidence  interface{} `json:"evidence"`
				Line      interface{} `json:"line"`
				Character interface{} `json:"character"`
				Severity  string      `json:"severity"`
				Reason    string      `json:"reason"`
			} `json:"files"`
			AdditionalFileCount int    `json:"additionalFileCount"`
			Code                string `json:"code"`
		} `json:"warnings"`
		Ignored []interface{} `json:"ignored"`
	} `json:"issues"`
	Summary struct {
		Sessions   []interface{} `json:"sessions"`
		Subjects   []string      `json:"subjects"`
		Tasks      []string      `json:"tasks"`
		Modalities []string      `json:"modalities"`
		TotalFiles int           `json:"totalFiles"`
		Size       int           `json:"size"`
	} `json:"summary"`
}

// Results returns the results of a previously run BIDS validation.
func Results(w http.ResponseWriter, r *http.Request) {
	user := mux.Vars(r)["user"]
	repo := mux.Vars(r)["repo"]
	log.Write("[Info] results for repo '%s/%s'\n", user, repo)

	srvcfg := config.Read()
	fp := filepath.Join(srvcfg.Dir.Result, user, repo, srvcfg.Label.ResultsFolder, srvcfg.Label.ResultsFile)
	content, err := ioutil.ReadFile(fp)
	if err != nil {
		log.Write("[Error] serving '%s/%s' result: %s\n", user, repo, err.Error())
		http.ServeContent(w, r, "unavailable", time.Now(), bytes.NewReader([]byte("404 Nothing to see here...")))
		return
	}

	var parseBIDS BidsRoot
	err = json.Unmarshal(content, &parseBIDS)
	if err != nil {
		log.Write("[Error] unmarshalling '%s/%s' result: %s\n", user, repo, err.Error())
		http.ServeContent(w, r, "unavailable", time.Now(), bytes.NewReader([]byte("500 Something went wrong...")))
		return
	}

	http.ServeContent(w, r, "results", time.Now(), bytes.NewReader(content))
}
