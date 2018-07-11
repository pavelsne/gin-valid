package log

import (
	"fmt"
	"log"
	"os"
	"path"

	"github.com/mpsonntag/gin-valid/config"
)

var logfile *os.File
var logger *log.Logger

// trim reduces the size of a file to the provided filesize.
// It reads the contents and writes them back, removing the
// initial bytes to fit the limit. If any error occurs, it returns silently.
func trim(file *os.File, filesize int) {
	filestat, err := file.Stat()
	if err != nil {
		return
	}
	if filestat.Size() < int64(filesize) {
		return
	}
	contents := make([]byte, filestat.Size())
	nbytes, err := file.ReadAt(contents, 0)
	if err != nil {
		return
	}
	file.Truncate(0)
	file.Write(contents[nbytes-filesize : nbytes])
}

// Init initialises log file and logger.
func Init() error {
	srvcfg := config.Read()

	fp := path.Join(srvcfg.Dir.Log, srvcfg.Label.LogFile)
	logfile, err := os.OpenFile(fp, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		return err
	}

	trim(logfile, srvcfg.Settings.LogSize)

	logger = log.New(logfile, "", 0)
	Write("\n")

	logger.SetFlags(log.Ldate | log.Ltime)
	Write("=== LOGINIT ===")

	return nil
}

// Write writes a string to the log file if there is an initialized logger.
// Depending on the number of arguments, Write behaves like Print or Printf,
// the first argument must always be a string.
func Write(fmtstr string, args ...interface{}) {
	if logger == nil {
		return
	}
	if len(args) == 0 {
		logger.Print(fmtstr)
	} else {
		logger.Printf(fmtstr, args...)
	}
}

// ShowWrite writes a string to Stdout and passes
// the arguments on to the log Writer function.
func ShowWrite(fmtstr string, args ...interface{}) {
	if len(args) == 0 {
		fmt.Fprint(os.Stdout, fmtstr)
	} else {
		fmt.Fprintf(os.Stdout, fmtstr, args...)
	}
	Write(fmtstr, args...)
}

// Close trims and closes the log file, errors are ignored.
func Close() {
	srvcfg := config.Read()

	Write("=== LOGEND ===")
	trim(logfile, srvcfg.Settings.LogSize)

	_ = logfile.Close()
}
