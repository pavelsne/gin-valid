package web

import (
	//"encoding/json"
	"github.com/G-Node/gin-valid/internal/config"
	//"io/ioutil"
	"testing"
)

var password = "student"
var weburl = "https://gin.dev.g-node.org:443"
var giturl = "git@gin.dev.g-node.org:22"

func TestUserCookieExp(t *testing.T) {
	cookieExp()
}
func TestUserDoLoginFailed(t *testing.T) {
	srvcfg := config.Read()
	srvcfg.GINAddresses.WebURL = weburl
	srvcfg.GINAddresses.GitURL = giturl
	config.Set(srvcfg)
	doLogin("wtf", "wtf")
}
func TestUserDoLoginOK(t *testing.T) {
	srvcfg := config.Read()
	srvcfg.GINAddresses.WebURL = weburl
	srvcfg.GINAddresses.GitURL = giturl
	config.Set(srvcfg)
	doLogin(username, password)
}
