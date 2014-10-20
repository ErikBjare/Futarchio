package db

import (
	"appengine/aetest"
	"appengine/datastore"
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
	c, err := aetest.NewContext(nil)
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()

	user := NewUser("user", "pass", "name", "email")
	_ = NewYesNoPoll("title", "desc", "erb")
	pollkey := datastore.NewIncompleteKey(c, "Poll", nil)
	vote1, _, _ := NewYesNoVote(pollkey, user, true, 0)
	vote2, _, _ := NewYesNoVote(pollkey, user, false, 0)

	// TODO: Verify that polls are actually made by their creator
	// TODO: Test different levels of privacy

	if vote1.Weights()["yes"] != vote2.Weights()["no"] {
		t.Fatal("")
	}

	votes := []Vote{vote1, vote2}
	sumvotes := SumVotes(votes)
	shouldbe := map[string]float32{"yes": 1.0, "no": 1.0}
	for k := range sumvotes {
		if sumvotes[k] != shouldbe[k] {
			t.Fatal("SumVotes was wrong")
		}
	}
}
