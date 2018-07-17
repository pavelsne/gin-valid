package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/G-Node/gin-valid/config"
	"github.com/G-Node/gin-valid/log"
	"github.com/G-Node/gin-valid/valutils"
	"github.com/G-Node/gin-valid/web"
	"github.com/docopt/docopt-go"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

const usage = `Server validating BIDS files

Usage:
  ginvalid [--listen <port>]
  ginvalid -h | --help
  ginvalid --version

Options:
  -h --help           Show this screen.
  --version           Print version.
  --listen            Port to listen at [default:3033]
  `

func root(w http.ResponseWriter, r *http.Request) {
	http.ServeContent(w, r, "root", time.Now(), bytes.NewReader([]byte("alive")))
}

func registerRoutes(r *mux.Router) {
	r.HandleFunc("/", root)
	r.HandleFunc("/validate/{user}/{repo}", web.Validate)
	r.HandleFunc("/status/{user}/{repo}", web.Status)
	r.HandleFunc("/results/{user}/{repo}", web.Results)
}

func main() {

	// Initialize and read the default server config
	srvcfg := config.Read()

	// Parse commandline arguments
	args, err := docopt.ParseArgs(usage, nil, "v1.0.0")
	if err != nil {
		fmt.Fprintf(os.Stderr, "\n[Error] parsing cli arguments: '%s', abort...\n\n", err.Error())
		os.Exit(-1)
	}

	// Parse and load custom server confguration
	if args["--config"] != nil {
		content, err := ioutil.ReadFile(args["--config"].(string))
		if err != nil {
			fmt.Fprintf(os.Stderr, "[Error] reading config file %v\n", args["--config"])
			os.Exit(-1)
		}

		// Overwrite any default settings with information from the
		// provided config file.
		err = json.Unmarshal(content, &srvcfg)
		if err != nil {
			fmt.Fprintf(os.Stderr, "[Error] unmarshalling config file %s\n", err.Error())
			os.Exit(-1)
		}
		config.Set(srvcfg)
	}

	err = log.Init()
	if err != nil {
		fmt.Fprintf(os.Stderr, "[Error] opening logfile '%s'\n", err.Error())
		os.Exit(-1)
	}
	defer log.Close()

	// Log cli arguments
	log.Write("[Warmup] cli arguments: %v\n", args)

	// Check whether the required directories are available and accessible
	if !valutils.ValidDirectory(srvcfg.Dir.Temp) {
		os.Exit(-1)
	}

	log.ShowWrite("[Warmup] using temp directory: '%s'\n", srvcfg.Dir.Temp)

	if !valutils.ValidDirectory(srvcfg.Dir.Result) {
		os.Exit(-1)
	}
	log.ShowWrite("[Warmup] using results directory '%s'\n", srvcfg.Dir.Result)

	// Check gin is installed and available
	outstr, err := valutils.AppVersionCheck(srvcfg.Exec.Gin)
	if err != nil {
		log.ShowWrite("\n[Error] checking gin client '%s'\n", err.Error())
		os.Exit(-1)
	}
	log.ShowWrite("[Warmup] using %s", outstr)

	// Check bids-validator is installed
	outstr, err = valutils.AppVersionCheck(srvcfg.Exec.BIDS)
	if err != nil {
		log.ShowWrite("\n[Error] checking bids-validator '%s'\n", err.Error())
		os.Exit(-1)
	}
	log.ShowWrite("[Warmup] using bids-validator v%s", outstr)

	// Use port if provided.
	port := fmt.Sprintf(":%s", srvcfg.Settings.Port)
	if valutils.IsValidPort(args["<port>"]) {
		p := args["<port>"]
		port = fmt.Sprintf(":%s", p.(string))
	} else {
		log.ShowWrite("[Warning] could not parse a valid port number, using default\n")
	}
	log.ShowWrite("[Warmup] using port: '%s'\n", port)

	log.ShowWrite("[Warmup] registering routes\n")
	router := mux.NewRouter()
	registerRoutes(router)

	handler := handlers.CORS(
		handlers.AllowedHeaders([]string{"Accept", "Content-Type", "Authorization"}),
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowedMethods([]string{"GET"}),
	)(router)

	server := http.Server{
		Addr:    port,
		Handler: handler,
	}

	// Monitor the environment for shutdown signals to
	// gracefully shutdown the server.
	go func() {
		sigchan := make(chan os.Signal, 1)
		signal.Notify(sigchan, os.Interrupt)
		<-sigchan
		log.ShowWrite("[Info] System interrupt, shutting down server\n")
		err := server.Shutdown(context.Background())
		if err != nil {
			log.ShowWrite("[Error] on server shutdown: %v\n", err)
		}
	}()

	log.ShowWrite("[Start] Listen and serve\n")
	err = server.ListenAndServe()
	if err == http.ErrServerClosed {
		log.Close()
		os.Exit(0)
	} else if err != nil {
		log.ShowWrite("[Error] Server startup: '%v', abort...\n\n", err)
		log.Close()
		os.Exit(-1)
	}
}
