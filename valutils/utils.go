package valutils

import (
	"os"

	"github.com/mpsonntag/gin-valid/log"
)

// ValidDirectory checks whether a given path exists and refers to a valid directory.
func ValidDirectory(path string) bool {
	var fi os.FileInfo
	var err error
	if fi, err = os.Stat(path); err != nil {
		log.Write("[Error] checking temp directory %s\n", err.Error())
		return false
	} else if !fi.IsDir() {
		log.Write("[Error] invalid temp directory '%s' \n", fi.Name())
		return false
	}
	return true
}
