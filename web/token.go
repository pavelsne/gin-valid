package web

import (
	"encoding/gob"
	"os"
	"path/filepath"
	"strings"

	gweb "github.com/G-Node/gin-cli/web"
	"github.com/G-Node/gin-valid/config"
)

// saveToken writes a token to disk using the username as filename.
// The location is defined by config.Dir.Tokens.
func saveToken(ut gweb.UserToken) error {
	cfg := config.Read()
	filename := filepath.Join(cfg.Dir.Tokens, ut.Username)
	tokenfile, err := os.Create(filename)
	defer tokenfile.Close()
	if err != nil {
		return err
	}
	encoder := gob.NewEncoder(tokenfile)
	return encoder.Encode(ut)
}

// loadUserToken reads a token from disk using the username as filename.
// The location is defined by config.Dir.Tokens.
func loadUserToken(username string) (gweb.UserToken, error) {
	cfg := config.Read()
	filename := filepath.Join(cfg.Dir.Tokens, username)
	return loadToken(filename)
}

// loadToken loads a token from the provided path
func loadToken(path string) (gweb.UserToken, error) {
	ut := gweb.UserToken{}
	tokenfile, err := os.Open(path)
	if err != nil {
		return ut, err
	}
	defer tokenfile.Close()

	decoder := gob.NewDecoder(tokenfile)
	err = decoder.Decode(ut)
	return ut, err
}

// linkToSession links a sessionID to a user's token.
func linkToSession(username string, sessionid string) error {
	cfg := config.Read()
	tokendir := cfg.Dir.Tokens
	utfile := filepath.Join(tokendir, username)
	sidfile := filepath.Join(tokendir, "by-sessionid", sessionid)
	return os.Link(utfile, sidfile)
}

// getTokenBySession loads a user's access token using the session ID found in
// the user's cookie store.
func getTokenBySession(sessionid string) (gweb.UserToken, error) {
	cfg := config.Read()
	tokendir := cfg.Dir.Tokens
	filename := filepath.Join(tokendir, "by-sessionid", sessionid)
	return loadToken(filename)
}

// linkToRepo links a repository name to a user's token.
// This token will be used for cloning a repository to run a validator when a
// web hook is triggered.
func linkToRepo(username string, repopath string) error {
	cfg := config.Read()
	tokendir := cfg.Dir.Tokens
	utfile := filepath.Join(tokendir, username)
	sidfile := filepath.Join(tokendir, "by-repo", strings.Replace(repopath, "/", "-", -1))
	return os.Link(utfile, sidfile)
}

// getTokenByRepo loads a user's access token using a repository path.
func getTokenByRepo(repopath string) (gweb.UserToken, error) {
	cfg := config.Read()
	tokendir := cfg.Dir.Tokens
	filename := filepath.Join(tokendir, "by-repo", strings.Replace(repopath, "/", "-", -1))
	return loadToken(filename)
}
