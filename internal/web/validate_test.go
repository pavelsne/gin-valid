package web

import (
	//"fmt"
	//"io"
	//"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestFirst(t *testing.T) {
	r, _ := http.NewRequest("GET", "wtf", strings.NewReader("{}"))
	r.Header.Add("X-Gogs-Signature", "hooksecret")
	w := httptest.NewRecorder()
	Validate(w, r)
}
