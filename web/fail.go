package web

import (
	"net/http"

	"github.com/G-Node/gin-valid/log"
)

// fail logs an error and renders an error page with the given message,
// returning the given status code to the user.
func fail(w http.ResponseWriter, status int, message string) {
	log.Write("[error] %s", message)
	w.WriteHeader(status)
	w.Write([]byte(message))
}
