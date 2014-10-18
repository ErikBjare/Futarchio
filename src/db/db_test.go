package db

import (
	"testing"
)

func TestPassword(t *testing.T) {
	pass := "123password!\"#¤"
	user := NewUser("erb", pass, "erik@bjareho.lt", "Erik")
	valid := user.CheckPassword(pass)
	if !valid {
		t.Fatal("Password check failed")
	}
}

func TestAuth(t *testing.T) {
	auth := NewAuth()
	// TODO: testing
	if auth.Key == "" {
		t.Fatal("Authkey was empty string")
	}
}

func TestPoll(t *testing.T) {
	user := NewUser("test", "test", "test", "test")
	vote1, _ := NewYesNoVote(true, user, true)
	vote2, _ := NewYesNoVote(false, user, true)
	votes := []Vote{vote1, vote2}

	sumvotes := SumVotes(votes)
	shouldbe := map[string]float32{"yes": 1.0, "no": 1.0}
	for k := range sumvotes {
		if sumvotes[k] != shouldbe[k] {
			t.Fatal("SumVotes was wrong")
		}
	}
}
