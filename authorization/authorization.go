package authorization

import (
	"errors"
	"log"
	"net/http"

	"github.com/damondouglas/oauth/common"
	"github.com/damondouglas/oauth/query"
	"golang.org/x/oauth2"
)

const (
	codeKey = "code"
)

// Authorization handles OAuth2 authorization requests.
type Authorization struct {
	oauthConfig *oauth2.Config
	force       bool
	offline     bool
	redirectURL string
	scopes      []string
}

// New instantiates Authorization from oauth2.Config and specified details.
// redirectURL dictates where to redirect user after authorizing permissions.
// Setting force or offline to true sets, query parameters prompt and access_type, respectively.
func New(oauthConfig *oauth2.Config, redirectURL string, scopes []string, force, offline bool) *Authorization {
	a := new(Authorization)
	a.oauthConfig = oauthConfig
	a.redirectURL = redirectURL
	a.scopes = scopes
	a.force = force
	a.offline = offline
	return a
}

func (a *Authorization) isValid(aq *query.AuthorizationQuery) bool {
	return a.oauthConfig.ClientID == aq.ClientID() &&
		a.oauthConfig.RedirectURL == aq.RedirectURI() &&
		aq.ResponseType() == codeKey
}

// Handle builds and forwards user authorization URL
func (a *Authorization) Handle(w http.ResponseWriter, r *http.Request) {
	aq, err := query.ParseAuthorizationQuery(r.URL.RawQuery)
	common.HandleError(w, err, "authorization query string could not be parsed")
	if !aq.IsValid() {
		log.Println(aq)
		common.HandleError(w, errors.New(r.URL.RawQuery), "authorization query string could not be parsed")
	}
	if !a.isValid(aq) {
		log.Println(a.oauthConfig.ClientID, aq.ClientID(), a.oauthConfig.ClientID == aq.ClientID())
		log.Println(a.oauthConfig.RedirectURL, aq.RedirectURI(), a.oauthConfig.RedirectURL == aq.RedirectURI())
		log.Println(aq.ResponseType(), codeKey, aq.ResponseType() == codeKey)
		common.HandleError(w, errors.New(r.URL.RawQuery), "authorization query is not valid")
	}
	a.oauthConfig.ClientID = aq.ClientID()
	a.oauthConfig.RedirectURL = aq.RedirectURI()
	url := a.buildURL(aq.State())
	http.Redirect(w, r, url, http.StatusFound)
}

func (a *Authorization) buildURL(state string) string {
	var url string
	var offlineOpt oauth2.AuthCodeOption
	if a.offline {
		offlineOpt = oauth2.AccessTypeOffline
	} else {
		offlineOpt = oauth2.AccessTypeOnline
	}
	if a.force {
		url = a.oauthConfig.AuthCodeURL(state, oauth2.ApprovalForce, offlineOpt)
	} else {
		url = a.oauthConfig.AuthCodeURL(state, offlineOpt)
	}

	return url
}
