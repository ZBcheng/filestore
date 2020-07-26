package repository

import (
	"testing"
)

func TestGenToken(t *testing.T) {
	token := GenToken("bee")
	t.Log(len(token))
	t.Log(token)
}

func TestAuthToken(t *testing.T) {
	token := GenToken("bee")
	AuthToken("bee", token)
}
