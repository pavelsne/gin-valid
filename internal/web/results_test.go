package web

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/G-Node/gin-valid/internal/config"
	"github.com/gorilla/mux"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func TestResultsSomeResults(t *testing.T) {
	username := "cervemar"
	reponame := "Testing"
	id := "1"
	content := []byte("wtf")
	body := []byte("{}")
	router := mux.NewRouter()
	router.HandleFunc("/results/{validator}/{user}/{repo}/{id}", Results).Methods("GET")
	r, _ := http.NewRequest("GET", filepath.Join("/results/bids", username, "/", reponame, "/", id), bytes.NewReader(body))
	w := httptest.NewRecorder()
	srvcfg := config.Read()
	fmt.Printf("%+v\n", os.Mkdir("bids", 0755))
	os.Mkdir(filepath.Join(srvcfg.Dir.Result, "bids", username), 0755)
	os.Mkdir(filepath.Join(srvcfg.Dir.Result, "bids", username, reponame), 0755)
	os.Mkdir(filepath.Join(srvcfg.Dir.Result, "bids", username, reponame, id), 0755)
	os.WriteFile(filepath.Join(srvcfg.Dir.Result, "bids", username, reponame, id, srvcfg.Label.ResultsFile), content, 0644)
	sig := hmac.New(sha256.New, []byte(srvcfg.Settings.HookSecret))
	sig.Write(body)
	r.Header.Add("X-Gogs-Signature", hex.EncodeToString(sig.Sum(nil)))
	router.ServeHTTP(w, r)
}
func TestResultsNoResults(t *testing.T) {
	body := []byte("{}")
	router := mux.NewRouter()
	router.HandleFunc("/results/{validator}/{user}/{repo}/{id}", Results).Methods("GET")
	r, _ := http.NewRequest("GET", "/results/bids/whatever/whatever/whatever", bytes.NewReader(body))
	w := httptest.NewRecorder()
	srvcfg := config.Read()
	sig := hmac.New(sha256.New, []byte(srvcfg.Settings.HookSecret))
	sig.Write(body)
	r.Header.Add("X-Gogs-Signature", hex.EncodeToString(sig.Sum(nil)))
	router.ServeHTTP(w, r)
}
func TestResultsUnsupportedValidator(t *testing.T) {
	body := []byte("{}")
	router := mux.NewRouter()
	router.HandleFunc("/results/{validator}/{user}/{repo}/{id}", Results).Methods("GET")
	r, _ := http.NewRequest("GET", "/results/wtf/whatever/whatever/whatever", bytes.NewReader(body))
	w := httptest.NewRecorder()
	srvcfg := config.Read()
	sig := hmac.New(sha256.New, []byte(srvcfg.Settings.HookSecret))
	sig.Write(body)
	r.Header.Add("X-Gogs-Signature", hex.EncodeToString(sig.Sum(nil)))
	router.ServeHTTP(w, r)
}
func TestResultsIDNotSpecified(t *testing.T) {
	body := []byte("{}")
	router := mux.NewRouter()
	router.HandleFunc("/results/{validator}/{user}/{repo}/", Results).Methods("GET")
	r, _ := http.NewRequest("GET", "/results/bids/whatever/whatever/", bytes.NewReader(body))
	w := httptest.NewRecorder()
	srvcfg := config.Read()
	sig := hmac.New(sha256.New, []byte(srvcfg.Settings.HookSecret))
	sig.Write(body)
	r.Header.Add("X-Gogs-Signature", hex.EncodeToString(sig.Sum(nil)))
	router.ServeHTTP(w, r)
}
