package web

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	//"encoding/json"
	"errors"
	//"fmt"
	gweb "github.com/G-Node/gin-cli/web"
	"github.com/G-Node/gin-valid/internal/config"
	"github.com/gorilla/mux"
	//"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

var username = "valid-testing"
var reponame = "Testing"
var token = "4c82d07cccf103e071ad9ee8aec82c34d7003c6c"

func TestValidateBadgeFail(t *testing.T) { //TODO
	body := []byte("{}")
	router := mux.NewRouter()
	router.HandleFunc("/validate/{validator}/{user}/{repo}", Validate).Methods("POST")
	srvcfg := config.Read()
	srvcfg.Dir.Tokens = "."
	os.Mkdir("tmp", 0755)
	srvcfg.Dir.Temp = "./tmp"
	srvcfg.GINAddresses.WebURL = "https://gin.dev.g-node.org:443"
	srvcfg.GINAddresses.GitURL = "git@gin.dev.g-node.org:22"
	config.Set(srvcfg)
	var tok gweb.UserToken
	tok.Username = username
	tok.Token = token
	saveToken(tok)
	os.Mkdir(filepath.Join(srvcfg.Dir.Tokens, "by-repo"), 0755)
	linkToRepo(username, filepath.Join(username, "/", reponame))
	r, _ := http.NewRequest("POST", filepath.Join("/validate/bids/", username, "/", reponame), bytes.NewReader(body))
	w := httptest.NewRecorder()
	sig := hmac.New(sha256.New, []byte(srvcfg.Settings.HookSecret))
	sig.Write(body)
	r.Header.Add("X-Gogs-Signature", hex.EncodeToString(sig.Sum(nil)))
	os.Mkdir("tmp", 0755)
	router.ServeHTTP(w, r)
	time.Sleep(5 * time.Second) //TODO HACK
	os.RemoveAll(filepath.Join(srvcfg.Dir.Tokens, "by-repo"))
}
func TestValidateTMPFail(t *testing.T) {
	body := []byte("{}")
	router := mux.NewRouter()
	router.HandleFunc("/validate/{validator}/{user}/{repo}", Validate).Methods("POST")
	srvcfg := config.Read()
	srvcfg.Dir.Tokens = "."
	config.Set(srvcfg)
	var tok gweb.UserToken
	tok.Username = username
	tok.Token = token
	saveToken(tok)
	os.Mkdir(filepath.Join(srvcfg.Dir.Tokens, "by-repo"), 0755)
	linkToRepo(username, filepath.Join(username, "/", reponame))
	r, _ := http.NewRequest("POST", filepath.Join("/validate/bids/", username, "/", reponame), bytes.NewReader(body))
	w := httptest.NewRecorder()
	sig := hmac.New(sha256.New, []byte(srvcfg.Settings.HookSecret))
	sig.Write(body)
	r.Header.Add("X-Gogs-Signature", hex.EncodeToString(sig.Sum(nil)))
	router.ServeHTTP(w, r)
	time.Sleep(5 * time.Second) //TODO HACK
}
func TestValidateRepoDoesNotExists(t *testing.T) {
	token2 := "wtf"
	body := []byte("{}")
	router := mux.NewRouter()
	router.HandleFunc("/validate/{validator}/{user}/{repo}", Validate).Methods("POST")
	srvcfg := config.Read()
	srvcfg.Dir.Tokens = "."
	config.Set(srvcfg)
	var tok gweb.UserToken
	tok.Username = username
	tok.Token = token2
	saveToken(tok)
	os.Mkdir(filepath.Join(srvcfg.Dir.Tokens, "by-repo"), 0755)
	linkToRepo(username, filepath.Join(username, "/", reponame))
	r, _ := http.NewRequest("POST", filepath.Join("/validate/bids/", username, "/", reponame), bytes.NewReader(body))
	w := httptest.NewRecorder()
	sig := hmac.New(sha256.New, []byte(srvcfg.Settings.HookSecret))
	sig.Write(body)
	r.Header.Add("X-Gogs-Signature", hex.EncodeToString(sig.Sum(nil)))
	router.ServeHTTP(w, r)
	time.Sleep(5 * time.Second) //TODO HACK
}
func TestValidateBadToken(t *testing.T) {
	body := []byte("{}")
	router := mux.NewRouter()
	router.HandleFunc("/validate/{validator}/{user}/{repo}", Validate).Methods("POST")
	r, _ := http.NewRequest("POST", "/validate/bids/whatever/whatever", bytes.NewReader(body))
	w := httptest.NewRecorder()
	srvcfg := config.Read()
	sig := hmac.New(sha256.New, []byte(srvcfg.Settings.HookSecret))
	sig.Write(body)
	r.Header.Add("X-Gogs-Signature", hex.EncodeToString(sig.Sum(nil)))
	router.ServeHTTP(w, r)
}
func TestValidateUnsupportedValidator(t *testing.T) {
	body := []byte("{}")
	r, _ := http.NewRequest("GET", "wtf", bytes.NewReader(body))
	srvcfg := config.Read()
	srvcfg.Settings.HookSecret = "hooksecret"
	config.Set(srvcfg)
	sig := hmac.New(sha256.New, []byte(srvcfg.Settings.HookSecret))
	sig.Write(body)
	r.Header.Add("X-Gogs-Signature", hex.EncodeToString(sig.Sum(nil)))
	w := httptest.NewRecorder()
	Validate(w, r)
}
func TestValidateHookSecretFailed(t *testing.T) {
	r, _ := http.NewRequest("GET", "wtf", strings.NewReader("{}"))
	srvcfg := config.Read()
	srvcfg.Settings.HookSecret = "hooksecret"
	config.Set(srvcfg)
	r.Header.Add("X-Gogs-Signature", "wtf")
	w := httptest.NewRecorder()
	Validate(w, r)
}
func TestValidateBodyNotJSON(t *testing.T) {
	r, _ := http.NewRequest("GET", "wtf", strings.NewReader("wtf"))
	w := httptest.NewRecorder()
	Validate(w, r)
}

type errReader int

func (errReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("test error")
}
func TestValidateBadBody(t *testing.T) {
	testRequest := httptest.NewRequest(http.MethodPost, "/something", errReader(0))
	w := httptest.NewRecorder()
	Validate(w, testRequest)
}
