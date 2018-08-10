package web

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/G-Node/gin-cli/ginclient"
	"github.com/G-Node/gin-cli/ginclient/config"
	glog "github.com/G-Node/gin-cli/ginclient/log"
	"github.com/G-Node/gin-cli/web"
	"github.com/G-Node/gin-valid/log"
	"github.com/G-Node/gin-valid/resources/templates"
	gogs "github.com/gogits/go-gogs-client"
	"github.com/gorilla/mux"
)

type usersession struct {
	sessionID string
	web.UserToken
}

var (
	sessions = make(map[string]*usersession)
	hookregs = make(map[string]web.UserToken)
)

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
	configpath, _ := config.Path(false)
	keyfilepath := filepath.Join(configpath, fmt.Sprintf("%s.key", serveralias))
	os.Remove(keyfilepath)
}

func doLogin(username, password string) (*usersession, error) {
	// TODO: remove this function when it becomes a standalone function in gin-cli
	// see https://github.com/G-Node/gin-cli/issues/212
	clientID := "gin-valid"
	gincl := ginclient.New(serveralias)
	glog.Init("")
	glog.Write("Performing login from gin-valid")
	tokenCreate := &gogs.CreateAccessTokenOption{Name: clientID}
	address := fmt.Sprintf("/api/v1/users/%s/tokens", username)
	res, err := gincl.PostBasicAuth(address, username, password, tokenCreate)
	if err != nil {
		return nil, err // return error from PostBasicAuth directly
	}
	if res.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf(res.Status)
	}
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	log.Write("Got response: %s", res.Status)
	token := ginclient.AccessToken{}
	err = json.Unmarshal(data, &token)
	if err != nil {
		return nil, err
	}
	gincl.Username = username
	gincl.Token = token.Sha1
	log.Write("Login successful. Username: %s", username)

	sessionid := "unique session-id " + username
	return &usersession{sessionid, gincl.UserToken}, nil
}

// Login renders the login form and logs in the user to the GIN server, storing
// a session token and key.
func Login(w http.ResponseWriter, r *http.Request) {
	fail := func(status int, message string) {
		log.Write("[error] %s", message)
		w.WriteHeader(status)
		w.Write([]byte(message))
	}
	if r.Method == http.MethodGet {
		log.Write("Login page")
		tmpl := template.New("layout")
		tmpl, err := tmpl.Parse(templates.Layout)
		if err != nil {
			log.Write("[Error] failed to parse html layout page")
			fail(http.StatusInternalServerError, "something went wrong")
			return
		}
		tmpl, err = tmpl.Parse(templates.Login)
		if err != nil {
			log.Write("[Error] failed to render login page")
			fail(http.StatusInternalServerError, "something went wrong")
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
			fail(http.StatusUnauthorized, "authentication failed")
			return
		}

		sessions[session.sessionID] = session
		cookie := http.Cookie{Name: "gin-valid-session", Value: session.sessionID, Expires: cookieExp()}
		http.SetCookie(w, &cookie)
		// Redirect to repo listing
		http.Redirect(w, r, fmt.Sprintf("/repos/%s", username), http.StatusFound)
	}
}

func getSessionOrRedirect(w http.ResponseWriter, r *http.Request) (*usersession, error) {
	cookie, err := r.Cookie("gin-valid-session")
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
	fail := func(status int, message string) {
		log.Write("[error] %s", message)
		w.WriteHeader(status)
		w.Write([]byte(message))
	}
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
	fmt.Printf("Requesting repository listing for user %s\n", user)
	fmt.Printf("Server alias: %s\n", serveralias)
	fmt.Println("Server configuration:")
	fmt.Println(cl.Host)

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
		fail(http.StatusInternalServerError, "something went wrong")
		return
	}
	tmpl, err = tmpl.Parse(templates.RepoList)
	if err != nil {
		log.Write("[Error] failed to render login page")
		fail(http.StatusInternalServerError, "something went wrong")
		return
	}
	type repoHooksInfo struct {
		gogs.Repository
		Hooks map[string]bool
	}

	repos := make([]repoHooksInfo, len(userrepos))
	// TODO: For each supported hook type, check if it's active
	for idx, r := range userrepos {
		bids := false
		if _, ok := hookregs[r.FullName]; ok {
			bids = true
		}
		repos[idx] = repoHooksInfo{r, map[string]bool{"BIDS": bids, "NIX": false}}
	}
	tmpl.Execute(w, &repos)
}
func Repo(w http.ResponseWriter, r *http.Request) {
	fail := func(status int, message string) {
		log.Write("[error] %s", message)
		w.WriteHeader(status)
		w.Write([]byte(message))
	}
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
		fail(http.StatusInternalServerError, "something went wrong")
		return
	}
	tmpl, err = tmpl.Parse(templates.RepoPage)
	if err != nil {
		log.Write("[Error] failed to render repository page")
		fail(http.StatusInternalServerError, "something went wrong")
		return
	}
	type repoHooksInfo struct {
		gogs.Repository
		Hooks map[string]bool
	}

	bids := false
	if _, ok := hookregs[repoinfo.FullName]; ok {
		bids = true
	}
	repohi := repoHooksInfo{repoinfo, map[string]bool{"BIDS": bids, "NIX": false}}
	tmpl.Execute(w, &repohi)
}
