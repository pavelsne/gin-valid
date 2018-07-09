package web

import (
	"bytes"
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

// Validate temporarily clones a provided repository from
// a gin server and checks whether the content of the
// repository is a valid BIDS dataset.
// Any cloned files are cleaned up after the check is done.
func Validate(w http.ResponseWriter, r *http.Request) {
	srvconfig := config.Read()

	user := mux.Vars(r)["user"]
	repo := mux.Vars(r)["repo"]
	fmt.Fprintf(w, "validate repo '%s/%s'\n", user, repo)
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
	fmt.Fprintf(w, "Directory created: %s\n", tmpdir)

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
	fmt.Fprintf(w, "running in %s, gin get: %s\n", cmd.Dir, out.String())

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
			fmt.Fprintf(os.Stderr, "[Error] Was not able to write output badge for '%s/%s'\n", user, repo)
		}
		return
	}
	fmt.Fprintln(w, "validation successful")

	// CHECK: can this lead to a race condition, if a job for the same user/repo combination is started twice in short succession?
	err = ioutil.WriteFile(outFile, out.Bytes(), os.ModePerm)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[Error] Was not able to write output file for '%s/%s'\n", user, repo)
	}
	err = ioutil.WriteFile(outBadge, []byte(resources.BidsSuccess), os.ModePerm)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[Error] Was not able to write output badge for '%s/%s'\n", user, repo)
	}
}
