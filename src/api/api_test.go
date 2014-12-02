package api

import (
	//	"gopkg.in/mgo.v2"
	"appengine/aetest"
	"appengine/datastore"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ErikBjare/Futarchio/src/db"
	"io/ioutil"
	"log"
	"net/http"
	"testing"
	"time"
)

func init() {
	_, err := http.Get("http://localhost:8080/api/0/init")
	if err != nil {
		log.Fatal(err)
	}

	// Required to let datastore init
	time.Sleep(time.Second)
}

func getAuthkey() (string, error) {
	client := &http.Client{}
	bodybuf := &bytes.Buffer{}
	bodybuf.Write([]byte("{\"username\": \"erb\", \"password\": \"password\"}"))
	resp, err := client.Post("http://localhost:8080/api/0/auth", "application/json", bodybuf)
	if err != nil {
		return "", err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	msg := db.Auth{}
	err = json.Unmarshal(body, &msg)
	if err != nil {
		return "", err
	}

	return msg.Key, nil
}

type UserWithKey struct {
	db.User
	Key *datastore.Key `json:"key"`
}

func TestAuth(t *testing.T) {
	authkey, err := getAuthkey()
	if err != nil {
		t.Fatal(err)
	}

	// Apparently required to allow the datastore time to be able to store Auth
	// TODO: Should probably be removed once Memcache is implemented
	time.Sleep(time.Second)

	body, err := getBody("http://localhost:8080/api/0/users/me", authkey)
	if err != nil {
		t.Fatal(err)
	}

	var user UserWithKey
	err = json.Unmarshal(body, &user)
	if err != nil {
		t.Error(err)
	}

	body, err = getBody("http://localhost:8080/api/0/users?key="+user.Key.Encode(), authkey)
	if err != nil {
		t.Fatal(err)
	}

	var users []UserResponse
	err = json.Unmarshal(body, &users)
	if err != nil {
		t.Error(err)
	}

	if user.Username != users[0].Username {
		t.Error("usernames did not match")
	}
}

func getBody(url string, authkey string) ([]byte, error) {
	client := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	if authkey != "" {
		req.Header.Add("Authorization", authkey)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		err := errors.New(fmt.Sprintf("Status code was not 200, was %d with message: %s", resp.StatusCode, resp.Status))
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, err
}

func TestUsers(t *testing.T) {
	c, err := aetest.NewContext(nil)
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()

	key := datastore.NewKey(c, "User", "erb", 0, nil)
	user1 := db.NewUser("erb", "secretpassword", "Erik", "erik@bjareho.lt")
	_, err = datastore.Put(c, key, &user1)
	if err != nil {
		t.Fatal(err)
	}

	var user2 db.User
	err = datastore.Get(c, key, &user2)
	if err != nil {
		t.Error(err)
	}

	var users []db.User
	q := datastore.NewQuery("User").Filter("Email =", "erik@bjareho.lt")
	keys, err := q.GetAll(c, &users)
	if err != nil {
		t.Error(err)
	}

	if len(keys) == 0 {
		t.Error("Couldn't find Erik in database")
	} else if len(keys) > 1 {
		t.Error("More than one user with email erik@bjareho.lt in database")
	}
}

func TestPolls(t *testing.T) {
	_, err := getBody("http://localhost:8080/api/0/polls", "")
	if err != nil {
		t.Fatal(err)
	}
}
