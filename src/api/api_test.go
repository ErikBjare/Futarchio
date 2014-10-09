package api

import (
	//	"gopkg.in/mgo.v2"
	"appengine/aetest"
	"appengine/datastore"
	"github.com/ErikBjare/Futarchio/src/db"
	"log"
	"testing"
)

func TestUsers(t *testing.T) {
	c, err := aetest.NewContext(nil)
	log.Println("hello")
	if err != nil {
		t.Error(err)
	}
	defer c.Close()
	result := []db.User{}
	q := datastore.NewQuery("User").Filter("email =", "erik@bjareho.lt")
	var users []db.User
	_, err = q.GetAll(c, &users)
	if err != nil {
		t.Error(err)
	}
	if len(users) == 0 {
		t.Error("Couldn't find Erik in database")
	} else if len(result) > 1 {
		t.Error("More than one user with email erik@bjareho.lt in database")
	}
}

/*
func BenchmarkUserExistanceCycle(b *testing.B) {
	c, err := aetest.NewContext(nil)
	if err != nil {
		panic(err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		user := db.NewUser("tester", "password", "Tester", "test@example.com", []string{})
		userkey, err := datastore.Put(c, datastore.NewIncompleteKey(c, "User", nil), user)
		if err != nil {
			b.Error("Error when creating user")
			log.Println("I'm here")
		}
		datastore.Delete(c, userkey)
	}
}
*/

func TestNotDone(t *testing.T) {
	t.Skip("Not implemented")
}
