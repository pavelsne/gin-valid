package web

import (
	"html/template"
	"net/http"

	"github.com/G-Node/gin-valid/log"
	"github.com/G-Node/gin-valid/resources/templates"
)

// fail logs an error and renders an error page with the given message,
// returning the given status code to the user.
func fail(w http.ResponseWriter, status int, message string) {
	log.Write("[error] %s", message)
	w.WriteHeader(status)

	tmpl := template.New("layout")
	tmpl, err := tmpl.Parse(templates.Layout)
	if err != nil {
		log.Write("[Error] failed to parse html layout page. Displaying error message without layout.")
		tmpl = template.New("content")
	}
	tmpl, err = tmpl.Parse(templates.Fail)
	if err != nil {
		log.Write("[Error] failed to render fail page. Displaying plain error message.")
		w.Write([]byte(message))
		return
	}
	errinfo := struct {
		StatusCode int
		StatusText string
		Message    string
	}{
		status,
		http.StatusText(status),
		message,
	}
	tmpl.Execute(w, &errinfo)
}
