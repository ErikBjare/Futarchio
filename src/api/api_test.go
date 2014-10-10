package api

import (
	//	"gopkg.in/mgo.v2"
	"appengine/aetest"
	"appengine/datastore"
	"bytes"
	"encoding/json"
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
}

func TestAuth(t *testing.T) {
	client := &http.Client{}
	bodybuf := &bytes.Buffer{}
	bodybuf.Write([]byte("{\"username\": \"erb\", \"password\": \"password\"}"))
	resp, err := client.Post("http://localhost:8080/api/0/auth", "application/json", bodybuf)
	if err != nil {
		t.Fatal(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	msg := map[string]map[string]string{}
	//fmt.Println(string(body))
	err = json.Unmarshal(body, &msg)
	if err != nil {
		t.Fatal(err)
	}

	// Apparently required to allow the datastore time to be able to store Auth
	// Should probably be removed once Memcache is implemented
	time.Sleep(1 * time.Second)

	urls := []string{"http://localhost:8080/api/0/users", "http://localhost:8080/api/0/users/me"}
	for _, url := range urls {
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Add("Authorization", msg["auth"]["key"])
		resp, err = client.Do(req)
		if err != nil {
			t.Error("Is the server running?")
			t.Fatal(err)
		}
		defer resp.Body.Close()
		if resp.StatusCode != 200 {
			t.Fatal(fmt.Sprintf("Status code was not 200, was %d with message: %s", resp.StatusCode, resp.Status))
		}

		body, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Error(err)
		}

		// The following comment line can be useful for debugging by printing body
		//log.Println(string(body))

		data := map[string][]interface{}{}
		err = json.Unmarshal(body, &data)
		if err != nil {
			t.Error(err)
		}
		if len(data["data"]) == 0 {
			t.Fatal(fmt.Sprintf("Got zero results or non-array"))
		}
	}
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

func TestUsers(t *testing.T) {
	c, err := aetest.NewContext(nil)
	if err != nil {
		t.Error(err)
	}
	defer c.Close()

	key := datastore.NewKey(c, "User", "", 1, nil)
	_, err = datastore.Put(c, key, db.NewUser("erb", "secretpassword", "Erik", "erik@bjareho.lt", []string{}))
	if err != nil {
		t.Fatal(err)
	}

	var user db.User
	err = datastore.Get(c, key, &user)
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
