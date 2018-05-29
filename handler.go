package oauth

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"

	"golang.org/x/oauth2"
)

const tokenRefreshEndpoint = "https://www.googleapis.com/oauth2/v4/token"

type ClientBuilder interface {
	Client(context.Context) *http.Client
}

type ContextBuilder interface {
	Context(*http.Request) context.Context
}

// Handler provides oauth flow function handlers
type Handler struct {
	clientBuilder  ClientBuilder
	contextBuilder ContextBuilder
	cred           *webOauthConfig
	config         *oauth2.Config
	force          bool
	offline        bool
	redirectURL    string
	authURL        string
}

func NewHandler(clientBuilder ClientBuilder, contextBuilder ContextBuilder, pathToSecret string, force bool, offline bool, redirectURL string, scopes []string) (*Handler, error) {
	cred, err := configureOauthFromFilePath(pathToSecret)
	if err != nil {
		return nil, err
	}
	config, err := buildOauthClientFromCredentialsPath(pathToSecret, redirectURL, scopes)
	if err != nil {
		return nil, err
	}
	handler := new(Handler)
	handler.clientBuilder = clientBuilder
	handler.contextBuilder = contextBuilder
	handler.cred = cred
	handler.config = config
	handler.force = force
	handler.redirectURL = redirectURL

	return handler, nil
}

func (h *Handler) buildAuthorizeURL(state string) string {
	var url string

	var offlineOpt oauth2.AuthCodeOption
	if h.offline {
		offlineOpt = oauth2.AccessTypeOffline
	} else {
		offlineOpt = oauth2.AccessTypeOnline
	}
	if h.force {
		url = h.config.AuthCodeURL(state, oauth2.ApprovalForce, offlineOpt)
	} else {
		url = h.config.AuthCodeURL(state, offlineOpt)
	}

	return url
}

// Authorize handles oauth flow initialization.
func (h *Handler) Authorize(w http.ResponseWriter, r *http.Request) {
	qry := r.URL.Query()
	state := qry.Get("state")
	url := h.buildAuthorizeURL(state)

	http.Redirect(w, r, url, http.StatusFound)
}

type tokenRedirect struct {
	State    string
	Code     string
	Redirect string
}

func (h *Handler) Token(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		h.code(w, r)
	}
	if r.Method == "POST" {
		h.exchange(w, r)
	}
}

func (h *Handler) code(w http.ResponseWriter, r *http.Request) {
	qry := r.URL.Query()
	state := qry.Get("state")
	code := qry.Get("code")

	redirect := tokenRedirect{
		State:    state,
		Code:     code,
		Redirect: h.redirectURL,
	}

	urlTemplate, err := template.New("").Parse("{{.Redirect}}?code={{.Code}}&state={{.State}}")
	if err != nil {
		log.Panic(err)
	}
	var urlBytes bytes.Buffer
	if err := urlTemplate.Execute(&urlBytes, redirect); err != nil {
		log.Panic(err)
	}
	http.Redirect(w, r, urlBytes.String(), http.StatusFound)
}

func read(r *http.Request) ([]byte, error) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()
	return body, nil
}

func (h *Handler) exchange(w http.ResponseWriter, r *http.Request) {
	data, err := read(r)
	if err != nil {
		log.Panic(err)
	}
	qry, err := url.ParseQuery(string(data))
	if qry.Get("code") != "" {
		h.exchangeCode(w, r, data)
	}
	if qry.Get("refresh_token") != "" {
		h.exchangeRefresh(w, r, data)
	}
}

func (h *Handler) exchangeCode(w http.ResponseWriter, r *http.Request, data []byte) {
	qry, err := url.ParseQuery(string(data))
	if err != nil {
		log.Panic(err)
	}
	code := qry.Get("code")
	ctx := h.contextBuilder.Context(r)
	tok, err := h.config.Exchange(ctx, code)
	if err != nil {
		log.Panic(err)
	}
	tokData, err := json.Marshal(tok)
	if err != nil {
		log.Panic(err)
	}
	fmt.Fprintln(w, string(tokData))
}

func (h *Handler) exchangeRefresh(w http.ResponseWriter, r *http.Request, data []byte) {
	qry, err := url.ParseQuery(string(data))
	if err != nil {
		log.Panic(err)
	}
	code := qry.Get("refresh_token")
	qry, err = url.ParseQuery("")
	if err != nil {
		log.Panic(err)
	}

	qry.Add("client_id", h.cred.ClientID)
	qry.Add("client_secret", h.cred.ClientSecret)
	qry.Add("refresh_token", code)
	qry.Add("grant_type", "refresh_token")
	body := strings.NewReader(qry.Encode())
	ctx := h.contextBuilder.Context(r)
	client := h.clientBuilder.Client(ctx)
	resp, err := client.Post(tokenRefreshEndpoint, "application/x-www-form-urlencoded", body)
	if err != nil {
		log.Panic(err)
	}
	tokData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Panic(err)
	}
	defer resp.Body.Close()
	fmt.Fprintln(w, string(tokData))
}
