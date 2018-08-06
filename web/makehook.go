package web

import (
	"fmt"

	"github.com/G-Node/gin-cli/ginclient"
	"github.com/G-Node/gin-valid/log"
	gogs "github.com/gogits/go-gogs-client"
)

const serveralias = "localgogs" // change to "gin" for live server

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
	config["url"] = "https://example.com"
	config["content_type"] = "json"
	data := gogs.CreateHookOption{
		Type:   "gogs",
		Config: config,
		Active: false,
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
