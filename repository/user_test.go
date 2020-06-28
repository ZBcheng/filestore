package repository

import "testing"

func TestSingup(t *testing.T) {
	_, err := UserSignup("Amy", "amy")
	if err != nil {
		t.Error(err)
	}
}

func TestSignin(t *testing.T) {
	_, err := UserSignin("Amy", "amy or WHERE 1=1")
	if err != nil {
		t.Error(err)
	}
}

func TestUpdateToken(t *testing.T) {
	_, err := UpdateToken("Amy", "0101010")
	if err != nil {
		t.Error(err)
	}
}
