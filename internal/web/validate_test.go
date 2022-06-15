package web

import (
	//"fmt"
	//"io"
	//"log"
	"github.com/G-Node/gin-valid/internal/config"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHookSecretOK(t *testing.T) {
	r, _ := http.NewRequest("GET", "wtf", strings.NewReader("{}"))
	srvcfg := config.Read()
	srvcfg.Settings.HookSecret = "hooksecret"
	config.Set(srvcfg)
	r.Header.Add("X-Gogs-Signature", "hooksecret")
	w := httptest.NewRecorder()
	Validate(w, r)
}
func TestHokkSecretFailed(t *testing.T) {
	r, _ := http.NewRequest("GET", "wtf", strings.NewReader("{}"))
	srvcfg := config.Read()
	srvcfg.Settings.HookSecret = "hooksecret"
	config.Set(srvcfg)
	r.Header.Add("X-Gogs-Signature", "wtf")
	w := httptest.NewRecorder()
	Validate(w, r)
}
