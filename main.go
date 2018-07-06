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
)

const usage = `Server validating BIDS files

Usage:
  ginvalid
  ginvalid -h | --help
  ginvalid --version

Options:
  -h --help           Show this screen.
  --version           Print version.
  `

func root(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "ginvalid running")
}

func validate(w http.ResponseWriter, r *http.Request) {
	repo := mux.Vars(r)["repo"]
	user := mux.Vars(r)["user"]
	fmt.Fprintf(w, "validate repo '%s/%s'\n", repo, user)

	cmd := exec.Command("gin", "repoinfo", fmt.Sprintf("%s/%s", repo, user))
	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "[Error] accessing '%s/%s': '%s'\n", repo, user, err.Error())
		return
	}

	tmpdir, err := ioutil.TempDir("/home/msonntag/Chaos/DL/val", "validator")
	if err != nil {
		fmt.Fprintf(os.Stderr, "[Error] creating temporary directory: '%s'\n", err.Error())
		return
	}
	fmt.Fprintf(w, "Directory created: %s", tmpdir)
}

func registerRoutes(r *mux.Router) {
	r.HandleFunc("/", root)
	r.HandleFunc("/validate/{repo}/{user}", validate)
}

func main() {

	// Check gin installed and available
	cmd := exec.Command("gin", "--version")
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "\n[Error] checking gin client '%s'\n", err.Error())
		os.Exit(-1)
	}
	fmt.Fprintf(os.Stdout, "[Warmup] using %s", out.String())

	// Check npm installed and available
	cmd = exec.Command("npm", "--version")
	out.Reset()
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "\n[Error] checking npm '%s'\n", err.Error())
		os.Exit(-1)
	}
	fmt.Fprintf(os.Stdout, "[Warmup] using npm v%s", out.String())

	// Check bids-validator is installed
	cmd = exec.Command("npm", "show", "bids-validator", "version")
	out.Reset()
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "\n[Error] checking bids-validator '%s'\n", err.Error())
		os.Exit(-1)
	}
	fmt.Fprintf(os.Stdout, "[Warmup] using bids-validator v%s", out.String())

	args, err := docopt.ParseArgs(usage, nil, "v1.0.0")
	if err != nil {
		fmt.Fprintf(os.Stderr, "\n[Error] parsing cli arguments: '%s', abort...\n\n", err.Error())
		os.Exit(-1)
	}
	fmt.Fprintf(os.Stdout, "[Warmup] cli arguments: %v\n", args)

	fmt.Fprintln(os.Stdout, "[Warmup] registering routes")
	router := mux.NewRouter()
	registerRoutes(router)

	handler := handlers.CORS(
		handlers.AllowedHeaders([]string{"Accept", "Content-Type", "Authorization"}),
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowedMethods([]string{"GET"}),
	)(router)

	server := http.Server{
		Addr:    ":3033",
		Handler: handler,
	}

	fmt.Fprintln(os.Stdout, "[Start] Listen and serve")
	err = server.ListenAndServe()
	if err != nil {
		fmt.Fprintf(os.Stderr, "\n[Error] Server startup: '%v', abort...\n\n", err)
		os.Exit(-1)
	}
}
