package web

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	gweb "github.com/G-Node/gin-cli/web"
	"github.com/G-Node/gin-valid/internal/config"
	"github.com/G-Node/gin-valid/internal/resources/templates"
	"github.com/gorilla/mux"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

var username = "valid-testing"
var reponame = "Testing"
var token = "4c82d07cccf103e071ad9ee8aec82c34d7003c6c"

func TestValiateBadConfig(t *testing.T) {
	valcfg, err := handleValidationConfig("wtf")
	if err == nil {
		t.Fatalf(`handleValidationConfig(cfgpath string) = %v, %v`, valcfg, err)
	}
}

func TestValidateNotYAML(t *testing.T) {
	f, _ := os.Create("testing-config.json")
	f.WriteString("foo: somebody said I should put a colon here: so I did")
	f.Close()
	valcfg, err := handleValidationConfig("testing-config.json")
	os.RemoveAll("testing-config.json")
	if err == nil {
		t.Fatalf(`handleValidationConfig(cfgpath string) = %v, %v`, valcfg, err)
	}

}

func TestValidateGoodConfig(t *testing.T) {
	f, _ := os.Create("testing-config.json")
	f.WriteString("empty: \"true\"")
	f.Close()
	valcfg, err := handleValidationConfig("testing-config.json")
	os.RemoveAll("testing-config.json")
	if err != nil {
		t.Fatalf(`handleValidationConfig(cfgpath string) = %v, %v`, valcfg, err)
	}
}

func TestValidateBIDSNoData(t *testing.T) {
	err := validateBIDS("wtf", "wtf")
	if err == nil {
		t.Fatalf(`validateBIDS(valroot, resdir string) = %v`, err)
	}
}

func TestValidateNIXNoData(t *testing.T) {
	err := validateNIX("wtf", "wtf")
	if err == nil {
		t.Fatalf(`validateNIX(valroot, resdir string) = %v`, err)
	}
}

func TestValidateODMLNoData(t *testing.T) {
	err := validateODML("wtf", "wtf")
	if err == nil {
		t.Fatalf(`validateODML(valroot, resdir string) = %v`, err)
	}
}

func TestValidateBadgeFail(t *testing.T) { //TODO
	body := []byte("{}")
	router := mux.NewRouter()
	router.HandleFunc("/validate/{validator}/{user}/{repo}", Validate).Methods("POST")
	srvcfg := config.Read()
	srvcfg.Dir.Tokens = "."
	os.Mkdir("tmp", 0755)
	srvcfg.Dir.Temp = "./tmp"
	srvcfg.GINAddresses.WebURL = weburl
	srvcfg.GINAddresses.GitURL = giturl
	config.Set(srvcfg)
	var tok gweb.UserToken
	tok.Username = username
	tok.Token = token
	saveToken(tok)
	os.Mkdir(filepath.Join(srvcfg.Dir.Tokens, "by-repo"), 0755)
	linkToRepo(username, filepath.Join(username, "/", reponame))
	r, err := http.NewRequest("POST", filepath.Join("/validate/bids/", username, "/", reponame), bytes.NewReader(body))
	if err != nil {
		t.Fatalf(`Validate(w http.ResponseWriter, r *http.Request) = %v`, err)
	}
	w := httptest.NewRecorder()
	sig := hmac.New(sha256.New, []byte(srvcfg.Settings.HookSecret))
	sig.Write(body)
	r.Header.Add("X-Gogs-Signature", hex.EncodeToString(sig.Sum(nil)))
	os.Mkdir("tmp", 0755)
	router.ServeHTTP(w, r)
	time.Sleep(5 * time.Second) //TODO HACK
	os.RemoveAll(filepath.Join(srvcfg.Dir.Tokens, "by-repo"))
	status := w.Code
	if status != http.StatusOK {
		t.Fatalf(`Validate(w http.ResponseWriter, r *http.Request) status code = %v`, status)
	}

}

func TestValidatePubBrokenPubValidate(t *testing.T) {
	original := templates.PubValidate
	templates.PubValidate = "{{ WTF? }"
	srvcfg := config.Read()
	body := []byte("{}")
	router := mux.NewRouter()
	router.HandleFunc("/pubvalidate", PubValidateGet).Methods("GET")
	r, _ := http.NewRequest("GET", "/pubvalidate", bytes.NewReader(body))
	w := httptest.NewRecorder()
	sig := hmac.New(sha256.New, []byte(srvcfg.Settings.HookSecret))
	sig.Write(body)
	r.Header.Add("X-Gogs-Signature", hex.EncodeToString(sig.Sum(nil)))
	router.ServeHTTP(w, r)
	templates.PubValidate = original
	status := w.Code
	if status != http.StatusInternalServerError {
		t.Fatalf(`Validate(w http.ResponseWriter, r *http.Request) status code = %v`, status)
	}
}

func TestValidatePubBrokenLayout(t *testing.T) {
	original := templates.Layout
	templates.Layout = "{{ WTF? }"
	srvcfg := config.Read()
	body := []byte("{}")
	router := mux.NewRouter()
	router.HandleFunc("/pubvalidate", PubValidateGet).Methods("GET")
	r, _ := http.NewRequest("GET", "/pubvalidate", bytes.NewReader(body))
	w := httptest.NewRecorder()
	sig := hmac.New(sha256.New, []byte(srvcfg.Settings.HookSecret))
	sig.Write(body)
	r.Header.Add("X-Gogs-Signature", hex.EncodeToString(sig.Sum(nil)))
	router.ServeHTTP(w, r)
	templates.Layout = original
	status := w.Code
	if status != http.StatusInternalServerError {
		t.Fatalf(`Validate(w http.ResponseWriter, r *http.Request) status code = %v`, status)
	}
}

func TestValidatePub(t *testing.T) {
	srvcfg := config.Read()
	body := []byte("{}")
	router := mux.NewRouter()
	router.HandleFunc("/pubvalidate", PubValidateGet).Methods("GET")
	r, _ := http.NewRequest("GET", "/pubvalidate", bytes.NewReader(body))
	w := httptest.NewRecorder()
	sig := hmac.New(sha256.New, []byte(srvcfg.Settings.HookSecret))
	sig.Write(body)
	r.Header.Add("X-Gogs-Signature", hex.EncodeToString(sig.Sum(nil)))
	router.ServeHTTP(w, r)
	status := w.Code
	if status != http.StatusOK {
		t.Fatalf(`Validate(w http.ResponseWriter, r *http.Request) status code = %v`, status)
	}
}

func TestValidateTMPFail(t *testing.T) {
	body := []byte("{}")
	router := mux.NewRouter()
	router.HandleFunc("/validate/{validator}/{user}/{repo}", Validate).Methods("POST")
	srvcfg := config.Read()
	srvcfg.Dir.Tokens = "."
	config.Set(srvcfg)
	var tok gweb.UserToken
	tok.Username = username
	tok.Token = token
	saveToken(tok)
	os.Mkdir(filepath.Join(srvcfg.Dir.Tokens, "by-repo"), 0755)
	linkToRepo(username, filepath.Join(username, "/", reponame))
	r, _ := http.NewRequest("POST", filepath.Join("/validate/bids/", username, "/", reponame), bytes.NewReader(body))
	w := httptest.NewRecorder()
	sig := hmac.New(sha256.New, []byte(srvcfg.Settings.HookSecret))
	sig.Write(body)
	r.Header.Add("X-Gogs-Signature", hex.EncodeToString(sig.Sum(nil)))
	router.ServeHTTP(w, r)
	time.Sleep(5 * time.Second) //TODO HACK
	status := w.Code
	if status != http.StatusOK {
		t.Fatalf(`Validate(w http.ResponseWriter, r *http.Request) status code = %v`, status)
	}
}

func TestValidateRepoDoesNotExists(t *testing.T) {
	token2 := "wtf"
	body := []byte("{}")
	router := mux.NewRouter()
	router.HandleFunc("/validate/{validator}/{user}/{repo}", Validate).Methods("POST")
	srvcfg := config.Read()
	srvcfg.Dir.Tokens = "."
	config.Set(srvcfg)
	var tok gweb.UserToken
	tok.Username = username
	tok.Token = token2
	saveToken(tok)
	os.Mkdir(filepath.Join(srvcfg.Dir.Tokens, "by-repo"), 0755)
	linkToRepo(username, filepath.Join(username, "/", reponame))
	r, _ := http.NewRequest("POST", filepath.Join("/validate/bids/", username, "/", reponame), bytes.NewReader(body))
	w := httptest.NewRecorder()
	sig := hmac.New(sha256.New, []byte(srvcfg.Settings.HookSecret))
	sig.Write(body)
	r.Header.Add("X-Gogs-Signature", hex.EncodeToString(sig.Sum(nil)))
	router.ServeHTTP(w, r)
	time.Sleep(5 * time.Second) //TODO HACK
	status := w.Code
	if status != http.StatusNotFound {
		t.Fatalf(`Validate(w http.ResponseWriter, r *http.Request) status code = %v`, status)
	}
}

func TestValidateBadToken(t *testing.T) {
	body := []byte("{}")
	router := mux.NewRouter()
	router.HandleFunc("/validate/{validator}/{user}/{repo}", Validate).Methods("POST")
	r, _ := http.NewRequest("POST", "/validate/bids/whatever/whatever", bytes.NewReader(body))
	w := httptest.NewRecorder()
	srvcfg := config.Read()
	sig := hmac.New(sha256.New, []byte(srvcfg.Settings.HookSecret))
	sig.Write(body)
	r.Header.Add("X-Gogs-Signature", hex.EncodeToString(sig.Sum(nil)))
	router.ServeHTTP(w, r)
	status := w.Code
	if status != http.StatusUnauthorized {
		t.Fatalf(`Validate(w http.ResponseWriter, r *http.Request) status code = %v`, status)
	}
}

func TestValidateUnsupportedValidator(t *testing.T) {
	body := []byte("{}")
	r, _ := http.NewRequest("GET", "wtf", bytes.NewReader(body))
	srvcfg := config.Read()
	srvcfg.Settings.HookSecret = "hooksecret"
	config.Set(srvcfg)
	sig := hmac.New(sha256.New, []byte(srvcfg.Settings.HookSecret))
	sig.Write(body)
	r.Header.Add("X-Gogs-Signature", hex.EncodeToString(sig.Sum(nil)))
	w := httptest.NewRecorder()
	Validate(w, r)
	status := w.Code
	if status != http.StatusNotFound {
		t.Fatalf(`Validate(w http.ResponseWriter, r *http.Request) status code = %v`, status)
	}
}

func TestValidateHookSecretFailed(t *testing.T) {
	r, _ := http.NewRequest("GET", "wtf", strings.NewReader("{}"))
	srvcfg := config.Read()
	srvcfg.Settings.HookSecret = "hooksecret"
	config.Set(srvcfg)
	r.Header.Add("X-Gogs-Signature", "wtf")
	w := httptest.NewRecorder()
	Validate(w, r)
	status := w.Code
	if status != http.StatusBadRequest {
		t.Fatalf(`Validate(w http.ResponseWriter, r *http.Request) status code = %v`, status)
	}
}

func TestValidateBodyNotJSON(t *testing.T) {
	r, _ := http.NewRequest("GET", "wtf", strings.NewReader("wtf"))
	w := httptest.NewRecorder()
	Validate(w, r)
	status := w.Code
	if status != http.StatusBadRequest {
		t.Fatalf(`Validate(w http.ResponseWriter, r *http.Request) status code = %v`, status)
	}
}

type errReader int

func (errReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("test error")
}

func TestValidateBadBody(t *testing.T) {
	testRequest := httptest.NewRequest(http.MethodPost, "/something", errReader(0))
	w := httptest.NewRecorder()
	Validate(w, testRequest)
	status := w.Code
	if status != http.StatusBadRequest {
		t.Fatalf(`Validate(w http.ResponseWriter, r *http.Request) status code = %v`, status)
	}
}
