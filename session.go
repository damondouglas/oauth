package oauth

import uuid "github.com/satori/go.uuid"

func buildSessionID() (string, error) {
	var sessionID string
	u, err := uuid.NewV4()
	if err != nil {
		return "", err
	}
	sessionID = u.String()
	return sessionID, nil
}
