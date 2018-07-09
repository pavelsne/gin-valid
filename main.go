package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"

	"github.com/docopt/docopt-go"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/mpsonntag/gin-valid/config"
	"github.com/mpsonntag/gin-valid/valutils"
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

func validate(w http.ResponseWriter, r *http.Request) {
	user := mux.Vars(r)["user"]
	repo := mux.Vars(r)["repo"]
	fmt.Fprintf(w, "validate repo '%s/%s'\n", user, repo)

	cmd := exec.Command("gin", "repoinfo", fmt.Sprintf("%s/%s", user, repo))
	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "[Error] accessing '%s/%s': '%s'\n", user, repo, err.Error())
		return
	}

	tmpdir, err := ioutil.TempDir("/home/msonntag/Chaos/DL/val", "bidsval_")
	if err != nil {
		fmt.Fprintf(os.Stderr, "[Error] creating temporary directory: '%s'\n", err.Error())
		return
	}
	fmt.Fprintf(w, "Directory created: %s\n", tmpdir)

	// enable cleanup once tried and tested
	defer os.RemoveAll(tmpdir)

	cmd = exec.Command("gin", "get", fmt.Sprintf("%s/%s", user, repo))
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Dir = tmpdir
	if err = cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "[Error] running gin get: '%s'\n", err.Error())
		return
	}
	fmt.Fprintf(w, "running in %s, gin get: %s\n", cmd.Dir, out.String())

	// Ignoring NiftiHeaders for now, since it seems to be a common error
	cmd = exec.Command("/home/msonntag/node_modules/.bin/bids-validator", "--ignoreNiftiHeaders", "--json", fmt.Sprintf("%s/%s", tmpdir, repo))
	out.Reset()
	cmd.Stdout = &out
	var serr bytes.Buffer
	cmd.Stderr = &serr
	cmd.Dir = tmpdir
	if err = cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "[Error] running bids validation (%s): '%s', '%s', '%s'", fmt.Sprintf("%s/%s", tmpdir, repo), err.Error(), serr.String(), out.String())
		return
	}
	fmt.Fprintf(w, "validation successful: \n%s\n", out.String())
}

func registerRoutes(r *mux.Router) {
	r.HandleFunc("/", root)
	r.HandleFunc("/validate/{user}/{repo}", validate)
}

func main() {

	srvconfig := config.Read()

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
