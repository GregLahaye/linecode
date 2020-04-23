package leetcode

import (
	"bytes"
	"encoding/json"
	"github.com/GregLahaye/linecode/config"
	"github.com/GregLahaye/yoyo"
	"github.com/GregLahaye/yoyo/styles"
	"io/ioutil"
	"net/http"
)

type dict map[string]interface{}

const (
	BaseURL = "https://leetcode.com"
)

func request(method, path string, data dict) ([]byte, error) {
	s := yoyo.Start(styles.Bounce)
	defer s.End()

	client := &http.Client{}

	// create bytes from data
	b, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	// create request
	req, err := http.NewRequest(method, BaseURL+path, bytes.NewReader(b))
	if err != nil {
		return nil, err
	}

	// load config
	c, err := config.Config()
	if err != nil {
		return nil, err
	}

	// add cookies to request
	req.AddCookie(&http.Cookie{Name: "csrftoken", Value: c.CSRFToken, Domain: ".leetcode.com"})
	req.AddCookie(&http.Cookie{Name: "LEETCODE_SESSION", Value: c.SessionID, Domain: ".leetcode.com"})

	// add headers
	req.Header.Set("X-CSRFToken", c.CSRFToken)
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	req.Header.Set("Referer", BaseURL)
	req.Header.Set("Content-Type", "application/json")

	// make request
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// read body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
