package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/damondouglas/oauth"
	"google.golang.org/appengine"
)

func main() {
	auth, err := oauth.New("./secret/client_secret.json", "http://localhost:8080/oauth2callback", nil)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	login := oauth.LoginHandler(auth, true, false)
	callback := oauth.CallbackHandler(auth)
	http.HandleFunc("/", login)
	http.HandleFunc("/oauth2callback", callback)
	http.HandleFunc("/profile", getProfileHandler(auth))
	appengine.Main()
}

func getProfileHandler(auth *oauth.Auth) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello, world!")
	}
}
