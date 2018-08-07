package web

import (
	"fmt"
	"net/http"

	"github.com/G-Node/gin-cli/ginclient"
	"github.com/G-Node/gin-valid/log"
	gogs "github.com/gogits/go-gogs-client"
	"github.com/gorilla/mux"
)

const serveralias = "gin"

func EnableHook(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		return
	}
	vars := mux.Vars(r)
	user := vars["user"]
	repo := vars["repo"]
	createValidHook(fmt.Sprintf("%s/%s", user, repo))
}

// TODO: Return error
func createValidHook(repopath string) {
	log.Write("Adding hook to %s\n", repopath)

	client := ginclient.New(serveralias)
	err := client.LoadToken()
	if err != nil {
		log.Write("[error] failed to load token %s\n", err.Error())
		return
	}
	config := make(map[string]string)
	// TODO: proper host:port
	// TODO: proper secret
	config["url"] = fmt.Sprintf("http://ginvalid:3033/validate/bids/%s", repopath)
	config["content_type"] = "json"
	config["secret"] = "TODO: Make a proper secret"
	data := gogs.CreateHookOption{
		Type:   "gogs",
		Config: config,
		Active: true,
		Events: []string{"push"},
	}
	res, err := client.Post(fmt.Sprintf("/api/v1/repos/%s/hooks", repopath), data)
	if err != nil {
		log.Write("[error] bad response from server %s\n", err.Error())
		return
	}
	defer res.Body.Close()
	// log.Write("Got response: %s\n", res.Status)
	// bdy, _ := ioutil.ReadAll(res.Body)
	// fmt.Println(string(bdy))
}

func deleteValidHook(repopath string, id int) {
	log.Write("Deleting %d from %s\n", id, repopath)

	client := ginclient.New(serveralias)
	err := client.LoadToken()
	if err != nil {
		log.Write("[error] failed to load token %s\n", err.Error())
		return
	}
	res, err := client.Delete(fmt.Sprintf("/api/v1/repos/%s/hooks/%d", repopath, id))
	if err != nil {
		log.Write("[error] bad response from server %s\n", err.Error())
		return
	}
	defer res.Body.Close()
	// fmt.Printf("Got response: %s\n", res.Status)
	// bdy, _ := ioutil.ReadAll(res.Body)
	// fmt.Println(string(bdy))
}

func CommCheck(user, pass string) error {
	client := ginclient.New(serveralias)
	return client.Login(user, pass, "gin-valid")
}
