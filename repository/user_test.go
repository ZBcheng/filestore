package repository

import (
	"testing"

	"github.com/zbcheng/filestore/models"
)

var user = &models.User{}

func TestSingup(t *testing.T) {
	suc := UserSignup("Amy2", "amy", "", "")
	if !suc {
		t.Error("Error")
	}
}

func TestGenToken(t *testing.T) {
	token := GenToken("bee")
	t.Log(len(token))
	t.Log(token)
}

func TestAuthToken(t *testing.T) {
	token := GenToken("bee")
	AuthToken("bee", token)
}
