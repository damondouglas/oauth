package oauth

import (
	"testing"

	"github.com/damondouglas/assert"
)

func TestBuildSessionID(t *testing.T) {
	a := assert.New(t)
	sessionID, err := buildSessionID()
	a.Equals(err, nil)
	a.NotEquals(sessionID, nil)
	a.NotEquals(sessionID, "")
}
