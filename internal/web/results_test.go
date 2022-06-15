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
	"testing"
)

func TestResults(t *testing.T) {
	body := []byte("{}")
	router := mux.NewRouter()
	router.HandleFunc("/results/{validator}/{user}/{repo}/{id}", Results).Methods("GET")
	r, _ := http.NewRequest("POST", "/results/bids/whatever/whatever/whatever", bytes.NewReader(body))
	w := httptest.NewRecorder()
	srvcfg := config.Read()
	sig := hmac.New(sha256.New, []byte(srvcfg.Settings.HookSecret))
	sig.Write(body)
	r.Header.Add("X-Gogs-Signature", hex.EncodeToString(sig.Sum(nil)))
	router.ServeHTTP(w, r)
}
