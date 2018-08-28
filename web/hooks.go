package web

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strings"

	"github.com/G-Node/gin-cli/ginclient"
	"github.com/G-Node/gin-valid/config"
	"github.com/G-Node/gin-valid/helpers"
	"github.com/G-Node/gin-valid/log"
	gogs "github.com/gogits/go-gogs-client"
	"github.com/gorilla/mux"
)

func EnableHook(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		return
	}
	vars := mux.Vars(r)
	user := vars["user"]
	repo := vars["repo"]
	validator := strings.ToLower(vars["validator"])
	cfg := config.Read()
	cookiename := cfg.Settings.CookieName
	sessionid, err := r.Cookie(cookiename)
	if err != nil {
		msg := fmt.Sprintf("Hook creation failed: unauthorised")
		fail(w, http.StatusUnauthorized, msg)
		return
	}

	session, ok := sessions[sessionid.Value]
	if !ok {
		msg := fmt.Sprintf("Hook creation failed: unauthorised")
		fail(w, http.StatusUnauthorized, msg)
		return
	}
	if !helpers.SupportedValidator(validator) {
		fail(w, http.StatusNotFound, "unsupported validator")
		return
	}
	repopath := fmt.Sprintf("%s/%s", user, repo)
	err = createValidHook(repopath, validator, session)
	if err != nil {
		// TODO: Check if failure is for other reasons and maybe return 500 instead
		fail(w, http.StatusUnauthorized, err.Error())
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/repos/%s", session.Username), http.StatusFound)
}

func validateHookSecret(data []byte, secret string) bool {
	cfg := config.Read()
	hooksecret := cfg.Settings.HookSecret
	sig := hmac.New(sha256.New, []byte(hooksecret))
	sig.Write(data)
	signature := hex.EncodeToString(sig.Sum(nil))
	return signature == secret
}

func createValidHook(repopath string, validator string, session *usersession) error {
	// TODO: AVOID DUPLICATES:
	//   - If it's already hooked and we have it on record, do nothing
	//   - If it's already hooked, but we don't know about it, check if it's valid and don't recreate
	log.Write("Adding %s hook to %s\n", validator, repopath)

	gvconfig := config.Read()
	client := ginclient.New(serveralias)
	client.UserToken = session.UserToken
	hookconfig := make(map[string]string)
	cfg := config.Read()
	hooksecret := cfg.Settings.HookSecret

	host := fmt.Sprintf("%s:%s", gvconfig.Settings.RootURL, gvconfig.Settings.Port)
	u, err := url.Parse(host)
	u.Path = path.Join(u.Path, "validate", validator, repopath)
	hookconfig["url"] = u.String()
	hookconfig["content_type"] = "json"
	hookconfig["secret"] = hooksecret
	data := gogs.CreateHookOption{
		Type:   "gogs",
		Config: hookconfig,
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
	// store user's token to disk so that the service can use it to clone the
	// repository when needed
	err = saveToken(session.UserToken)
	return err
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
