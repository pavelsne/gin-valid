package web

import (
	"testing"
)

func TestTokenLinkToSessionWrong(t *testing.T) {
	err := linkToSession("wtf", "wtf")
	if err == nil {
		t.Fatalf(`linkToSession(username string, sessionid string) = %v`, err)
	}
}

func TestTokenGetTokenBySessionWrong(t *testing.T) {
	_, err := getTokenBySession("wtf")
	if err == nil {
		t.Fatalf(`getTokenBySession(sessionid string) = %v`, err)
	}
}

func TestTokenRmTokenRepoLinkWrong(t *testing.T) {
	err := rmTokenRepoLink("wtf")
	if err == nil {
		t.Fatalf(`rmTokenRepoLink(repopath string) = %v`, err)
	}
}

func TestTokenGetTokenByUsernameWrong(t *testing.T) {
	_, err := getTokenByUsername("wtf")
	if err == nil {
		t.Fatalf(`getTokenByUsername(username string) = %v`, err)
	}
}
