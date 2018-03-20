package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/damondouglas/oauth"
	"google.golang.org/appengine"
)

func main() {
	login, err := oauth.LoginHandler("./secret/client_secret.json", "http://localhost:8080/oauth2callback", nil, true, false)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	http.HandleFunc("/", login)
	appengine.Main()
}

func handle(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}
