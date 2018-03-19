package oauth

import uuid "github.com/satori/go.uuid"

const (
	defaultSessionID        = "default"
	googleProfileSessionKey = "google_profile"
	oauthTokenSessionKey    = "oauth_token"
	oauthFlowRedirectKey    = "redirect"
)

func buildSessionID() (string, error) {
	var sessionID string
	u, err := uuid.NewV4()
	if err != nil {
		return "", err
	}
	sessionID = u.String()
	return sessionID, nil
}
