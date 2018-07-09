package config

// Executables used by the server.
type Executables struct {
	Gin string
	NPM string
	BIDS string
}

// Directories used by the server for
// temporary and long term storage.
type Directories struct {
	Temp string
	Permanent string
}

// ServerCfg holds the config used to setup
// the gin validation server and the paths
// to all required executables, temporary
// and permanent folders 
type ServerCfg struct {
	Exec Executables
	Dir Directories
}

var ginValidDefaultServer = ServerCfg{
	Executables{
		Gin: "gin",
		NPM: "npm",
		BIDS: "/home/msonntag/node_modules/.bin/bids-validator",
	},
	Directories{
		Temp: "/home/msonntag/Chaos/DL/val",
		Permanent: "/home/msonntag/node_modules/.bin/bids-validator",
	},
}

// Read returns the default server configuration.
func Read() ServerCfg {
	return ginValidDefaultServer
}
