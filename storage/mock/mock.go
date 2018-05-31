package mock

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/damondouglas/oauth/storage"
)

// Mock implements Storage for testing purposes.
// It uses the local file system.
// You should not use this in production.
type Mock struct {
	root             string
	data             string
	emailIndex       string
	authCodeIndex    string
	refreshCodeIndex string
}

// New instantiates Mock from pathToRoot.
func New(pathToRoot string) (*Mock, error) {
	m := new(Mock)
	m.root = pathToRoot
	m.data = path.Join(m.root, "data")
	m.emailIndex = path.Join(m.root, "email")
	m.authCodeIndex = path.Join(m.root, "auth")
	m.refreshCodeIndex = path.Join(m.root, "refresh")
	return m, nil
}

func writeFile(path string, content []byte) error {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		file, err := os.Create(path)
		if err != nil {
			return err
		}
		defer file.Close()
	}

	err = ioutil.WriteFile(path, content, 0644)
	if err != nil {
		return err
	}

	return nil
}

func encodeEmail(email string) string {
	// encoded := strings.Replace(email, "@", "_at_", -1)
	encoded := strings.Replace(email, ".", "_dot_", -1)
	return encoded
}

// Save saves user into mock database.
func (m *Mock) Save(u *storage.User) error {
	data, err := json.Marshal(u)
	if err != nil {
		return err
	}
	pathToData := path.Join(m.data, u.ID)
	err = writeFile(pathToData, data)
	if err != nil {
		return err
	}

	pathToEmail := path.Join(m.emailIndex, encodeEmail(u.Email))
	err = writeFile(pathToEmail, []byte(pathToData))
	if err != nil {
		return err
	}

	pathToAuth := path.Join(m.authCodeIndex, u.AuthorizationCode)
	err = writeFile(pathToAuth, []byte(pathToData))
	if err != nil {
		return err
	}

	pathToRefresh := path.Join(m.refreshCodeIndex, u.Token.RefreshToken)
	err = writeFile(pathToRefresh, []byte(pathToData))

	return nil
}

// GetByUserID retrieves storage.User from ID
func (m *Mock) GetByUserID(ID string) (*storage.User, error) {
	pathToData := path.Join(m.data, ID)
	return unmarshal(pathToData)
}

func unmarshal(pathToData string) (*storage.User, error) {
	var u *storage.User
	data, err := ioutil.ReadFile(pathToData)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(data, &u)
	if err != nil {
		return nil, err
	}

	return u, nil
}

// GetByEmail acquires storage.User from mock database using email.
func (m *Mock) GetByEmail(email string) (*storage.User, error) {
	pathToEmailIndex := path.Join(m.emailIndex, encodeEmail(email))
	pathData, err := ioutil.ReadFile(pathToEmailIndex)
	if err != nil {
		return nil, err
	}
	return unmarshal(string(pathData))
}

// GetByAuthorizationCode acquires storage.User from mock database using granted authorization code.
func (m *Mock) GetByAuthorizationCode(code string) (*storage.User, error) {
	pathToIndex := path.Join(m.authCodeIndex, code)
	pathData, err := ioutil.ReadFile(pathToIndex)
	if err != nil {
		return nil, err
	}
	return unmarshal(string(pathData))
}

// GetByRefreshToken acquires storage.User from mock database using granted refresh token code.
func (m *Mock) GetByRefreshToken(code string) (*storage.User, error) {
	pathToIndex := path.Join(m.refreshCodeIndex, code)
	pathData, err := ioutil.ReadFile(pathToIndex)
	if err != nil {
		return nil, err
	}
	return unmarshal(string(pathData))
}
