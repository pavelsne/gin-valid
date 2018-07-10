package log

import (
	"log"
	"os"
	"path"

	"github.com/mpsonntag/gin-valid/config"
)

var logfile *os.File
var logger *log.Logger

// Init initialises log file and logger.
func Init() error {
	srvcfg := config.Read()

	fp := path.Join(srvcfg.Dir.Log, srvcfg.Label.LogFile)
	logfile, err := os.OpenFile(fp, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		return err
	}

	logger = log.New(logfile, "", log.Ldate | log.Ltime | log.Lshortfile)

	return nil
}
