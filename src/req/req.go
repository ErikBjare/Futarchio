package req

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

func getJSON(url string) (map[string]interface{}, error) {
	body, err := getBytes(url)
	if err != nil {
		return nil, err
	}

	data := map[string]interface{}{}
	json.Unmarshal(body, &data)

	return data, nil
}

func getBytes(url string) ([]byte, error) {
	// Build the request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	// Send the request via a client
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	// Defer the closing of the body
	defer resp.Body.Close()
	// Read the content into a byte array
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// At this point we're done - simply return the bytes
	return body, nil
}
