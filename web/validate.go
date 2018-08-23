package web

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"gopkg.in/yaml.v2"

	"github.com/G-Node/gin-cli/ginclient"
	glog "github.com/G-Node/gin-cli/ginclient/log"
	"github.com/G-Node/gin-cli/git"
	"github.com/G-Node/gin-cli/git/shell"
	"github.com/G-Node/gin-valid/config"
	"github.com/G-Node/gin-valid/helpers"
	"github.com/G-Node/gin-valid/log"
	"github.com/G-Node/gin-valid/resources"
	"github.com/G-Node/gin-valid/resources/templates"
	gogs "github.com/gogits/go-gogs-client"
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

// Validationcfg is used to unmarshall a config file
// holding information specific for running the
// various validations. e.g. where the root
// folder of a bids directory can be found or
// whether the NiftiHeaders should be ignored.
type Validationcfg struct {
	Bidscfg struct {
		BidsRoot      string `yaml:"bidsroot"`
		ValidateNifti bool   `yaml:"validatenifti"`
	} `yaml:"bidsconfig"`
}

// handleValidationConfig unmarshalles a yaml config file
// from file and returns the resulting Validationcfg struct.
func handleValidationConfig(cfgpath string) (Validationcfg, error) {
	valcfg := Validationcfg{}

	content, err := ioutil.ReadFile(cfgpath)
	if err != nil {
		return valcfg, err
	}

	err = yaml.Unmarshal(content, &valcfg)
	if err != nil {
		return valcfg, err
	}

	return valcfg, nil
}

func runValidator(validator, repopath, commit string, gcl *ginclient.Client) (int, error) {
	log.Write("[Info] Commit hash: %s", commit)

	repopathparts := strings.SplitN(repopath, "/", 2)
	user, repo := repopathparts[0], repopathparts[1]

	srvcfg := config.Read()

	// TODO add check if a repo is currently being validated. Since the cloning
	// can potentially take quite some time prohibit running the same
	// validation at the same time. Could also move this to a mapped go
	// routine and if the same repo is validated twice, the first occurrence is
	// stopped and cleaned up while the second starts anew - to make sure its
	// always the latest state of the repository that is being validated.

	// TODO: Use the payload data to check if the specific commit has already
	// been validated

	_, err := gcl.GetRepo(repopath)
	if err != nil {
		code, converr := strconv.Atoi(err.Error()[:3])
		if converr != nil {
			// First three chars should be error code. If not, default to NotFound
			code = http.StatusNotFound
		}
		log.Write("[Error] Repository not found: %s", repopath)
		return code, fmt.Errorf("accessing '%s': %s", repopath, err.Error())
	}

	log.Write("[Info] Found repository on server")

	tmpdir, err := ioutil.TempDir(srvcfg.Dir.Temp, validator)
	if err != nil {
		log.Write("[Error] Internal error: Couldn't create temporary gin directory: %s", err.Error())
		return http.StatusInternalServerError, fmt.Errorf("validation on %s failed", repopath)
	}

	// Enable cleanup once tried and tested
	defer os.RemoveAll(tmpdir)

	// TODO: if (annexed) content is not available yet, wait and retry.  We
	// would have to set a max timeout for this.  The issue is that when a user
	// does a 'gin upload' a push happens immediately and the hook is
	// triggered, but annexed content is only transferred after the push and
	// could take a while (hours?). The validation service should try to
	// download content after the transfer is complete, or should keep retrying
	// until it's available, with a timeout. We could also make it more
	// efficient by only downloading the content in the directories which are
	// specified in the validator config (if it exists).

	glog.Init("")
	clonechan := make(chan git.RepoFileStatus)
	os.Chdir(tmpdir)
	go gcl.CloneRepo(repopath, clonechan)
	for stat := range clonechan {
		if stat.Err != nil {
			e := stat.Err.(shell.Error)
			log.Write("[Error] %s", e.UError)
			log.Write("[Error] %s", e.Description)
			log.Write("[Error] %s", e.Origin)
			return http.StatusInternalServerError, fmt.Errorf("failed to fetch repository data")
		}
		log.Write("[Info] %s %s", stat.State, stat.Progress)
	}
	log.Write("[Info] clone complete for '%s'", repopath)

	// checkout specific commit then download all content
	log.Write("[Info] git checkout %s", commit)
	err = git.Checkout(commit, nil)
	if err != nil {
		log.Write("[Error] failed to checkout commit '%s': %s", commit, err.Error())
		return http.StatusInternalServerError, fmt.Errorf("failed to fetch repository data")
	}

	log.Write("[Info] Downloading content")
	getcontentchan := make(chan git.RepoFileStatus)
	go gcl.GetContent([]string{"."}, getcontentchan)
	for stat := range getcontentchan {
		if stat.Err != nil {
			e := stat.Err.(shell.Error)
			log.Write("[Error] %s", e.UError)
			log.Write("[Error] %s", e.Description)
			log.Write("[Error] %s", e.Origin)
			return http.StatusInternalServerError, fmt.Errorf("failed to fetch repository data")
		}
		log.Write("[Info] %s %s %s", stat.State, stat.FileName, stat.Progress)
	}
	log.Write("[Info] get-content complete")

	// Create results folder if necessary
	// CHECK: can this lead to a race condition, if a job for the same user/repo combination is started twice in short succession?
	resdir := filepath.Join(srvcfg.Dir.Result, validator, repopath, srvcfg.Label.ResultsFolder)
	err = os.MkdirAll(resdir, os.ModePerm)
	if err != nil {
		log.Write("[Error] creating '%s' results folder: %s", repopath, err.Error())
		return http.StatusInternalServerError, fmt.Errorf("failed to generate results")
	}

	// Use validation config file if available
	valroot := filepath.Join(tmpdir, repo)
	var validateNifti bool

	cfgpath := filepath.Join(tmpdir, repo, srvcfg.Label.ValidationConfigFile)
	log.Write("[Info] looking for config file at '%s'", cfgpath)
	if fi, err := os.Stat(cfgpath); err == nil && !fi.IsDir() {
		valcfg, err := handleValidationConfig(cfgpath)
		if err == nil {
			checkdir := filepath.Join(tmpdir, repo, valcfg.Bidscfg.BidsRoot)
			if fi, err = os.Stat(checkdir); err == nil && fi.IsDir() {
				valroot = checkdir
				log.Write("[Info] using validation root directory: %s\n%s", valroot, checkdir)
			} else {
				log.Write("[Error] reading validation root directory: %s", err.Error())
			}
			validateNifti = valcfg.Bidscfg.ValidateNifti
		} else {
			log.Write("[Error] unmarshalling validation config file: %s", err.Error())
		}
	} else {
		log.Write("[Info] no validation config file found or processed, running from repo root (%s)", err.Error())
	}

	// Ignoring NiftiHeaders for now, since it seems to be a common error
	outBadge := filepath.Join(resdir, srvcfg.Label.ResultsBadge)
	log.Write("[Info] Running bids validation: '%s %t --json %s'", srvcfg.Exec.BIDS, validateNifti, valroot)

	// Make sure the validator arguments are in the right order
	var args []string
	if !validateNifti {
		args = append(args, "--ignoreNiftiHeaders")
	}
	args = append(args, "--json")
	args = append(args, valroot)

	// cmd = exec.Command(srvcfg.Exec.BIDS, validateNifti, "--json", valroot)
	var out, serr bytes.Buffer
	cmd := exec.Command(srvcfg.Exec.BIDS, args...)
	out.Reset()
	serr.Reset()
	cmd.Stdout = &out
	cmd.Stderr = &serr
	cmd.Dir = tmpdir
	if err = cmd.Run(); err != nil {
		log.Write("[Error] running bids validation (%s): '%s', '%s', '%s'",
			valroot, err.Error(), serr.String(), out.String())

		err = ioutil.WriteFile(outBadge, []byte(resources.BidsFailure), os.ModePerm)
		if err != nil {
			log.Write("[Error] writing results badge for '%s'", repopath)
		}
		return 0, err
	}

	// We need this for both the writing of the result and the badge
	output := out.Bytes()

	// CHECK: can this lead to a race condition, if a job for the same user/repo combination is started twice in short succession?
	outFile := filepath.Join(resdir, srvcfg.Label.ResultsFile)
	err = ioutil.WriteFile(outFile, []byte(output), os.ModePerm)
	if err != nil {
		log.Write("[Error] writing results file for '%s'", repopath)
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
		log.Write("[Error] writing results badge for '%s/%s'", user, repo)
		return 0, err
	}

	log.Write("[Info] finished validating repo '%s/%s'", user, repo)
	return 0, nil
}

// Root renders the root page of the gin-valid service, which allows the user
// to manually run a validator on a publicly accessible repository, without
// using a web hook.
func Root(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		tmpl := template.New("layout")
		tmpl, err := tmpl.Parse(templates.Layout)
		if err != nil {
			log.Write("[Error] failed to parse html layout page")
			fail(w, http.StatusInternalServerError, "something went wrong")
			return
		}
		tmpl, err = tmpl.Parse(templates.Root)
		if err != nil {
			log.Write("[Error] failed to render root page")
			fail(w, http.StatusInternalServerError, "something went wrong")
			return
		}
		tmpl.Execute(w, nil)
	}
}

// PubValidate parses the POST data from the root form and calls the validator
// using the built-in ServiceWaiter.
func PubValidate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		// Do nothing
		log.Write("[Error] no post request: %s", r.Method)
		return
	}

	srvcfg := config.Read()
	ginuser := srvcfg.Settings.GINUser

	r.ParseForm()
	repopath := r.Form["repopath"][0]
	validator := "bids" // vars["validator"] // TODO: add options to root form

	log.Write("[Info] About to validate repository '%s' with %s", repopath, ginuser)
	log.Write("[Info] Logging in to GIN server")
	gcl := ginclient.New(serveralias)
	err := gcl.Login(ginuser, srvcfg.Settings.GINPassword, srvcfg.Settings.ClientID)
	if err != nil {
		log.Write("[error] failed to login as %s", ginuser)
		msg := fmt.Sprintf("failed to validate '%s': %s", repopath, err.Error())
		fail(w, http.StatusUnauthorized, msg)
		return
	}
	defer gcl.Logout()

	// check if repository is accessible
	repoinfo, err := gcl.GetRepo(repopath)
	if err != nil {
		fail(w, http.StatusNotFound, err.Error())
		return
	}
	if repoinfo.Private {
		// We (the built in ServiceWaiter) have access, but the repository is
		// marked private. This can happen if an owner of the private
		// repository adds the user as a collaborator to the repository. We
		// don't allow this.
		fail(w, http.StatusNotFound, fmt.Sprintf("repository '%s' does not exist", repopath))
		return
	}

	statuscode, err := runValidator(validator, repopath, "HEAD", gcl)
	if err != nil {
		if statuscode == 0 {
			statuscode = http.StatusInternalServerError
		}
		fail(w, statuscode, err.Error())
		return
	}

	// TODO redirect to results
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

// Validate temporarily clones a provided repository from
// a gin server and checks whether the content of the
// repository is a valid BIDS dataset.
// Any cloned files are cleaned up after the check is done.
func Validate(w http.ResponseWriter, r *http.Request) {
	// TODO: Simplify/split this function
	if r.Method != http.MethodPost {
		// Do nothing
		log.Write("[Error] no post request: %s", r.Method)
		// TODO: Redirect to results
		return
	}

	secret := r.Header.Get("X-Gogs-Signature")

	var hookdata gogs.PushPayload
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Write("[Error] failed to parse hook payload")
		fail(w, http.StatusBadRequest, "bad request")
		return
	}
	err = json.Unmarshal(b, &hookdata)
	if err != nil {
		log.Write("[Error] failed to parse hook payload")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("bad request"))
		return
	}
	if !validateHookSecret(b, secret) {
		log.Write("[Error] authorisation failed: bad secret")
		fail(w, http.StatusBadRequest, "bad request")
		return
	}

	commithash := hookdata.After

	log.Write("[Info] Hook secret: %s", secret)
	log.Write("[Info] Commit hash: %s", commithash)

	vars := mux.Vars(r)
	validator := vars["validator"]
	if !helpers.SupportedValidator(validator) {
		fail(w, http.StatusNotFound, "unsupported validator")
		return
	}
	user := vars["user"]
	repo := vars["repo"]
	repopath := fmt.Sprintf("%s/%s", user, repo)
	log.Write("[Info] '%s' validation for repo '%s'", validator, repopath)

	// TODO add check if a repo is currently being validated. Since the cloning
	// can potentially take quite some time prohibit running the same
	// validation at the same time. Could also move this to a mapped go
	// routine and if the same repo is validated twice, the first occurrence is
	// stopped and cleaned up while the second starts anew - to make sure its
	// always the latest state of the repository that is being validated.

	// TODO: Use the payload data to check if the specific commit has already
	// been validated

	// get the username + token from the registered hooks map

	ut, ok := hookregs[repopath]
	if !ok {
		// We don't have a valid token for this repository: can't clone
		msg := fmt.Sprintf("accessing '%s': no access token found", repopath)
		fail(w, http.StatusUnauthorized, msg)
		return
	}
	log.Write("[Info] Using user %s", ut.Username)
	gcl := ginclient.New(serveralias)
	gcl.UserToken = ut
	log.Write("[Info] Got user %s. Checking repo", gcl.Username)
	// TODO: make key with unique name in tmp and delete when done
	// Currently, multiple simultaneous validations will override each-others keys
	err = gcl.MakeSessionKey()
	if err != nil {
		log.Write("[error] failed to create session key")
		msg := fmt.Sprintf("failed to clone '%s': %s", repopath, err.Error())
		fail(w, http.StatusUnauthorized, msg)
		return
	}
	defer deleteSessionKey(gcl)

	statuscode, err := runValidator(validator, repopath, commithash, gcl)
	if err != nil {
		if statuscode == 0 {
			statuscode = http.StatusInternalServerError
		}
		w.WriteHeader(statuscode)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
