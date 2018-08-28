package web

import (
	"encoding/gob"
	"os"
	"path/filepath"

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

// loadToken reads a token from disk using the username as filename.
// The location is defined by config.Dir.Tokens.
func loadToken(username string) (gweb.UserToken, error) {
	cfg := config.Read()
	ut := gweb.UserToken{}
	filename := filepath.Join(cfg.Dir.Tokens, username)
	tokenfile, err := os.Open(filename)
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
	return nil
}

// getTokenBySession loads a user's access token using the session ID found in
// the user's cookie store.
func getTokenBySession(sessionid string) (gweb.UserToken, error) {
	return gweb.UserToken{}, nil
}

// linkToRepo links a repository name to a user's token.
// This token will be used for cloning a repository to run a validator when a
// web hook is triggered.
func linkToRepo(username string, repopath string) error {
	return nil
}

// getTokenByRepo loads a user's access token using a repository path.
func getTokenByRepo(repopath string) (gweb.UserToken, error) {
	return gweb.UserToken{}, nil
}
