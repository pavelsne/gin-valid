package web

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/G-Node/gin-cli/ginclient"
	gogs "github.com/gogits/go-gogs-client"
)

const serveralias = "gin" // change to "gin" for live server

func createValidHook(repopath string) {
	fmt.Printf("Adding hook to %s\n", repopath)

	client := ginclient.New(serveralias)
	err := client.LoadToken()
	if err != nil {
		fmt.Fprintf(os.Stderr, "[error] failed to load token %s\n", err.Error())
		log.Fatal()
	}
	fmt.Printf("Loaded token %s\n", client.Token)
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
		log.Fatal(fmt.Sprintf("[error] bad response from server %s\n", err.Error()))
	}
	defer res.Body.Close()
	fmt.Printf("Got response: %s\n", res.Status)
	bdy, _ := ioutil.ReadAll(res.Body)
	fmt.Println(string(bdy))
}

func deleteValidHook(repopath string, id int) {
	fmt.Printf("Deleting %d from %s\n", id, repopath)

	client := ginclient.New(serveralias)
	err := client.LoadToken()
	if err != nil {
		fmt.Fprintf(os.Stderr, "[error] failed to load token %s\n", err.Error())
		log.Fatal()
	}
	fmt.Printf("Loaded token %s\n", client.Token)
	res, err := client.Delete(fmt.Sprintf("/api/v1/repos/%s/hooks/%d", repopath, id))
	if err != nil {
		log.Fatal(fmt.Sprintf("[error] bad response from server %s\n", err.Error()))
	}
	defer res.Body.Close()
	fmt.Printf("Got response: %s\n", res.Status)
	bdy, _ := ioutil.ReadAll(res.Body)
	fmt.Println(string(bdy))
}
