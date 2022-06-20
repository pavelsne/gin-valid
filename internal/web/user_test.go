package web

import (
	//"encoding/json"
	"github.com/G-Node/gin-valid/internal/config"
	//"io/ioutil"
	"fmt"
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
	f, _ := filepath.Abs(srvcfg.Dir.Tokens)
	tokendir, _ := filepath.Abs(srvcfg.Dir.Tokens)
	os.Mkdir(f, 0755)
	os.Mkdir(tokendir, 0755)
	os.Mkdir(filepath.Join(tokendir, "by-sessionid"), 0755)
	t2, e := doLogin(username, password)
}
