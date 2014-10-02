package main

import (
	//	"gopkg.in/mgo.v2"
	"encoding/json"
	"fmt"
	"gopkg.in/mgo.v2/bson"
	"io/ioutil"
	"log"
	"net/http"
	"testing"
)

func TestUsers(t *testing.T) {
	u := NewUserResource(NewSession())
	u.Init()
	c := u.collection
	result := []User{}
	c.Find(bson.M{"email": "erik@bjareho.lt"}).All(&result)
	if len(result) == 0 {
		t.Error("Couldn't find Erik in database")
	} else if len(result) > 1 {
		t.Error("More than one user with email erik@bjareho.lt in database")
	}

}

func BenchmarkUserExistanceCycle(b *testing.B) {
	u := NewUserResource(NewSession())
	c := u.collection
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		user := NewUser("tester", "Tester", "test@example.com", []string{})
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

func BenchmarkAPICall(b *testing.B) {
	client := &http.Client{}
	for i := 0; i < b.N; i++ {
		req, err := http.NewRequest("GET", "http://localhost:8080/api/0/users", nil)
		if err != nil {
			b.Fatal(err)
		}
		req.Header.Add("Authorization", "Basic YWRtaW46YWRtaW4=")
		resp, err := client.Do(req)
		if err != nil {
			b.Error("Is the server running?")
			b.Fatal(err)
		}
		defer resp.Body.Close()
		if resp.StatusCode != 200 {
			b.Fatal(fmt.Sprintf("Status code was not 200, was %d with message: %s", resp.StatusCode, resp.Status))
		}

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			b.Error(err)
		}
		data := map[string]interface{}{}
		err = json.Unmarshal(body, &data)
		if err != nil {
			b.Error(err)
		}
		if data["length"] == float64(0) {
			b.Fatal(fmt.Sprintf("Got zero results: %f", data["length"]))
		}
	}
}
