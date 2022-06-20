package web

import (
	"bytes"
	"github.com/G-Node/gin-valid/internal/config"
	"github.com/G-Node/gin-valid/internal/resources/templates"
	"github.com/gorilla/mux"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

var password = "student"
var weburl = "https://gin.dev.g-node.org:443"
var giturl = "git@gin.dev.g-node.org:22"

func TestUserCookieExp(t *testing.T) {
	cookieExp()
}
func TestUserDoLoginFailed(t *testing.T) {
	srvcfg := config.Read()
	srvcfg.GINAddresses.WebURL = weburl
	srvcfg.GINAddresses.GitURL = giturl
	config.Set(srvcfg)
	doLogin("wtf", "wtf")
}
func TestUserDoLoginOK(t *testing.T) {
	srvcfg := config.Read()
	srvcfg.GINAddresses.WebURL = weburl
	srvcfg.GINAddresses.GitURL = giturl
	config.Set(srvcfg)
	pth, _ := filepath.Abs(srvcfg.Dir.Tokens)
	tokendir := filepath.Join(pth, "by-sessionid")
	os.MkdirAll(tokendir, 0755)
	doLogin(username, password)
	os.RemoveAll(tokendir)
}
func TestUserLoginPost(t *testing.T) {
	data := url.Values{}
	data.Set("username", username)
	data.Set("password", password)
	router := mux.NewRouter()
	router.HandleFunc("/login", LoginPost).Methods("POST")
	r, _ := http.NewRequest("POST", "/login", strings.NewReader(data.Encode()))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, r)
}
func TestUserLoginGet(t *testing.T) {
	body := []byte("{}")
	router := mux.NewRouter()
	router.HandleFunc("/login", LoginGet).Methods("GET")
	r, _ := http.NewRequest("GET", "/login", bytes.NewReader(body))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, r)
}
func TestUserLoginGetBadLoginPage(t *testing.T) {
	original := templates.Login
	templates.Login = "{{ WTF? }"
	body := []byte("{}")
	router := mux.NewRouter()
	router.HandleFunc("/login", LoginGet).Methods("GET")
	r, _ := http.NewRequest("GET", "/login", bytes.NewReader(body))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, r)
	templates.Login = original
}
func TestUserLoginGetBadLayout(t *testing.T) {
	original := templates.Layout
	templates.Layout = "{{ WTF? }"
	body := []byte("{}")
	router := mux.NewRouter()
	router.HandleFunc("/login", LoginGet).Methods("GET")
	r, _ := http.NewRequest("GET", "/login", bytes.NewReader(body))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, r)
	templates.Layout = original
}
