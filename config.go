package oauth

import (
	"encoding/json"
	"io/ioutil"
)

type config struct {
	Web *WebOauthConfig `json:"web"`
}

// WebOauthConfig stores oauth web client credentials
type WebOauthConfig struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}

func configureOauthFromData(data []byte) *WebOauthConfig {
	var c *config
	json.Unmarshal(data, &c)
	return c.Web
}

func configureOauthFromFilePath(pathToCredentials string) (*WebOauthConfig, error) {
	data, err := ioutil.ReadFile(pathToCredentials)
	if err != nil {
		return nil, err
	}
	return configureOauthFromData(data), nil
}
