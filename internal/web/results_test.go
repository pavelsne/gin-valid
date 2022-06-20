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

func TestResultsUnsupportedV2(t *testing.T) {
	id := "1"
	content := "{\"empty\":\"json\"}"
	body := []byte("{}")
	router := mux.NewRouter()
	router.HandleFunc("/results/{validator}/{user}/{repo}/{id}", Results).Methods("GET")
	r, _ := http.NewRequest("GET", filepath.Join("/results/wtf", username, "/", reponame, "/", id), bytes.NewReader(body))
	w := httptest.NewRecorder()
	srvcfg := config.Read()
	srvcfg.Settings.Validators = append(srvcfg.Settings.Validators, "wtf")
	config.Set(srvcfg)
	os.MkdirAll(filepath.Join(srvcfg.Dir.Result, "nix", username, reponame, id), 0755)
	f, _ := os.Create(filepath.Join(srvcfg.Dir.Result, "nix", username, reponame, id, srvcfg.Label.ResultsFile))
	defer f.Close()
	f.WriteString(content)
	sig := hmac.New(sha256.New, []byte(srvcfg.Settings.HookSecret))
	sig.Write(body)
	r.Header.Add("X-Gogs-Signature", hex.EncodeToString(sig.Sum(nil)))
	router.ServeHTTP(w, r)
	os.RemoveAll(filepath.Join(srvcfg.Dir.Result, "nix", username, reponame, id))
	srvcfg.Settings.Validators = srvcfg.Settings.Validators[:len(srvcfg.Settings.Validators)-1]
	config.Set(srvcfg)
}
func TestResultsODML(t *testing.T) {
	id := "1"
	content := "{\"empty\":\"json\"}"
	body := []byte("{}")
	router := mux.NewRouter()
	router.HandleFunc("/results/{validator}/{user}/{repo}/{id}", Results).Methods("GET")
	r, _ := http.NewRequest("GET", filepath.Join("/results/odml", username, "/", reponame, "/", id), bytes.NewReader(body))
	w := httptest.NewRecorder()
	srvcfg := config.Read()
	os.MkdirAll(filepath.Join(srvcfg.Dir.Result, "odml", username, reponame, id), 0755)
	f, _ := os.Create(filepath.Join(srvcfg.Dir.Result, "odml", username, reponame, id, srvcfg.Label.ResultsFile))
	defer f.Close()
	f.WriteString(content)
	sig := hmac.New(sha256.New, []byte(srvcfg.Settings.HookSecret))
	sig.Write(body)
	r.Header.Add("X-Gogs-Signature", hex.EncodeToString(sig.Sum(nil)))
	router.ServeHTTP(w, r)
	os.RemoveAll(filepath.Join(srvcfg.Dir.Result, "odml", username, reponame, id))
}
func TestResultsNIX(t *testing.T) {
	id := "1"
	content := "{\"empty\":\"json\"}"
	body := []byte("{}")
	router := mux.NewRouter()
	router.HandleFunc("/results/{validator}/{user}/{repo}/{id}", Results).Methods("GET")
	r, _ := http.NewRequest("GET", filepath.Join("/results/nix", username, "/", reponame, "/", id), bytes.NewReader(body))
	w := httptest.NewRecorder()
	srvcfg := config.Read()
	os.MkdirAll(filepath.Join(srvcfg.Dir.Result, "nix", username, reponame, id), 0755)
	f, _ := os.Create(filepath.Join(srvcfg.Dir.Result, "nix", username, reponame, id, srvcfg.Label.ResultsFile))
	defer f.Close()
	f.WriteString(content)
	sig := hmac.New(sha256.New, []byte(srvcfg.Settings.HookSecret))
	sig.Write(body)
	r.Header.Add("X-Gogs-Signature", hex.EncodeToString(sig.Sum(nil)))
	router.ServeHTTP(w, r)
	os.RemoveAll(filepath.Join(srvcfg.Dir.Result, "nix", username, reponame, id))
}
func TestResultsInJSON(t *testing.T) {
	id := "1"
	content := "{\"empty\":\"json\"}"
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
func TestResultsInProgress(t *testing.T) {
	id := "1"
	content := progressmsg
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
