package api

import (
	//	"gopkg.in/mgo.v2"
	"github.com/ErikBjare/Futarchio/src/db"
	"gopkg.in/mgo.v2/bson"
	"log"
	"testing"
)

func TestUsers(t *testing.T) {
	c := Users.collection
	result := []db.User{}
	c.Find(bson.M{"email": "erik@bjareho.lt"}).All(&result)
	if len(result) == 0 {
		t.Error("Couldn't find Erik in database")
	} else if len(result) > 1 {
		t.Error("More than one user with email erik@bjareho.lt in database")
	}
}

func BenchmarkUserExistanceCycle(b *testing.B) {
	c := Users.collection
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		user := db.NewUser("tester", "password", "Tester", "test@example.com", []string{})
		err := c.Insert(user)
		if err != nil {
			b.Error("Error when creating user")
			log.Println("I'm here")
		}
		c.Remove(user)
	}
}

func TestNotDone(t *testing.T) {
	t.Skip("Not implemented")
}
