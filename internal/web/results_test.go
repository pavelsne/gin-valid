package web

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"github.com/G-Node/gin-valid/internal/config"
	"github.com/gorilla/mux"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func TestResultsSomeResults(t *testing.T) {
	id := "1"
	content := "wtf"
	body := []byte("{}")
	router := mux.NewRouter()
	router.HandleFunc("/results/{validator}/{user}/{repo}/{id}", Results).Methods("GET")
	r, _ := http.NewRequest("GET", filepath.Join("/results/bids", username, "/", reponame, "/", id), bytes.NewReader(body))
	w := httptest.NewRecorder()
	srvcfg := config.Read()
	os.MkdirAll(filepath.Join(srvcfg.Dir.Result, "bids", username, reponame, id), 0755)
	f, _ := os.Create(filepath.Join(srvcfg.Dir.Result, "bids", username, reponame, id, srvcfg.Label.ResultsFile))
	defer f.Close()
	f.WriteString(content)
	sig := hmac.New(sha256.New, []byte(srvcfg.Settings.HookSecret))
	sig.Write(body)
	r.Header.Add("X-Gogs-Signature", hex.EncodeToString(sig.Sum(nil)))
	router.ServeHTTP(w, r)
	os.RemoveAll(filepath.Join(srvcfg.Dir.Result, "bids", username, reponame, id))
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
