package web

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"

	"github.com/G-Node/gin-cli/ginclient"
	"github.com/G-Node/gin-valid/log"
	gogs "github.com/gogits/go-gogs-client"
	"github.com/gorilla/mux"
)

const (
	serveralias = "gin"
	hooksecret  = "omg so sekrit"
)

func EnableHook(w http.ResponseWriter, r *http.Request) {
	fail := func(status int, message string) {
		log.Write("[error] %s", message)
		w.WriteHeader(status)
		w.Write([]byte(message))
	}
	if r.Method != "GET" {
		return
	}
	vars := mux.Vars(r)
	user := vars["user"]
	repo := vars["repo"]
	sessionid, err := r.Cookie("gin-valid-session")
	if err != nil {
		msg := fmt.Sprintf("Hook creation failed: unauthorised")
		fail(http.StatusUnauthorized, msg)
		return
	}

	session, ok := sessions[sessionid.Value]
	if !ok {
		msg := fmt.Sprintf("Hook creation failed: unauthorised")
		fail(http.StatusUnauthorized, msg)
		return
	}
	repopath := fmt.Sprintf("%s/%s", user, repo)
	err = createValidHook(repopath, session)
	if err != nil {
		fail(http.StatusUnauthorized, err.Error())
		return
	}
	msg := fmt.Sprintf("Successfully created hook for %s", repopath)
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(msg))
}

func validateHookSecret(data []byte, secret string) bool {
	sig := hmac.New(sha256.New, []byte(hooksecret))
	sig.Write(data)
	signature := hex.EncodeToString(sig.Sum(nil))
	return signature == secret
}

func createValidHook(repopath string, session *usersession) error {
	log.Write("Adding hook to %s\n", repopath)

	client := ginclient.New(serveralias)
	client.UserToken = session.UserToken
	config := make(map[string]string)
	// TODO: proper host:port
	// TODO: proper secret
	config["url"] = fmt.Sprintf("http://ginvalid:3033/validate/bids/%s", repopath)
	config["content_type"] = "json"
	config["secret"] = hooksecret
	data := gogs.CreateHookOption{
		Type:   "gogs",
		Config: config,
		Active: true,
		Events: []string{"push"},
	}
	res, err := client.Post(fmt.Sprintf("/api/v1/repos/%s/hooks", repopath), data)
	if err != nil {
		log.Write("[error] failed to post: %s\n", err.Error())
		return fmt.Errorf("Hook creation failed: %s", err.Error())
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusCreated {
		log.Write("[error] non-OK response: %s\n", res.Status)
		return fmt.Errorf("Hook creation failed: %s", res.Status)
	}
	hookregs[repopath] = session.UserToken
	return nil
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
