package repository

import "testing"

func TestSingup(t *testing.T) {
	_, err := UserSignup("Amy", "amy")
	if err != nil {
		t.Error(err)
	}
}

func TestSignin(t *testing.T) {
	_, err := UserSignin("Amy", "amy")
	if err != nil {
		t.Error(err)
	}
}