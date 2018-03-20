package oauth

import (
	"net/http"

	"github.com/gorilla/sessions"
	uuid "github.com/satori/go.uuid"
)

const (
	defaultSessionID        = "default"
	googleProfileSessionKey = "google_profile"
	oauthTokenSessionKey    = "oauth_token"
	oauthFlowRedirectKey    = "redirect"
)

var (
	sessionStore sessions.Store
)

func init() {
	sessionStore = sessions.NewCookieStore([]byte(buildRandomString()))
}

func buildRandomString() string {
	return uuid.Must(uuid.NewV4()).String()
}

func getOauthFlowSession(r *http.Request, sessionID string) (*sessions.Session, error) {
	var oauthFlowSession *sessions.Session
	oauthFlowSession, err := sessionStore.New(r, sessionID)
	if err != nil {
		return nil, err
	}
	oauthFlowSession.Options.MaxAge = 10 * 60 // 10 minutes
	return oauthFlowSession, nil
}
