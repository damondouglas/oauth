package oauth

import (
	"encoding/gob"
	"errors"
	"net/url"

	"golang.org/x/oauth2"
)

func init() {
	// Gob encoding for gorilla/sessions
	gob.Register(&oauth2.Token{})
}

// RedirectURL holds oauth flow url data
type RedirectURL struct {
	SessionID string
	Value     string
}

func newRedirectURL(config *oauth2.Config, force bool, offline bool) (*RedirectURL, error) {
	var r *RedirectURL
	var err error
	r = new(RedirectURL)
	r.SessionID, err = buildSessionID()
	if err != nil {
		return nil, err
	}

	var offlineOpt oauth2.AuthCodeOption
	if offline {
		offlineOpt = oauth2.AccessTypeOffline
	} else {
		offlineOpt = oauth2.AccessTypeOnline
	}
	if force {
		r.Value = config.AuthCodeURL(r.SessionID, oauth2.ApprovalForce, offlineOpt)
	} else {
		r.Value = config.AuthCodeURL(r.SessionID, offlineOpt)
	}
	return r, nil
}

// Auth encapsulates Oauth2 details
type Auth struct {
	config *oauth2.Config
}

// New instantiates Auth
func New(pathToCredentials string, oauthFlowURL string, scopes []string) (*Auth, error) {
	a := new(Auth)
	var err error
	a.config, err = buildOauthClientFromCredentialsPath(pathToCredentials, oauthFlowURL, scopes)
	if err != nil {
		return nil, err
	}
	return a, nil
}

// BuildRedirectURL builds url for oauth flow
func (a *Auth) BuildRedirectURL(force bool, offline bool) (*RedirectURL, error) {
	return newRedirectURL(a.config, force, offline)
}

func validateRedirectURL(path string) (string, error) {
	if path == "" {
		return "/", nil
	}

	// Ensure redirect URL is valid and not pointing to a different server.
	parsedURL, err := url.Parse(path)
	if err != nil {
		return "/", err
	}
	if parsedURL.IsAbs() {
		return "/", errors.New("URL must not be absolute")
	}
	return path, nil
}
