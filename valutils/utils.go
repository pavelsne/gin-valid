package valutils

import (
	"fmt"
	"os"
)

// ValidDirectory checks whether a given path exists and refers to a valid directory.
func ValidDirectory(path string) bool {
	var fi os.FileInfo
	var err error
	if fi, err = os.Stat(path); err != nil {
		fmt.Fprintf(os.Stderr, "[Error] checking temp directory %s\n", err.Error())
		return false
	} else if !fi.IsDir() {
		fmt.Fprintf(os.Stderr, "[Error] invalid temp directory '%s' \n", fi.Name())
		return false
	}
	return true
}
