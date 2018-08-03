package web

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/G-Node/gin-cli/ginclient"
	"github.com/G-Node/gin-cli/web"
	"github.com/gorilla/mux"
)

const tokenform = `
<html>
    <head>
    <title></title>
    </head>
    <body>
        <form action="/newuser" method="post">
            Username: <input type="text" name="username">
            Token:    <input type="text" name="token">
            <input type="submit" value="Submit">
        </form>
    </body>
</html>
`

func SetToken(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		tf := template.New("token")
		tf.Parse(tokenform)
		tf.Execute(w, nil)
	} else {
		r.ParseForm()
		username := r.Form["username"][0]
		token := r.Form["token"][0]
		ut := web.UserToken{Username: username, Token: token}
		ut.StoreToken(serveralias)
	}
}

const repostmpl = `
<html>
    <head>
    <title></title>
    </head>
    <body>
        {{ range . }}

			<p><b>{{.FullName}}</b></p>
			<p>{{.Description}} {{.Website}}</p>
			<hr>
		{{ end }}
    </body>
</html>
`

func ListRepos(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		return
	}
	vars := mux.Vars(r)
	user := vars["user"]
	cl := ginclient.New(serveralias)
	cl.LoadToken()
	fmt.Printf("Requesting repository listing for user %s\n", user)
	fmt.Printf("Server alias: %s\n", serveralias)
	fmt.Println("Server configuration:")
	fmt.Println(cl.Host)

	repos, err := cl.ListRepos(user)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Printf("Got %d repos\n", len(repos))

	rl := template.New("repos")
	rl.Parse(repostmpl)
	rl.Execute(w, &repos)
}
