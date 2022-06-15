package web

import (
	//"fmt"
	//"io"
	//"log"
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"github.com/G-Node/gin-valid/internal/config"
	//"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHookSecretOK(t *testing.T) {
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
