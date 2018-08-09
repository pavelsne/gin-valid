package web

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/G-Node/gin-cli/ginclient"
	glog "github.com/G-Node/gin-cli/ginclient/log"
	"github.com/G-Node/gin-cli/web"
	"github.com/G-Node/gin-valid/log"
	"github.com/G-Node/gin-valid/resources/templates"
	gogs "github.com/gogits/go-gogs-client"
	"github.com/gorilla/mux"
)

var sessions = make(map[string]web.UserToken)

func doLogin(username, password string) error {
	// TODO: remove this function when it becomes a standalone function in gin-cli
	// see https://github.com/G-Node/gin-cli/issues/212
	clientID := "gin-valid"
	gincl := ginclient.New(serveralias)
	glog.Init("")
	glog.Write("Performing login from gin-valid")
	err := gincl.Login(username, password, "gin-valid")
	if err != nil {
		return err
	}
	tokenCreate := &gogs.CreateAccessTokenOption{Name: clientID}
	address := fmt.Sprintf("/api/v1/users/%s/tokens", username)
	res, err := gincl.PostBasicAuth(address, username, password, tokenCreate)
	if err != nil {
		return err // return error from PostBasicAuth directly
	}
	if res.StatusCode != http.StatusCreated {
		return fmt.Errorf(res.Status)
	}
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	log.Write("Got response: %s", res.Status)
	token := ginclient.AccessToken{}
	err = json.Unmarshal(data, &token)
	if err != nil {
		return err
	}
	gincl.Username = username
	gincl.Token = token.Sha1
	log.Write("Login successful. Username: %s", username)

	sessions[username] = gincl.UserToken
	return nil
}

// Login renders the login form and logs in the user to the GIN server, storing
// a session token and key.
func Login(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet {
		log.Write("Login page")
		tmpl := template.New("layout")
		tmpl, err := tmpl.Parse(templates.Layout)
		if err != nil {
			log.Write("[Error] failed to parse html layout page")
			http.ServeContent(w, r, "unavailable", time.Now(), bytes.NewReader([]byte("500 Something went wrong...")))
			return
		}
		tmpl, err = tmpl.Parse(templates.Login)
		if err != nil {
			log.Write("[Error] failed to render login page")
			http.ServeContent(w, r, "unavailable", time.Now(), bytes.NewReader([]byte("500 Something went wrong...")))
			return
		}
		tmpl.Execute(w, nil)
	} else if r.Method == http.MethodPost {
		log.Write("Doing login")
		r.ParseForm()
		username := r.Form["username"][0]
		password := r.Form["password"][0]
		err := doLogin(username, password)
		if err != nil {
			log.Write("[error] Login failed: %s", err.Error())
			http.ServeContent(w, r, "auth failed", time.Now(), bytes.NewReader([]byte("auth failed")))
			return
		}

		// Redirect to repo listing
		http.Redirect(w, r, fmt.Sprintf("/repos/%s", username), http.StatusFound)
	}
}

const repostmpl = `
	{{ define "content" }}
	<br/><br/>
	<div>
	{{ range . }}
	<div><b><a href=/repos/{{.FullName}}/enable>{{.FullName}}</a></b></div>
	<div>{{.Description}} {{.Website}}</div>
	<hr>
	{{ end }}
	</div>
	{{ end }}
`

// ListRepos queries the GIN server for a list of repositories owned (or
// accessible) by a given user and renders the page which displays the
// repositories and their validation status.
func ListRepos(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		return
	}
	vars := mux.Vars(r)
	user := vars["user"]
	cl := ginclient.New(serveralias)
	cl.LoadToken()
	fmt.Printf("Requesting repository listing for user %s\n", user)
	fmt.Printf("Server alias: %s\n", serveralias)
	fmt.Println("Server configuration:")
	fmt.Println(cl.Host)

	repos, err := cl.ListRepos(user)
	if err != nil {
		errmsg := fmt.Sprintf("404 %s", err.Error())
		log.ShowWrite(err.Error())
		w.WriteHeader(http.StatusNotFound)
		http.ServeContent(w, r, "not found", time.Now(), bytes.NewReader([]byte(errmsg)))
		return
	}

	fmt.Printf("Got %d repos\n", len(repos))
	tmpl := template.New("layout")
	tmpl, err = tmpl.Parse(templates.Layout)
	if err != nil {
		log.Write("[Error] failed to parse html layout page")
		http.ServeContent(w, r, "unavailable", time.Now(), bytes.NewReader([]byte("500 Something went wrong...")))
		return
	}
	tmpl, err = tmpl.Parse(repostmpl)
	if err != nil {
		log.Write("[Error] failed to render login page")
		http.ServeContent(w, r, "unavailable", time.Now(), bytes.NewReader([]byte("500 Something went wrong...")))
		return
	}
	tmpl.Execute(w, &repos)
}
