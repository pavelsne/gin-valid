package config

import (
	"os"
	"path/filepath"
)

// Executables used by the server.
type Executables struct {
	Gin  string `json:"gin"`
	BIDS string `json:"bids"`
}

// Directories used by the server for temporary and long term storage.
type Directories struct {
	Temp   string `json:"temp"`
	Result string `json:"result"`
	Log    string `json:"log"`
}

// Denotations provide any freuquently used file names or other denotations
// e.g. validation result files, badge or result folder names.
type Denotations struct {
	LogFile       string `json:"logfile"`
	ResultsFolder string `json:"resultsfolder"`
	ResultsFile   string `json:"resultsfile"`
	ResultsBadge  string `json:"resultsbadge"`
}

// Settings provide the default server settings.
type Settings struct {
	Port         string `json:"port"`
	LogSize      int    `json:"logsize"`
	ResourcesDir string `json:"resourcesdir"`
}

// ServerCfg holds the config used to setup the gin validation server and
// the paths to all required executables, temporary and permanent folders.
type ServerCfg struct {
	Settings Settings    `json:"settings"`
	Exec     Executables `json:"executables"`
	Dir      Directories `json:"directories"`
	Label    Denotations `json:"denotations"`
}

var ginValidDefaultServer = ServerCfg{
	Settings{
		Port:         "3033",
		LogSize:      1048576,
		ResourcesDir: filepath.Join(os.Getenv("GOPATH"), "src", "github.com", "G-Node", "gin-valid", "resources"),
	},
	Executables{
		Gin:  "gin",
		BIDS: "bids-validator",
	},
	Directories{
		Temp:   filepath.Join(os.Getenv("GINVALIDTEMP")),
		Log:    filepath.Join(os.Getenv("GINVALIDHOME"), "gin-valid"),
		Result: filepath.Join(os.Getenv("GINVALIDHOME"), "gin-valid", "results"),
	},
	Denotations{
		LogFile:       "ginvalid.log",
		ResultsFolder: "latest",
		ResultsFile:   "results.json",
		ResultsBadge:  "results.svg",
	},
}

// Read returns the default server configuration.
func Read() ServerCfg {
	return ginValidDefaultServer
}

// Set sets the server configuration.
func Set(cfg ServerCfg) {
	ginValidDefaultServer = cfg
}
