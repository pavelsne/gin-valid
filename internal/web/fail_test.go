package web

import (
	"github.com/G-Node/gin-valid/internal/resources/templates"
	"net/http/httptest"
	"testing"
)

func TestFailFailedToParse(t *testing.T) {
	w := httptest.NewRecorder()
	templates.Layout = "{{ WTF? }"
	fail(w, 200, "WTF")
}
func TestFailFailedToParseFailPage(t *testing.T) {
	w := httptest.NewRecorder()
	templates.Fail = "{{ WTF? }"
	fail(w, 200, "WTF")
}
