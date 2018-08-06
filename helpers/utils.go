package helpers

import (
	"os"

	"github.com/G-Node/gin-valid/config"
	"github.com/G-Node/gin-valid/log"
)

// ValidDirectory checks whether a given path exists and refers to a valid directory.
func ValidDirectory(path string) bool {
	var fi os.FileInfo
	var err error
	if fi, err = os.Stat(path); err != nil {
		log.ShowWrite("[Error] checking directory '%s': '%s'\n", path, err.Error())
		return false
	} else if !fi.IsDir() {
		log.ShowWrite("[Error] invalid directory '%s'\n", fi.Name())
		return false
	}
	return true
}

// SupportedValidator checks whether a string matches
// Validators supported by the server.
func SupportedValidator(validator string) bool {
	validators := config.Read().Settings.Validators

	for _, val := range validators {
		if val == validator {
			return true
		}
	}
	return false
}
