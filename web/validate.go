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

	"github.com/gorilla/mux"
	"github.com/mpsonntag/gin-valid/config"
	"github.com/mpsonntag/gin-valid/resources"
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

// Validate temporarily clones a provided repository from
// a gin server and checks whether the content of the
// repository is a valid BIDS dataset.
// Any cloned files are cleaned up after the check is done.
func Validate(w http.ResponseWriter, r *http.Request) {
	srvconfig := config.Read()

	user := mux.Vars(r)["user"]
	repo := mux.Vars(r)["repo"]
	fmt.Fprintf(os.Stdout, "[Info] validating repo '%s/%s'\n", user, repo)

	cmd := exec.Command("gin", "repoinfo", fmt.Sprintf("%s/%s", user, repo))
	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "[Error] accessing '%s/%s': '%s'\n", user, repo, err.Error())
		return
	}

	tmpdir, err := ioutil.TempDir(srvconfig.Dir.Temp, "bidsval_")
	if err != nil {
		fmt.Fprintf(os.Stderr, "[Error] creating temporary directory: '%s'\n", err.Error())
		return
	}

	// enable cleanup once tried and tested
	defer os.RemoveAll(tmpdir)

	cmd = exec.Command(srvconfig.Exec.Gin, "get", fmt.Sprintf("%s/%s", user, repo))
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Dir = tmpdir
	if err = cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "[Error] running gin get: '%s'\n", err.Error())
		return
	}

	outBadge := filepath.Join(srvconfig.Dir.Result, user, repo, "latest", "results.svg")
	outFile := filepath.Join(srvconfig.Dir.Result, user, repo, "latest", "results.json")

	// Create results folder if necessary
	// CHECK: can this lead to a race condition, if a job for the same user/repo combination is started twice in short succession?
	latestPath := filepath.Join(srvconfig.Dir.Result, user, repo, "latest")
	err = os.MkdirAll(latestPath, os.ModePerm)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[Error] creating latest build folder '%s/%s': %s", user, repo, err.Error())
		// Think about whether we should do something at this point
		return
	}

	// Ignoring NiftiHeaders for now, since it seems to be a common error
	cmd = exec.Command(srvconfig.Exec.BIDS, "--ignoreNiftiHeaders", "--json", fmt.Sprintf("%s/%s", tmpdir, repo))
	out.Reset()
	cmd.Stdout = &out
	var serr bytes.Buffer
	cmd.Stderr = &serr
	cmd.Dir = tmpdir
	if err = cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "[Error] running bids validation (%s): '%s', '%s', '%s'",
			fmt.Sprintf("%s/%s", tmpdir, repo), err.Error(), serr.String(), out.String())

		err = ioutil.WriteFile(outBadge, []byte(resources.BidsFailure), os.ModePerm)
		if err != nil {
			fmt.Fprintf(os.Stderr, "[Error] writing output badge for '%s/%s'\n", user, repo)
		}
		return
	}

	// We need this for both the writing of the result and the badge
	output := out.Bytes()

	// CHECK: can this lead to a race condition, if a job for the same user/repo combination is started twice in short succession?
	err = ioutil.WriteFile(outFile, []byte(output), os.ModePerm)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[Error] writing output file for '%s/%s'\n", user, repo)
	}

	// Write proper badge according to result
	content := resources.BidsSuccess
	var parseBIDS BidsRoot
	err = json.Unmarshal(output, &parseBIDS)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[Error] unmarshalling json: %s", err.Error())
		content = resources.BidsFailure
	} else if len(parseBIDS.Issues.Errors) > 0 {
		content = resources.BidsFailure
	} else if len(parseBIDS.Issues.Warnings) > 0 {
		content = resources.BidsWarning
	}

	err = ioutil.WriteFile(outBadge, []byte(content), os.ModePerm)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[Error] writing output badge for '%s/%s'\n", user, repo)
	}

	fmt.Fprintf(os.Stdout, "[Info] finished validating repo '%s/%s'\n", user, repo)

	_, err = w.Write([]byte(content))
	if err != nil {
		fmt.Fprintf(os.Stderr, "[Error] returning badge %s", err.Error())
	}
}
