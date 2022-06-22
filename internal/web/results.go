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
					string       `json:"name"`
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

// Results returns the results of a previously run validation.
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
	resID, ok := vars["id"]
	if !ok {
		fmt.Println("Results ID not specified: Rendering default")
		resID = srvcfg.Label.ResultsFolder
	}
	resdir := filepath.Join(srvcfg.Dir.Result, validator, user, repo, resID)

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

	if string(content) == progressmsg {
		// validation in progress
		renderInProgress(w, r, badge, strings.ToUpper(validator), user, repo)
		return
	}

	switch validator {
	case "bids":
		renderBIDSResults(w, r, badge, content, user, repo)
	case "nix":
		renderNIXResults(w, r, badge, content, user, repo)
	case "odml":
		renderODMLResults(w, r, badge, content, user, repo)
	default:
		log.ShowWrite("[Error] Validator %q is supported but no render result function is set up", validator)
		http.ServeContent(w, r, "unavailable", time.Now(), bytes.NewReader([]byte("404 Validator results missing")))
	}
	return
}

func renderInProgress(w http.ResponseWriter, r *http.Request, badge []byte, validator string, user, repo string) {
	tmpl := template.New("layout")
	tmpl, err := tmpl.Parse(templates.Layout)
	if err != nil {
		log.ShowWrite("[Error] '%s/%s' result: %s\n", user, repo, err.Error())
		http.ServeContent(w, r, "unavailable", time.Now(), bytes.NewReader([]byte("500 Something went wrong...")))
		return
	}
	tmpl, err = tmpl.Parse(templates.GenericResults)
	if err != nil {
		log.ShowWrite("[Error] '%s/%s' result: %s\n", user, repo, err.Error())
		http.ServeContent(w, r, "unavailable", time.Now(), bytes.NewReader([]byte("500 Something went wrong...")))
		return
	}

	// Parse results into html template and serve it
	head := fmt.Sprintf("%s validation for %s/%s", validator, user, repo)
	srvcfg := config.Read()
	year, _, _ := time.Now().Date()
	info := struct {
		Badge       template.HTML
		Header      string
		Content     string
		GinURL      string
		CurrentYear int
	}{template.HTML(badge), head, string(progressmsg), srvcfg.GINAddresses.WebURL, year}

	err = tmpl.ExecuteTemplate(w, "layout", info)
	if err != nil {
		log.ShowWrite("[Error] '%s/%s' result: %s\n", user, repo, err.Error())
		http.ServeContent(w, r, "unavailable", time.Now(), bytes.NewReader([]byte("500 Something went wrong...")))
		return
	}
}

func renderBIDSResults(w http.ResponseWriter, r *http.Request, badge []byte, content []byte, user, repo string) {
	// Parse results file
	var resBIDS BidsResultStruct
	err := json.Unmarshal(content, &resBIDS)
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
	head := fmt.Sprintf("BIDS validation for %s/%s", user, repo)
	year, _, _ := time.Now().Date()
	srvcfg := config.Read()
	info := struct {
		Badge  template.HTML
		Header string
		*BidsResultStruct
		GinURL      string
		CurrentYear int
	}{template.HTML(badge), head, &resBIDS, srvcfg.GINAddresses.WebURL, year}

	err = tmpl.ExecuteTemplate(w, "layout", info)
	if err != nil {
		log.ShowWrite("[Error] '%s/%s' result: %s\n", user, repo, err.Error())
		http.ServeContent(w, r, "unavailable", time.Now(), bytes.NewReader([]byte("500 Something went wrong...")))
		return
	}
}

func renderNIXResults(w http.ResponseWriter, r *http.Request, badge []byte, content []byte, user, repo string) {
	// Parse results file
	// Parse html template
	tmpl := template.New("layout")
	tmpl, err := tmpl.Parse(templates.Layout)
	if err != nil {
		log.ShowWrite("[Error] '%s/%s' result: %s\n", user, repo, err.Error())
		http.ServeContent(w, r, "unavailable", time.Now(), bytes.NewReader([]byte("500 Something went wrong...")))
		return
	}
	tmpl, err = tmpl.Parse(templates.GenericResults)
	if err != nil {
		log.ShowWrite("[Error] '%s/%s' result: %s\n", user, repo, err.Error())
		http.ServeContent(w, r, "unavailable", time.Now(), bytes.NewReader([]byte("500 Something went wrong...")))
		return
	}

	// Parse results into html template and serve it
	head := fmt.Sprintf("NIX validation for %s/%s", user, repo)
	year, _, _ := time.Now().Date()
	srvcfg := config.Read()
	info := struct {
		Badge       template.HTML
		Header      string
		Content     string
		GinURL      string
		CurrentYear int
	}{template.HTML(badge), head, string(content), srvcfg.GINAddresses.WebURL, year}

	err = tmpl.ExecuteTemplate(w, "layout", info)
	if err != nil {
		log.ShowWrite("[Error] '%s/%s' result: %s\n", user, repo, err.Error())
		http.ServeContent(w, r, "unavailable", time.Now(), bytes.NewReader([]byte("500 Something went wrong...")))
		return
	}
}

func renderODMLResults(w http.ResponseWriter, r *http.Request, badge []byte, content []byte, user, repo string) {
	// Parse results file
	// Parse html template
	tmpl := template.New("layout")
	tmpl, err := tmpl.Parse(templates.Layout)
	if err != nil {
		log.ShowWrite("[Error] '%s/%s' result: %s\n", user, repo, err.Error())
		http.ServeContent(w, r, "unavailable", time.Now(), bytes.NewReader([]byte("500 Something went wrong...")))
		return
	}
	tmpl, err = tmpl.Parse(templates.GenericResults)
	if err != nil {
		log.ShowWrite("[Error] '%s/%s' result: %s\n", user, repo, err.Error())
		http.ServeContent(w, r, "unavailable", time.Now(), bytes.NewReader([]byte("500 Something went wrong...")))
		return
	}

	// Parse results into html template and serve it
	head := fmt.Sprintf("odML validation for %s/%s", user, repo)
	srvcfg := config.Read()
	year, _, _ := time.Now().Date()
	info := struct {
		Badge       template.HTML
		Header      string
		Content     string
		GinURL      string
		CurrentYear int
	}{template.HTML(badge), head, string(content), srvcfg.GINAddresses.WebURL, year}

	err = tmpl.ExecuteTemplate(w, "layout", info)
	if err != nil {
		log.ShowWrite("[Error] '%s/%s' result: %s\n", user, repo, err.Error())
		http.ServeContent(w, r, "unavailable", time.Now(), bytes.NewReader([]byte("500 Something went wrong...")))
		return
	}
}
