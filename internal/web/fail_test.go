package web

import (
	"github.com/G-Node/gin-valid/internal/resources/templates"
	"net/http/httptest"
	"testing"
)

func TestFailFailedToParse(t *testing.T) {
	w := httptest.NewRecorder()
	original := templates.Layout
	templates.Layout = "{{ WTF? }"
	fail(w, 200, "WTF")
	templates.Layout = original
}
func TestFailFailedToParseFailPage(t *testing.T) {
	w := httptest.NewRecorder()
	original := templates.Fail
	templates.Fail = "{{ WTF? }"
	fail(w, 200, "WTF")
	templates.Fail = original
}
