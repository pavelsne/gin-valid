package web

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"time"

	"github.com/G-Node/gin-cli/ginclient"
	gcfg "github.com/G-Node/gin-cli/ginclient/config"
	glog "github.com/G-Node/gin-cli/ginclient/log"
	gweb "github.com/G-Node/gin-cli/web"
	"github.com/G-Node/gin-valid/config"
	"github.com/G-Node/gin-valid/helpers"
	"github.com/G-Node/gin-valid/log"
	"github.com/G-Node/gin-valid/resources/templates"
	gogs "github.com/gogits/go-gogs-client"
	"github.com/gorilla/mux"
)

type repoHooksInfo struct {
	gogs.Repository
	Hooks map[string]bool
}

func cookieExp() time.Time {
	return time.Now().Add(7 * 24 * time.Hour)
}

func deleteSessionKey(gcl *ginclient.Client) {
	hostname, err := os.Hostname()
	if err != nil {
		log.Write("Could not retrieve hostname")
		hostname = "(unknown)"
	}
	description := fmt.Sprintf("GIN Client: %s@%s", gcl.Username, hostname)
	gcl.DeletePubKeyByTitle(description)
	configpath, _ := gcfg.Path(false)
	keyfilepath := filepath.Join(configpath, fmt.Sprintf("%s.key", serveralias))
	os.Remove(keyfilepath)
}

// generateNewSessionID simply generates a secure random 64-byte string (b64 encoded)
func generateNewSessionID() (string, error) {
	length := 64
	id := make([]byte, length)
	if _, err := io.ReadFull(rand.Reader, id); err != nil {
		// This will bubble up and result in an authentication failure. Is
		// there a better message to display to the user? Perhaps 500?
		log.Write("[error] Failed to generate random session ID")
		return "", err
	}
	return base64.StdEncoding.EncodeToString(id), nil
}

func doLogin(username, password string) (*usersession, error) {
	gincl := ginclient.New(serveralias)
	glog.Init("")
	cfg := config.Read()
	clientID := cfg.Settings.ClientID

	// retrieve user's active tokens
	log.Write("Retrieving tokens for user '%s'", username)
	tokens, err := gincl.GetTokens(username, password)
	if err != nil {
		return nil, err
	}

	// check if we have a gin-valid token
	log.Write("Checking for existing token")
	for _, token := range tokens {
		if token.Name == clientID {
			// found our token
			gincl.UserToken.Username = username
			gincl.UserToken.Token = token.Sha1
			log.Write("Found %s access token", clientID)
			break
		}
	}

	if len(gincl.UserToken.Token) == 0 {
		// no existing token; creating new one
		log.Write("Requesting new token from server")
		glog.Write("Performing login from gin-valid")
		err = gincl.NewToken(username, password, clientID)
		if err != nil {
			return nil, err
		}
		log.Write("Login successful. Username: %s", username)
	}

	sessionid, err := generateNewSessionID()
	if err != nil {
		return nil, err
	}
	return &usersession{sessionid, gincl.UserToken}, nil
}

// Login renders the login form and logs in the user to the GIN server, storing
// a session token.
func Login(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		log.Write("Login page")
		tmpl := template.New("layout")
		tmpl, err := tmpl.Parse(templates.Layout)
		if err != nil {
			log.Write("[Error] failed to parse html layout page")
			fail(w, http.StatusInternalServerError, "something went wrong")
			return
		}
		tmpl, err = tmpl.Parse(templates.Login)
		if err != nil {
			log.Write("[Error] failed to render login page")
			fail(w, http.StatusInternalServerError, "something went wrong")
			return
		}
		tmpl.Execute(w, nil)
	} else if r.Method == http.MethodPost {
		log.Write("Doing login")
		r.ParseForm()
		username := r.Form["username"][0]
		password := r.Form["password"][0]
		session, err := doLogin(username, password)
		if err != nil {
			log.Write("[error] Login failed: %s", err.Error())
			fail(w, http.StatusUnauthorized, "authentication failed")
			return
		}

		cfg := config.Read()
		sessions[session.sessionID] = session
		cookie := http.Cookie{
			Name:    cfg.Settings.CookieName,
			Value:   session.sessionID,
			Expires: cookieExp(),
			Secure:  false, // TODO: Switch when we go live
		}
		http.SetCookie(w, &cookie)
		// Redirect to repo listing
		http.Redirect(w, r, fmt.Sprintf("/repos/%s", username), http.StatusFound)
	}
}

func getSessionOrRedirect(w http.ResponseWriter, r *http.Request) (*usersession, error) {
	cfg := config.Read()
	cookie, err := r.Cookie(cfg.Settings.CookieName)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return nil, fmt.Errorf("No session cookie found")
	}
	session, ok := sessions[cookie.Value]
	if !ok {
		http.Redirect(w, r, "/login", http.StatusFound)
		return nil, fmt.Errorf("Invalid session found in cookie")
	}
	return session, nil
}

// ListRepos queries the GIN server for a list of repositories owned (or
// accessible) by a given user and renders the page which displays the
// repositories and their validation status.
func ListRepos(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		return
	}

	session, err := getSessionOrRedirect(w, r)
	if err != nil {
		return
	}

	vars := mux.Vars(r)
	user := vars["user"]
	cl := ginclient.New(serveralias)
	cl.UserToken = session.UserToken

	userrepos, err := cl.ListRepos(user)
	if err != nil {
		log.ShowWrite("[Error] ListRepos failed: %s", err.Error())
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("not found"))
		return
	}

	fmt.Printf("Got %d repos\n", len(userrepos))
	tmpl := template.New("layout")
	tmpl, err = tmpl.Parse(templates.Layout)
	if err != nil {
		log.Write("[Error] failed to parse html layout page")
		fail(w, http.StatusInternalServerError, "something went wrong")
		return
	}
	tmpl, err = tmpl.Parse(templates.RepoList)
	if err != nil {
		log.Write("[Error] failed to render repository list page")
		fail(w, http.StatusInternalServerError, "something went wrong")
		return
	}
	repos := make([]repoHooksInfo, len(userrepos))

	// TODO: check that we have a token configured for each repository with a
	// hook.
	// TODO: if a validator is present but disabled warn the user that a
	// matching hook exists (which means it was created at some point), but it
	// is disabled and offer to enable it. This required differentiating
	// between Enabled, Disabled, and Not Found, so the repohooks map values
	// need to be ternary.

	for idx, rinfo := range userrepos {
		repohooks, err := getRepoHooks(cl, rinfo.FullName)
		if err != nil {
			// simply initialise the map for now
			repohooks = make(map[string]bool)
		}
		repos[idx] = repoHooksInfo{rinfo, repohooks}
	}
	tmpl.Execute(w, &repos)
}

// matchValidator receives a URL path from a GIN hook and returns the validator
// it specifies.
func matchValidator(path string) (string, error) {
	re := regexp.MustCompile(`validate/(?P<validator>[^/]+)/.*`)
	if !re.MatchString(path) {
		return "", fmt.Errorf("URL does not match expected pattern for validator hooks")
	}
	match := re.FindStringSubmatch(path)
	validator := match[1]

	if !helpers.SupportedValidator(validator) {
		return "", fmt.Errorf("URL matches pattern but validator '%s' is not known", validator)
	}

	return validator, nil
}

// getRepoHooks queries the main GIN server and determines which validators are
// enabled via hooks (true), which are configured but disabled (false)
func getRepoHooks(cl *ginclient.Client, repopath string) (map[string]bool, error) {
	// fetch all hooks
	res, err := cl.Get(path.Join("api", "v1", "repos", repopath, "hooks"))
	if err != nil {
		// Bad request?
		log.Write("hook request failed for %s", repopath)
		return nil, fmt.Errorf("hook request failed")
	}
	if res.StatusCode != http.StatusOK {
		// Bad repo path? Unauthorised request?
		log.Write("hook request for %s returned non-OK exit status: %s", repopath, res.Status)
		return nil, fmt.Errorf("hook request returned non-OK exit status: %s", res.Status)
	}
	var ginhooks []gogs.Hook
	defer gweb.CloseRes(res.Body)
	b, err := ioutil.ReadAll(res.Body) // ignore potential read error on res.Body; catch later when trying to unmarshal
	if err != nil {
		// failed to read response body
		log.Write("failed to read response for %s", repopath)
		return nil, fmt.Errorf("failed to read response")
	}
	err = json.Unmarshal(b, &ginhooks)
	if err != nil {
		// failed to parse response body
		log.Write("failed to parse hooks list for %s", repopath)
		return nil, fmt.Errorf("failed to parse hooks list")
	}

	hooks := make(map[string]bool)
	for _, hook := range ginhooks {
		// parse URL to get validator
		hookurl, err := url.Parse(hook.Config["url"])
		if err != nil {
			// can't parse URL. Ignore
			log.Write("can't parse URL %s", hook.Config["url"])
			continue
		}
		validator, err := matchValidator(hookurl.Path)
		if err != nil {
			// Validator not recognised (either path was bad or validator is
			// not supported). Either way, just continue.
			log.Write("validator in path not recognised %s (%s)", hookurl.String(), hookurl.Path)
			log.Write("hook URL in config: %s", hook.Config["url"])
			log.Write(err.Error())
			continue
		}
		// Check if Active, and 'push' is in Events
		var pushenabled bool
		for _, event := range hook.Events {
			if event == "push" {
				pushenabled = true
				break
			}
		}
		if hook.Active && pushenabled {
			log.Write("found %s hook for %s", validator, repopath)
			hooks[validator] = true
		} else {
			log.Write("found disabled or invalid %s hook for %s", validator, repopath)
			hooks[validator] = false
		}
		// TODO: Check if the same validator is found twice
	}
	// add supported validators that were not found and mark them disabled
	supportedValidators := config.Read().Settings.Validators
	for _, validator := range supportedValidators {
		if _, ok := hooks[validator]; !ok {
			hooks[validator] = false
		}
	}
	return hooks, nil
}

// ShowRepo renders the repository information page where the user can enable or
// disable validator hooks.
func ShowRepo(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		return
	}

	session, err := getSessionOrRedirect(w, r)
	if err != nil {
		return
	}

	vars := mux.Vars(r)
	user := vars["user"]
	repo := vars["repo"]
	repopath := fmt.Sprintf("%s/%s", user, repo)
	cl := ginclient.New(serveralias)
	cl.UserToken = session.UserToken
	fmt.Printf("Requesting repository info %s\n", repopath)
	fmt.Printf("Server alias: %s\n", serveralias)
	fmt.Println("Server configuration:")
	fmt.Println(cl.Host)

	repoinfo, err := cl.GetRepo(repopath)
	if err != nil {
		log.ShowWrite("[Error] Repo info failed: %s", err.Error())
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("not found"))
		return
	}

	tmpl := template.New("layout")
	tmpl, err = tmpl.Parse(templates.Layout)
	if err != nil {
		log.Write("[Error] failed to parse html layout page")
		fail(w, http.StatusInternalServerError, "something went wrong")
		return
	}
	tmpl, err = tmpl.Parse(templates.RepoPage)
	if err != nil {
		log.Write("[Error] failed to render repository page")
		fail(w, http.StatusInternalServerError, "something went wrong")
		return
	}

	hooks, err := getRepoHooks(cl, repopath)
	if err != nil {
		hooks = make(map[string]bool)
	}
	repohi := repoHooksInfo{repoinfo, hooks}
	tmpl.Execute(w, &repohi)
}
