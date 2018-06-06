package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/damondouglas/oauth/storage/aedatastore"

	"github.com/damondouglas/oauth"
	"github.com/damondouglas/oauth/authorization"
	"github.com/damondouglas/oauth/exchange"
	helper "github.com/damondouglas/oauth/helper/appengine"

	"google.golang.org/appengine"
)

const pathToTunnelInfo = "tunnels"

const pathToSecret = "secret/client_secret.json"

var (
	scopes = []string{"email", "profile"}
)

func main() {
	// This is just a bunch of yada yada to set up this demo.
	// Just ignore it.
	tunnelURL, err := getTunnelURL(pathToTunnelInfo)
	redirectURL := fmt.Sprintf("%s/step1", tunnelURL.String())
	if err != nil {
		log.Panic(err)
	}
	http.HandleFunc("/", root)
	http.HandleFunc("/step1", step1)
	http.HandleFunc("/info", info)
	// End of the demo yada yada.

	// Here is the real tofu of this demo.
	// (I'm vegetarian so I say tofu instead of meat. ❤️🐂🐔🐐❤️)
	// You probably might put `redirectURL` in your app.yaml environment_variables section.
	config, err := oauth.ConfigFromPath(pathToSecret, redirectURL, scopes)
	if err != nil {
		log.Panic(err)
	}
	auth := authorization.New(config, redirectURL, scopes, true, true)
	ctxHelper := new(helper.ContextHelper)
	cltHelper := new(helper.ClientHelper)
	h := new(exchange.Helper)
	h.ClientHelper = cltHelper
	h.ContextHelper = ctxHelper
	h.StorageHelper = new(aedatastore.Datastore)
	if err != nil {
		log.Panic(err)
	}
	exch := exchange.New(config, h)

	http.HandleFunc("/auth", auth.Handle)
	http.HandleFunc("/token", exch.Handle)
	appengine.Main()
}

// The entire code below is just more yada yada to set up this demo.
// Just ignore it.

type infoRequest struct {
	IDToken string `json:"id_token"`
}

type infoResponse struct {
	Email string
}

func info(w http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Panic(err)
	}
	defer r.Body.Close()
	var i *infoRequest
	resp := new(infoResponse)
	err = json.Unmarshal(data, &i)
	datastore := new(aedatastore.Datastore)
	ctxHelper := new(helper.ContextHelper)
	ctx := ctxHelper.Context(r)
	u, err := datastore.GetByIDToken(ctx, i.IDToken)
	if err != nil {
		log.Panic(err)
	}
	resp.Email = u.Email
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func root(w http.ResponseWriter, r *http.Request) {
	content, err := buildClientHTTP(r, pathToSecret)
	if err != nil {
		log.Panic(err)
	}
	fmt.Fprintln(w, content)
}

func step1(w http.ResponseWriter, r *http.Request) {
	qry := r.URL.Query()
	lines := []string{}
	lines = append(lines, fmt.Sprintf("@code = %s", qry.Get("code")))

	content := strings.Join(lines, "\n")
	fmt.Fprintln(w, content)
}

func buildClientHTTP(r *http.Request, pathToSecret string) (string, error) {
	tunnelURL, err := getTunnelURL(pathToTunnelInfo)
	if err != nil {
		log.Panic(err)
	}
	redirectURL := fmt.Sprintf("%s/step1", tunnelURL.String())
	lines := []string{}
	lines = append(lines, fmt.Sprintf("@base = %s", tunnelURL))
	lines = append(lines, fmt.Sprintf("@redirect = %s", redirectURL))
	lines = append(lines, fmt.Sprintf("@state = %s", randomString(12)))
	lines = append(lines, fmt.Sprintf("@scope = %s", strings.Join(scopes, "+")))

	config, err := oauth.CredentialsFromPath(pathToSecret)
	if err != nil {
		return "", err
	}
	lines = append(lines, fmt.Sprintf("@client_id = %s", config.ClientID))
	lines = append(lines, fmt.Sprintf("@client_secret = %s", config.ClientSecret))
	lines = append(lines, "###################################")
	lines = append(lines, "# STEP 3")
	lines = append(lines, fmt.Sprintf("# Open the following URL and add %s to Authorized redirect URIs:", redirectURL))
	lines = append(lines, "# ")
	lines = append(lines, fmt.Sprintf("# https://console.developers.google.com/apis/credentials/oauthclient/%s", config.ClientID))
	lines = append(lines, "###################################")

	return strings.Join(lines, "\n"), nil
}

type tunnelData struct {
	Tunnels []struct {
		URL   string `json:"public_url"`
		Proto string
	}
}

func getTunnelURL(path string) (*url.URL, error) {
	var u *url.URL
	var t *tunnelData
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(data, &t)
	if err != nil {
		return nil, err
	}
	for _, info := range t.Tunnels {
		if info.Proto == "https" {
			u, err = url.Parse(info.URL)
			if err != nil {
				return nil, err
			}
		}
	}

	return u, nil
}

func randomString(n int) string {
	rand.Seed(time.Now().UnixNano())
	letterRunes := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
