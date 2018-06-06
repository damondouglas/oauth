package oauth

import (
	"encoding/json"
	"io/ioutil"
)

// Credentials represents OAuth2 secret credentials.
type Credentials struct {
	Web *WebOauthCredentials `json:"web"`
}

// WebOauthCredentials Web OAuth2 secret credentials.
type WebOauthCredentials struct {
	ClientID        string   `json:"client_id"`
	ClientSecret    string   `json:"client_secret"`
	RedirectURIList []string `json:"redirect_uris"`
}

func credentialsFromJSON(jsonData []byte) *WebOauthCredentials {
	var c *Credentials
	json.Unmarshal(jsonData, &c)
	return c.Web
}

// CredentialsFromPath loads WebOauthCredentials from given path.
func CredentialsFromPath(path string) (*WebOauthCredentials, error) {
	jsonData, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return credentialsFromJSON(jsonData), nil
}
