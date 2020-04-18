package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"golang.org/x/net/websocket"
	"io"
	"os"
	"os/exec"
	"path"
	"regexp"
)

type chrome struct {
	cmd       *exec.Cmd
	ws        *websocket.Conn
	messageID int
	targetID  string
	sessionID string
}

type message struct {
	ID     int             `json:"id"`
	Result json.RawMessage `json:"result"`
	Error  json.RawMessage `json:"error"`
	Method string          `json:"method"`
	Params json.RawMessage `json:"params"`
}

type response struct {
	RequestID string `json:"requestId"`
}

type request struct {
	RequestID string `json:"requestId"`
	Request   struct {
		URL    string `json:"url"`
		Method string `json:"method"`
	} `json:"request"`
}

type cookie struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

var TokensNotFound = errors.New("could not find tokens")

type dict map[string]interface{}

func Login(u *User) error {
	c, err := start()
	if err != nil {
		return err
	}
	defer c.close()

	if err = c.findSessionID(); err != nil {
		return err
	}

	LeetCodeSession, CSRFToken, err := c.findTokens()
	if err == TokensNotFound {
		if err = c.waitForLogin(); err != nil {
			return err
		}

		LeetCodeSession, CSRFToken, err = c.findTokens()
		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	}

	u.LeetCodeSession = LeetCodeSession
	u.CSRFToken = CSRFToken

	return nil
}

func start() (*chrome, error) {
	dir, err := createUserDirectory()
	if err != nil {
		return nil, err
	}

	args := []string{"https://leetcode.com/accounts/login/", "--remote-debugging-port=0", "--user-data-dir=" + dir}
	cmd := exec.Command("/Program Files (x86)/Google/Chrome/Application/chrome.exe", args...)

	pipe, err := cmd.StderrPipe()
	if err != nil {
		return nil, err
	}

	if err := cmd.Start(); err != nil {
		return nil, err
	}

	url, err := readWebSocketURL(pipe)
	if err != nil {
		return nil, err
	}

	ws, err := websocket.Dial(url, "", "http://localhost/")
	if err != nil {
		return nil, err
	}

	return &chrome{cmd: cmd, ws: ws}, nil
}

func createUserDirectory() (string, error) {
	dir, err := os.UserCacheDir()
	if err != nil {
		return "", err
	}

	dir = path.Join(dir, "leetcode-terminal", "chrome")

	err = os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return "", err
	}

	return dir, nil
}

func readWebSocketURL(rd io.ReadCloser) (string, error) {
	re := regexp.MustCompile("(ws://.*?)\r\n")
	br := bufio.NewReader(rd)
	for {
		if line, err := br.ReadString('\n'); err != nil {
			rd.Close()
			return "", err
		} else if m := re.FindStringSubmatch(line); m != nil {
			return m[1], nil
		}
	}
}

func (c *chrome) findSessionID() error {
	if err := c.send("Target.setDiscoverTargets", dict{"discover": true}); err != nil {
		return err
	}

	for {
		m := message{}
		if err := websocket.JSON.Receive(c.ws, &m); err != nil {
			return err
		} else if m.Method == "Target.targetCreated" {
			v := struct {
				TargetInfo struct {
					Type string `json:"type"`
					ID   string `json:"targetId"`
				} `json:"targetInfo"`
			}{}
			if err := json.Unmarshal(m.Params, &v); err != nil {
				return err
			} else if v.TargetInfo.Type == "page" {
				c.targetID = v.TargetInfo.ID
				break
			}
		}
	}

	if err := c.send("Target.attachToTarget", dict{"targetId": c.targetID}); err != nil {
		return err
	}

	for {
		m := message{}
		if err := websocket.JSON.Receive(c.ws, &m); err != nil {
			return err
		} else if m.ID == c.messageID {
			v := struct {
				ID string `json:"sessionId"`
			}{}
			if err := json.Unmarshal(m.Result, &v); err != nil {
				return err
			} else {
				c.sessionID = v.ID
				break
			}
		}
	}

	if err := c.sendToTarget("Network.enable", nil); err != nil {
		return err
	}

	return nil
}

func (c *chrome) waitForLogin() error {
	var requestID string
	for {
		m := message{}
		if err := websocket.JSON.Receive(c.ws, &m); err != nil {
			return err
		}

		v := struct {
			Message string `json:"message"`
		}{}
		json.Unmarshal(m.Params, &v)
		json.Unmarshal([]byte(v.Message), &m)

		if m.Method == "Network.requestWillBeSent" {
			r := request{}
			json.Unmarshal(m.Params, &r)
			if r.Request.Method == "POST" && r.Request.URL == "https://leetcode.com/accounts/login/" {
				requestID = r.RequestID
			}
		} else if m.Method == "Network.responseReceived" {
			r := response{}
			json.Unmarshal(m.Params, &r)
			if r.RequestID == requestID {
				return nil
			}
		}
	}
}
func (c *chrome) findTokens() (string, string, error) {
	cookies, err := c.getCookies()
	if err != nil {
		return "", "", err
	}

	var LeetCodeSession string
	var CSRFToken string

	for _, ck := range cookies {
		if ck.Name == "LEETCODE_SESSION" {
			LeetCodeSession = ck.Value
		} else if ck.Name == "csrftoken" {
			CSRFToken = ck.Value
		}
	}

	if LeetCodeSession == "" || CSRFToken == "" {
		return "", "", TokensNotFound
	}

	return LeetCodeSession, CSRFToken, nil
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

func (c *chrome) close() {
	c.send("Browser.close", nil)
	c.ws.Close()
}

func (c *chrome) send(method string, params interface{}) error {
	c.messageID++
	return websocket.JSON.Send(c.ws, dict{"id": c.messageID, "method": method, "params": params})
}

func (c *chrome) sendToTarget(method string, params interface{}) error {
	b, _ := json.Marshal(dict{"id": c.messageID + 1, "method": method, "params": params})
	return c.send("Target.sendMessageToTarget", dict{"message": string(b), "sessionId": c.sessionID})
}
