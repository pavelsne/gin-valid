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
)

func TestValidateOK(t *testing.T) { //TODO
	username := "cervemar"
	reponame := "Testing"
	token := "d1221b5670fad98c590c5540e83e4c4bbf641cbc"
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
}
func TestRepoDoesNotExists(t *testing.T) {
	username := "cervemar"
	reponame := "Testing"
	token := "wtf"
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
}
func TestBadToken(t *testing.T) {
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
func TestUnsupportedValidator(t *testing.T) {
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
func TestHookSecretFailed(t *testing.T) {
	r, _ := http.NewRequest("GET", "wtf", strings.NewReader("{}"))
	srvcfg := config.Read()
	srvcfg.Settings.HookSecret = "hooksecret"
	config.Set(srvcfg)
	r.Header.Add("X-Gogs-Signature", "wtf")
	w := httptest.NewRecorder()
	Validate(w, r)
}
func TestBodyNotJSON(t *testing.T) {
	r, _ := http.NewRequest("GET", "wtf", strings.NewReader("wtf"))
	w := httptest.NewRecorder()
	Validate(w, r)
}

type errReader int

func (errReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("test error")
}
func TestBadBody(t *testing.T) {
	testRequest := httptest.NewRequest(http.MethodPost, "/something", errReader(0))
	w := httptest.NewRecorder()
	Validate(w, testRequest)
}
