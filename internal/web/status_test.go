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

func TestStatusOK(t *testing.T) {
	content := "wtf"
	body := []byte("{}")
	router := mux.NewRouter()
	router.HandleFunc("/status/{validator}/{user}/{repo}", Status).Methods("GET")
	r, _ := http.NewRequest("GET", filepath.Join("/status/bids", username, reponame), bytes.NewReader(body))
	w := httptest.NewRecorder()
	srvcfg := config.Read()
	sig := hmac.New(sha256.New, []byte(srvcfg.Settings.HookSecret))
	sig.Write(body)
	r.Header.Add("X-Gogs-Signature", hex.EncodeToString(sig.Sum(nil)))
	os.Mkdir(srvcfg.Dir.Result, 0755)
	os.Mkdir(filepath.Join(srvcfg.Dir.Result, "bids"), 0755)
	os.Mkdir(filepath.Join(srvcfg.Dir.Result, "bids", username), 0755)
	os.Mkdir(filepath.Join(srvcfg.Dir.Result, "bids", username, reponame), 0755)
	os.Mkdir(filepath.Join(srvcfg.Dir.Result, "bids", username, reponame, srvcfg.Label.ResultsFolder), 0755)
	f, _ := os.Create(filepath.Join(srvcfg.Dir.Result, "bids", username, reponame, srvcfg.Label.ResultsFolder, srvcfg.Label.ResultsBadge))
	defer f.Close()
	f.WriteString(content)
	router.ServeHTTP(w, r)
}
func TestStatusNoConent(t *testing.T) {
	body := []byte("{}")
	router := mux.NewRouter()
	router.HandleFunc("/status/{validator}/{user}/{repo}", Status).Methods("GET")
	r, _ := http.NewRequest("GET", "/status/bids/whatever/whatever", bytes.NewReader(body))
	w := httptest.NewRecorder()
	srvcfg := config.Read()
	sig := hmac.New(sha256.New, []byte(srvcfg.Settings.HookSecret))
	sig.Write(body)
	r.Header.Add("X-Gogs-Signature", hex.EncodeToString(sig.Sum(nil)))
	router.ServeHTTP(w, r)
}
func TestStatusUnsupportedValidator(t *testing.T) {
	body := []byte("{}")
	router := mux.NewRouter()
	router.HandleFunc("/status/{validator}/{user}/{repo}", Status).Methods("GET")
	r, _ := http.NewRequest("GET", "/status/whatever/whatever/whatever", bytes.NewReader(body))
	w := httptest.NewRecorder()
	srvcfg := config.Read()
	sig := hmac.New(sha256.New, []byte(srvcfg.Settings.HookSecret))
	sig.Write(body)
	r.Header.Add("X-Gogs-Signature", hex.EncodeToString(sig.Sum(nil)))
	router.ServeHTTP(w, r)
}
