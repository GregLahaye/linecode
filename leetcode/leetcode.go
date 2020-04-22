package leetcode

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type dict map[string]interface{}

const (
	BaseURL = "https://leetcode.com"
)

var (
	sessionID string
	csrfToken string
)

func request(method, path string, data dict) ([]byte, error) {
	client := &http.Client{}

	b, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(method, BaseURL+path, bytes.NewReader(b))
	if err != nil {
		return nil, err
	}

	req.AddCookie(&http.Cookie{Name: "csrftoken", Value: csrfToken, Domain: ".leetcode.com"})
	req.AddCookie(&http.Cookie{Name: "LEETCODE_SESSION", Value: sessionID, Domain: ".leetcode.com"})

	req.Header.Set("X-csrfToken", csrfToken)
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	req.Header.Set("Referer", BaseURL)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

