package config

// Executables used by the server.
type Executables struct {
	Gin  string
	BIDS string
}

// Directories used by the server for temporary and long term storage.
type Directories struct {
	Temp   string
	Result string
	Log    string
}

// Denotations provide any freuquently used file names or other denotations
// e.g. validation result files, badge or result folder names.
type Denotations struct {
	LogFile       string
	ResultsFolder string
	ResultsFile   string
	ResultsBadge  string
}

// Settings provide the default server settings.
type Settings struct {
	Port    string
	LogSize int
}

// ServerCfg holds the config used to setup the gin validation server and
// the paths to all required executables, temporary and permanent folders.
type ServerCfg struct {
	Settings Settings
	Exec     Executables
	Dir      Directories
	Label    Denotations
}

var ginValidDefaultServer = ServerCfg{
	Settings{
		Port:    "3033",
		LogSize: 1048576,
	},
	Executables{
		Gin:  "gin",
		BIDS: "/home/msonntag/node_modules/.bin/bids-validator",
	},
	Directories{
		Temp:   "/home/msonntag/Chaos/DL/val",
		Result: "/home/msonntag/Chaos/DL/valresults",
		Log:    "/home/msonntag/Chaos/DL/val",
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
