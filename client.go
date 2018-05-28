package oauth

import (
	"sort"

	"golang.org/x/oauth2/google"

	"golang.org/x/oauth2"
)

var baseScopes = []string{"email", "profile"}

func mergeScopesFromBase(scopes []string) []string {
	var mergedScopes []string
	for _, value := range baseScopes {
		mergedScopes = append(mergedScopes, value)
	}
	sort.Strings(mergedScopes)
	for _, value := range scopes {
		i := sort.SearchStrings(mergedScopes, value)
		if i == len(mergedScopes) {
			mergedScopes = append(mergedScopes, value)
			sort.Strings(mergedScopes)
		}
		if i < len(mergedScopes) && mergedScopes[i] != value {
			mergedScopes = append(mergedScopes, value)
			sort.Strings(mergedScopes)
		}
	}
	return mergedScopes
}

func buildOauthClient(clientID, clientSecret, redirectURL string, scopes []string) *oauth2.Config {
	mergedScopes := mergeScopesFromBase(scopes)
	return &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURL,
		Scopes:       mergedScopes,
		Endpoint:     google.Endpoint,
	}
}

func buildOauthClientFromCredentialsPath(path string, redirectURL string, scopes []string) (*oauth2.Config, error) {
	config, err := configureOauthFromFilePath(path)
	if err != nil {
		return nil, err
	}
	return buildOauthClient(config.ClientID, config.ClientSecret, redirectURL, scopes), nil
}
