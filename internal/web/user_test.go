package web

import (
	"github.com/G-Node/gin-valid/internal/config"
	"os"
	"path/filepath"
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
	pth, _ := filepath.Abs(srvcfg.Dir.Tokens)
	tokendir := filepath.Join(pth, "by-sessionid")
	os.MkdirAll(tokendir, 0755)
	doLogin(username, password)
	os.RemoveAll(tokendir)
}
