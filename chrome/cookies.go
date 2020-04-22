package chrome

import (
	"encoding/json"
	"golang.org/x/net/websocket"
)

func (c *chrome) findTokens() (string, string, error) {
	var SessionID, CSRFToken string

	cookies, err := c.getCookies()
	if err != nil {
		return "", "", err
	}

	for _, ck := range cookies {
		if ck.Name == "LEETCODE_SESSION" {
			SessionID = ck.Value
		} else if ck.Name == "csrftoken" {
			CSRFToken = ck.Value
		}
	}

	if SessionID == "" || CSRFToken == "" {
		return SessionID, CSRFToken, tokenNotFoundError
	}

	return SessionID, CSRFToken, nil
}

func (c *chrome) getCookies() ([]cookie, error) {
	if err := c.send("Storage.getCookies", []string{".leetcode.com"}); err != nil {
		return nil, err
	}

	for {
		m := message{}
		websocket.JSON.Receive(c.ws, &m)
		if m.ID == c.messageID {
			v := struct {
				Cookies []cookie `json:"Cookies"`
			}{}
			json.Unmarshal(m.Result, &v)

			return v.Cookies, nil
		}
	}
}
