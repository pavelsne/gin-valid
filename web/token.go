package web

import (
	"encoding/gob"
	"os"
	"path/filepath"

	gweb "github.com/G-Node/gin-cli/web"
	"github.com/G-Node/gin-valid/config"
)

func saveToken(filename string, ut gweb.UserToken) error {
	cfg := config.Read()
	filename = filepath.Join(cfg.Dir.Tokens, filename)
	tokenfile, err := os.Create(filename)
	defer tokenfile.Close()
	if err != nil {
		return err
	}
	encoder := gob.NewEncoder(tokenfile)
	err = encoder.Encode(ut)
	return err
}

func loadToken(filename string) (gweb.UserToken, error) {
	cfg := config.Read()
	ut := gweb.UserToken{}
	filename = filepath.Join(cfg.Dir.Tokens, filename)
	tokenfile, err := os.Open(filename)
	if err != nil {
		return ut, err
	}
	defer tokenfile.Close()

	decoder := gob.NewDecoder(tokenfile)
	err = decoder.Decode(ut)
	if err != nil {
		return ut, err
	}
	return ut, nil
}
