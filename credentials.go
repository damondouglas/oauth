package oauth

import (
	"encoding/json"
	"io/ioutil"
)

type credentials struct {
	Web *webOauthCredentials `json:"web"`
}

type webOauthCredentials struct {
	ClientID        string   `json:"client_id"`
	ClientSecret    string   `json:"client_secret"`
	RedirectURIList []string `json:"redirect_uris"`
}

func credentialsFromJSON(jsonData []byte) *webOauthCredentials {
	var c *credentials
	json.Unmarshal(jsonData, &c)
	return c.Web
}

func credentialsFromPath(path string) (*webOauthCredentials, error) {
	jsonData, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return credentialsFromJSON(jsonData), nil
}
