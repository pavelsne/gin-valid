package web

import (
	//"encoding/json"
	//"github.com/G-Node/gin-valid/internal/config"
	//"io/ioutil"
	"testing"
)

var password = "student"

func TestUserCookieExp(t *testing.T) {
	cookieExp()
}
func TestUserDoLoginFailed(t *testing.T) {
	doLogin("wtf", "wtf")
}
func TestUserDoLoginOK(t *testing.T) {
	doLogin(username, password)
}
