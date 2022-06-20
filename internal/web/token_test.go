package web

import (
	//"bytes"
	//"crypto/hmac"
	//"crypto/sha256"
	//"encoding/hex"
	//"github.com/G-Node/gin-valid/internal/config"
	//"github.com/gorilla/mux"
	//"net/http"
	//"net/http/httptest"
	"testing"
)

func TestTokenLinkToSessionWrong(t *testing.T) {
	linkToSession("wtf", "wtf")
}
func TestTokenGetTokenBySessionWrong(t *testing.T) {
	getTokenBySession("wtf")
}
func TestTokenRmTokenRepoLinkWrong(t *testing.T) {
	rmTokenRepoLink("wtf")
}
func TestTokenGetTokenByUsernameWrong(t *testing.T) {
	getTokenByUsername("wtf")
}
