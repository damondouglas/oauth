package exchange

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"google.golang.org/appengine/urlfetch"

	"github.com/damondouglas/oauth/helper"
	"github.com/damondouglas/oauth/query"

	"github.com/damondouglas/oauth/common"
	"github.com/damondouglas/oauth/storage"
	"golang.org/x/oauth2"
)

const refreshTokenEndpoint = "https://www.googleapis.com/oauth2/v4/token"

// Helper provides context.Context and http.Client.
type Helper struct {
	ContextHelper helper.ContextHelper
	ClientHelper  helper.ClientHelper
	StorageHelper storage.Storage
}

// Exchange handles Oauth2 token exchange requests.
type Exchange struct {
	oauthConfig *oauth2.Config
	helper      *Helper
}

// New instantiates Exchange
func New(config *oauth2.Config, helper *Helper) *Exchange {
	e := new(Exchange)
	e.oauthConfig = config
	e.helper = helper
	return e
}

func read(r *http.Request) ([]byte, error) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()
	return body, nil
}

func (e *Exchange) isValid(qry *query.TokenQuery) bool {
	return qry.ClientID() == e.oauthConfig.ClientID &&
		qry.ClientSecret() == e.oauthConfig.ClientSecret
}

// Handle takes Oauth2 token exchange Post request and responds with token data.
// The method also stores resulting user data associated with the token response.
func (e *Exchange) Handle(w http.ResponseWriter, r *http.Request) {
	data, err := read(r)
	common.HandleError(w, err, "could not read request data")
	qry, err := query.ParseTokenQuery(string(data))
	common.HandleError(w, err, "could not parse token request form data")
	if !qry.IsValid() {
		common.HandleError(w, errors.New(string(data)), "query form data is invalid")
	}
	if qry.Mode == query.AuthorizationCodeMode {
		e.code(w, r, qry)
	}
	if qry.Mode == query.RefreshTokenMode {
		e.refresh(w, r, qry)
	}
}

/*
{
  token_type: "bearer",
  access_token: "ACCESS_TOKEN",
  refresh_token: "REFRESH_TOKEN",
  expires_in: SECONDS_TO_EXPIRATION
}
*/
type codeResponse struct {
	TokenType    string `json:"token_type"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token,omitempty"`
	Expiry       int    `json:"expires_in"`
}

// client_id=GOOGLE_CLIENT_ID&client_secret=GOOGLE_CLIENT_SECRET&grant_type=authorization_code&code=AUTHORIZATION_CODE
func (e *Exchange) code(w http.ResponseWriter, r *http.Request, qry *query.TokenQuery) {
	ctx := e.helper.ContextHelper.Context(r)
	code := qry.Code()
	tok, err := e.oauthConfig.Exchange(ctx, code)
	common.HandleError(w, err, "authorization code could not be exchanged")
	user, err := storage.GetUserFromToken(ctx, e.oauthConfig, tok)
	common.HandleError(w, err, "User info could not be loaded from exchanged token")
	user.AuthorizationCode = code
	err = e.helper.StorageHelper.Save(ctx, user)
	common.HandleError(w, err, "user info could not be stored")
	c := new(codeResponse)
	c.TokenType = "bearer"
	c.AccessToken = tok.AccessToken
	c.RefreshToken = tok.RefreshToken
	expiry := time.Until(tok.Expiry)
	c.Expiry = int(expiry.Seconds())
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(c)
}

//client_id=GOOGLE_CLIENT_ID&client_secret=GOOGLE_CLIENT_SECRET&grant_type=refresh_token&refresh_token=REFRESH_TOKEN
func (e *Exchange) refresh(w http.ResponseWriter, r *http.Request, qry *query.TokenQuery) {
	var c *codeResponse
	ctx := e.helper.ContextHelper.Context(r)
	client := urlfetch.Client(ctx)
	req, err := http.NewRequest("POST", refreshTokenEndpoint, bytes.NewBuffer([]byte(qry.Body)))
	common.HandleError(w, err, "Request could not be built")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := client.Do(req)
	common.HandleError(w, err, "Request could not be completed")
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	common.HandleError(w, err, "Response could not be read")
	err = json.Unmarshal(data, &c)
	refreshToken := qry.RefreshToken()
	u, err := e.helper.StorageHelper.GetByRefreshToken(ctx, refreshToken)
	common.HandleError(w, err, fmt.Sprintf("Could not retrieve user from refresh token: %s", refreshToken))
	u.Token.AccessToken = c.AccessToken
	err = e.helper.StorageHelper.Save(ctx, u)
	common.HandleError(w, err, fmt.Sprintf("Could not store user: %v", u))
	fmt.Fprintln(w, string(data))
}
