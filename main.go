package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/docopt/docopt-go"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/mpsonntag/gin-valid/config"
	"github.com/mpsonntag/gin-valid/valutils"
	"github.com/mpsonntag/gin-valid/web"
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
	fmt.Fprintln(w, "ginvalid running")
}

func registerRoutes(r *mux.Router) {
	r.HandleFunc("/", root)
	r.HandleFunc("/validate/{user}/{repo}", web.Validate)
}

func main() {

	srvconfig := config.Read()

	// Check whether the required directories are available and accessible
	if !valutils.ValidDirectory(srvconfig.Dir.Temp) {
		os.Exit(-1)
	}
	fmt.Fprintf(os.Stdout, "[Warmup] using temp directory: '%s'\n", srvconfig.Dir.Temp)

	if !valutils.ValidDirectory(srvconfig.Dir.Result) {
		os.Exit(-1)
	}
	fmt.Fprintf(os.Stdout, "[Warmup] using results directory '%s'\n", srvconfig.Dir.Result)

	// Check gin is installed and available
	outstr, err := valutils.AppVersionCheck(srvconfig.Exec.Gin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "\n[Error] checking gin client '%s'\n", err.Error())
		os.Exit(-1)
	}
	fmt.Fprintf(os.Stdout, "[Warmup] using %s", outstr)

	// Check npm is installed and available
	outstr, err = valutils.AppVersionCheck(srvconfig.Exec.NPM)
	if err != nil {
		fmt.Fprintf(os.Stderr, "\n[Error] checking npm '%s'\n", err.Error())
		os.Exit(-1)
	}
	fmt.Fprintf(os.Stdout, "[Warmup] using npm v%s", outstr)

	// Check bids-validator is installed
	outstr, err = valutils.AppVersionCheck(srvconfig.Exec.BIDS)
	if err != nil {
		fmt.Fprintf(os.Stderr, "\n[Error] checking bids-validator '%s'\n", err.Error())
		os.Exit(-1)
	}
	fmt.Fprintf(os.Stdout, "[Warmup] using bids-validator v%s", outstr)

	// Parse commandline arguments
	args, err := docopt.ParseArgs(usage, nil, "v1.0.0")
	if err != nil {
		fmt.Fprintf(os.Stderr, "\n[Error] parsing cli arguments: '%s', abort...\n\n", err.Error())
		os.Exit(-1)
	}
	fmt.Fprintf(os.Stdout, "[Warmup] cli arguments: %v\n", args)

	// Use port if provided.
	port := ":3033"
	if valutils.IsValidPort(args["<port>"]) {
		p := args["<port>"]
		port = fmt.Sprintf(":%s", p.(string))
	} else {
		fmt.Fprintln(os.Stderr, "[Info] could not parse a valid port number, using default")
	}
	fmt.Fprintf(os.Stdout, "[Warmup] using port: '%s'\n", port)

	fmt.Fprintln(os.Stdout, "[Warmup] registering routes")
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

	fmt.Fprintln(os.Stdout, "[Start] Listen and serve")
	err = server.ListenAndServe()
	if err != nil {
		fmt.Fprintf(os.Stderr, "\n[Error] Server startup: '%v', abort...\n\n", err)
		os.Exit(-1)
	}
}
