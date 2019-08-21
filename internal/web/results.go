package web

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/G-Node/gin-valid/internal/config"
	"github.com/G-Node/gin-valid/internal/helpers"
	"github.com/G-Node/gin-valid/internal/log"
	"github.com/G-Node/gin-valid/internal/resources/templates"
	"github.com/gorilla/mux"
)

// BidsResultStruct is the struct to parse a full BIDS validation json.
type BidsResultStruct struct {
	Issues struct {
		Errors []struct {
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
				} `json:"file"`
				Evidence  interface{} `json:"evidence"`
				Line      interface{} `json:"line"`
				Character interface{} `json:"character"`
				Severity  string      `json:"severity"`
				Reason    string      `json:"reason"`
			} `json:"files"`
			AdditionalFileCount int `json:"additionalFileCount"`
			Code                int `json:"code"`
		} `json:"errors"`
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
			AdditionalFileCount int `json:"additionalFileCount"`
			Code                int `json:"code"`
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
	vars := mux.Vars(r)
	user := vars["user"]
	repo := vars["repo"]
	validator := strings.ToLower(vars["validator"])
	if !helpers.SupportedValidator(validator) {
		log.ShowWrite("[Error] unsupported validator '%s'\n", validator)
		http.ServeContent(w, r, "unavailable", time.Now(), bytes.NewReader([]byte("404 Nothing to see here...")))
		return
	}
	log.ShowWrite("[Info] '%s' results for repo '%s/%s'\n", validator, user, repo)

	srvcfg := config.Read()
	resdir := filepath.Join(srvcfg.Dir.Result, validator, user, repo, srvcfg.Label.ResultsFolder)

	fp := filepath.Join(resdir, srvcfg.Label.ResultsBadge)
	badge, err := ioutil.ReadFile(fp)
	if err != nil {
		log.ShowWrite("[Error] serving '%s/%s' badge: %s\n", user, repo, err.Error())
	}

	fp = filepath.Join(resdir, srvcfg.Label.ResultsFile)
	content, err := ioutil.ReadFile(fp)
	if err != nil {
		log.ShowWrite("[Error] serving '%s/%s' result: %s\n", user, repo, err.Error())
		http.ServeContent(w, r, "unavailable", time.Now(), bytes.NewReader([]byte("404 Nothing to see here...")))
		return
	}

	// Parse results file
	var resBIDS BidsResultStruct
	err = json.Unmarshal(content, &resBIDS)
	if err != nil {
		log.ShowWrite("[Error] unmarshalling '%s/%s' result: %s\n", user, repo, err.Error())
		http.ServeContent(w, r, "unavailable", time.Now(), bytes.NewReader([]byte("500 Something went wrong...")))
		return
	}

	// Parse html template
	tmpl := template.New("layout")
	tmpl, err = tmpl.Parse(templates.Layout)
	if err != nil {
		log.ShowWrite("[Error] '%s/%s' result: %s\n", user, repo, err.Error())
		http.ServeContent(w, r, "unavailable", time.Now(), bytes.NewReader([]byte("500 Something went wrong...")))
		return
	}
	tmpl, err = tmpl.Parse(templates.BidsResults)
	if err != nil {
		log.ShowWrite("[Error] '%s/%s' result: %s\n", user, repo, err.Error())
		http.ServeContent(w, r, "unavailable", time.Now(), bytes.NewReader([]byte("500 Something went wrong...")))
		return
	}

	// Parse results into html template and serve it
	head := fmt.Sprintf("Validation for %s/%s", user, repo)
	info := struct {
		Badge  template.HTML
		Header string
		*BidsResultStruct
	}{template.HTML(badge), head, &resBIDS}

	err = tmpl.ExecuteTemplate(w, "layout", info)
	if err != nil {
		log.ShowWrite("[Error] '%s/%s' result: %s\n", user, repo, err.Error())
		http.ServeContent(w, r, "unavailable", time.Now(), bytes.NewReader([]byte("500 Something went wrong...")))
		return
	}
}
