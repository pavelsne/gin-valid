package web

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/G-Node/gin-valid/config"
	"github.com/G-Node/gin-valid/helpers"
	"github.com/G-Node/gin-valid/log"
	"github.com/G-Node/gin-valid/resources"
	"github.com/gorilla/mux"
)

// BidsMessages contains Errors, Warnings and Ignored messages.
// Currently its just the number of individual messages
// we are interested in. If this changes, the messages
// will be expanded into proper structs of their own.
type BidsMessages struct {
	Errors   []interface{} `json:"errors"`
	Warnings []interface{} `json:"warnings"`
	Ignored  []interface{} `json:"ignored"`
}

// BidsRoot contains only the root issues element.
type BidsRoot struct {
	Issues BidsMessages `json:"issues"`
}

// unavailable creates a log entry and writes the unavailable badge to the responseWriter.
func unavailable(w http.ResponseWriter, r *http.Request, badge string, message string) {
	log.Write(message)
	http.ServeContent(w, r, badge, time.Now(), bytes.NewReader([]byte(resources.BidsUnavailable)))
}

// Validate temporarily clones a provided repository from
// a gin server and checks whether the content of the
// repository is a valid BIDS dataset.
// Any cloned files are cleaned up after the check is done.
func Validate(w http.ResponseWriter, r *http.Request) {
	service := mux.Vars(r)["service"]
	if !helpers.SupportedValidator(service) {
		log.Write("[Error] unsupported validator '%s'\n", service)
		http.ServeContent(w, r, "unavailable", time.Now(), bytes.NewReader([]byte("404 Nothing to see here...")))
		return
	}
	user := mux.Vars(r)["user"]
	repo := mux.Vars(r)["repo"]
	log.Write("[Info] '%s' validation for repo '%s/%s'\n", service, user, repo)

	srvcfg := config.Read()

	var out bytes.Buffer
	var serr bytes.Buffer

	// login with gin servce user, gaining access to public repositories for now
	cmd := exec.Command(srvcfg.Exec.Gin, "login", "ServiceWaiter")
	cmd.Stdin = strings.NewReader(fmt.Sprintln(srvcfg.Settings.GPW))
	cmd.Stdout = &out
	cmd.Stderr = &serr
	if err := cmd.Run(); err != nil {
		msg := fmt.Sprintf("[Error] logging into gin: '%s, %s'\n", out.String(), serr.String())
		unavailable(w, r, srvcfg.Label.ResultsBadge, msg)
		return
	}

	// TODO add check if a repo is currently being validated. since
	// the cloning can potentially take quite some time prohibit
	// running the same validation at the same time.
	// could also move this to a mapped go routine and if the same
	// repo is validated twice, the first occurrence is stopped and
	// cleaned up while the second starts anew - to make sure its always
	// the latest state of the repository that is being validated.

	cmd = exec.Command(srvcfg.Exec.Gin, "repoinfo", fmt.Sprintf("%s/%s", user, repo))
	out.Reset()
	serr.Reset()
	cmd.Stdout = &out
	cmd.Stderr = &serr
	if err := cmd.Run(); err != nil {
		msg := fmt.Sprintf("[Error] accessing '%s/%s': '%s, %s'\n", user, repo, out.String(), serr.String())
		unavailable(w, r, srvcfg.Label.ResultsBadge, msg)
		return
	}

	tmpdir, err := ioutil.TempDir(srvcfg.Dir.Temp, "bidsval_")
	if err != nil {
		msg := fmt.Sprintf("[Error] creating temp gin directory: '%s'\n", err.Error())
		unavailable(w, r, srvcfg.Label.ResultsBadge, msg)
		return
	}

	// enable cleanup once tried and tested
	defer os.RemoveAll(tmpdir)

	cmd = exec.Command(srvcfg.Exec.Gin, "get", fmt.Sprintf("%s/%s", user, repo))
	out.Reset()
	serr.Reset()
	cmd.Stdout = &out
	cmd.Stderr = &serr
	cmd.Dir = tmpdir
	if err = cmd.Run(); err != nil {
		msg := fmt.Sprintf("[Error] running gin get: '%s', '%s'\n", out.String(), serr.String())
		unavailable(w, r, srvcfg.Label.ResultsBadge, msg)
		return
	}

	// Create results folder if necessary
	// CHECK: can this lead to a race condition, if a job for the same user/repo combination is started twice in short succession?
	resdir := filepath.Join(srvcfg.Dir.Result, "bids", user, repo, srvcfg.Label.ResultsFolder)
	err = os.MkdirAll(resdir, os.ModePerm)
	if err != nil {
		msg := fmt.Sprintf("[Error] creating '%s/%s' results folder: %s", user, repo, err.Error())
		unavailable(w, r, srvcfg.Label.ResultsBadge, msg)
		return
	}

	// Ignoring NiftiHeaders for now, since it seems to be a common error
	outBadge := filepath.Join(resdir, srvcfg.Label.ResultsBadge)
	cmd = exec.Command(srvcfg.Exec.BIDS, "--ignoreNiftiHeaders", "--json", fmt.Sprintf("%s/%s", tmpdir, repo))
	out.Reset()
	serr.Reset()
	cmd.Stdout = &out
	cmd.Stderr = &serr
	cmd.Dir = tmpdir
	if err = cmd.Run(); err != nil {
		log.Write("[Error] running bids validation (%s/%s): '%s', '%s', '%s'",
			tmpdir, repo, err.Error(), serr.String(), out.String())

		err = ioutil.WriteFile(outBadge, []byte(resources.BidsFailure), os.ModePerm)
		if err != nil {
			log.Write("[Error] writing results badge for '%s/%s'\n", user, repo)
		}
		return
	}

	// We need this for both the writing of the result and the badge
	output := out.Bytes()

	// CHECK: can this lead to a race condition, if a job for the same user/repo combination is started twice in short succession?
	outFile := filepath.Join(resdir, srvcfg.Label.ResultsFile)
	err = ioutil.WriteFile(outFile, []byte(output), os.ModePerm)
	if err != nil {
		log.Write("[Error] writing results file for '%s/%s'\n", user, repo)
	}

	// Write proper badge according to result
	content := resources.BidsSuccess
	var parseBIDS BidsRoot
	err = json.Unmarshal(output, &parseBIDS)
	if err != nil {
		log.Write("[Error] unmarshalling results json: %s", err.Error())
		content = resources.BidsFailure
	} else if len(parseBIDS.Issues.Errors) > 0 {
		content = resources.BidsFailure
	} else if len(parseBIDS.Issues.Warnings) > 0 {
		content = resources.BidsWarning
	}

	err = ioutil.WriteFile(outBadge, []byte(content), os.ModePerm)
	if err != nil {
		log.Write("[Error] writing results badge for '%s/%s'\n", user, repo)
	}

	log.Write("[Info] finished validating repo '%s/%s'\n", user, repo)

	http.ServeContent(w, r, srvcfg.Label.ResultsBadge, time.Now(), bytes.NewReader([]byte(content)))
}
