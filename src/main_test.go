package main

// TODO: Move relevant tests to API package

import (
	//	"gopkg.in/mgo.v2"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
)

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
	fmt.Println(string(body))
	json.Unmarshal(body, &msg)

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
		data := map[string]interface{}{}
		err = json.Unmarshal(body, &data)
		if err != nil {
			t.Error(err)
		}
		if data["length"] == float64(0) {
			t.Fatal(fmt.Sprintf("Got zero results: %f", data["length"]))
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
