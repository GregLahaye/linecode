package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"golang.org/x/net/websocket"
	"io"
	"log"
	"os"
	"os/exec"
	"path"
	"regexp"
)

type Chrome struct {
	Cmd       *exec.Cmd
	WS        *websocket.Conn
	MessageID int
	TargetID  string
	SessionID string
}

type Message struct {
	ID     int             `json:"id"`
	Result json.RawMessage `json:"result"`
	Error  json.RawMessage `json:"error"`
	Method string          `json:"method"`
	Params json.RawMessage `json:"params"`
}

type Response struct {
	RequestID string `json:"requestId"`
}

type Request struct {
	RequestID string `json:"requestId"`
	Request   struct {
		URL    string `json:"url"`
		Method string `json:"method"`
	} `json:"request"`
}

type Cookie struct {
	Name     string  `json:"name"`
	Value    string  `json:"value"`
	Domain   string  `json:"domain"`
	Path     string  `json:"path"`
	Expires  float64 `json:"expires"`
	Size     int     `json:"size"`
	HTTPOnly bool    `json:"httpOnly"`
	Secure   bool    `json:"secure"`
	Session  bool    `json:"session"`
	SameSite string  `json:"sameSite,omitempty"`
}

type s map[string]interface{}

func create() (string, error) {
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

func start() (*Chrome, error) {
	dir, err := create()
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

	url, err := wait(pipe)
	if err != nil {
		return nil, err
	}

	ws, err := websocket.Dial(url, "", "http://localhost/")
	if err != nil {
		return nil, err
	}

	return &Chrome{Cmd: cmd, WS: ws}, nil
}

func wait(rd io.ReadCloser) (string, error) {
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

func (c *Chrome) connect() error {
	if err := c.Send("Target.setDiscoverTargets", s{"discover": true}); err != nil {
		return err
	}

	for {
		m := Message{}
		if err := websocket.JSON.Receive(c.WS, &m); err != nil {
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
				c.TargetID = v.TargetInfo.ID
				break
			}
		}
	}

	if err := c.Send("Target.attachToTarget", s{"targetId": c.TargetID}); err != nil {
		return err
	}

	for {
		m := Message{}
		if err := websocket.JSON.Receive(c.WS, &m); err != nil {
			return err
		} else if m.ID == c.MessageID {
			v := struct {
				ID string `json:"sessionId"`
			}{}
			if err := json.Unmarshal(m.Result, &v); err != nil {
				return err
			} else {
				c.SessionID = v.ID
				break
			}
		}
	}

	if err := c.SendMessageToTarget("Network.enable", nil); err != nil {
		return err
	}

	return nil
}

func (c *Chrome) login() error {
	var requestID string
	for {
		m := Message{}
		if err := websocket.JSON.Receive(c.WS, &m); err != nil {
			return err
		}

		v := struct {
			Message string `json:"message"`
		}{}
		json.Unmarshal(m.Params, &v)
		json.Unmarshal([]byte(v.Message), &m)

		if m.Method == "Network.requestWillBeSent" {
			r := Request{}
			json.Unmarshal(m.Params, &r)
			if r.Request.Method == "POST" && r.Request.URL == "https://leetcode.com/accounts/login/" {
				requestID = r.RequestID
			}
		} else if m.Method == "Network.responseReceived" {
			r := Response{}
			json.Unmarshal(m.Params, &r)
			if r.RequestID == requestID {
				return nil
			}
		}
	}
}

func (c *Chrome) end() {
	c.WS.Close()
	c.Send("Browser.close", nil)
}

func main() {
	c, err := start()
	if err != nil {
		log.Fatal(err)
	}
	defer c.end()

	if err = c.connect(); err != nil {
		log.Fatal(err)
	}

	LeetCodeSession, CSRFToken, err := c.getTokens()
	if err != nil {
		log.Fatal(err)
	}

	if LeetCodeSession != "" && CSRFToken != "" {
		fmt.Println(LeetCodeSession, CSRFToken)
		return
	}

	if err = c.login(); err != nil {
		log.Fatal(err)
	}

	LeetCodeSession, CSRFToken, err = c.getTokens()
	if err != nil {
		log.Fatal(err)
	}

	if LeetCodeSession != "" && CSRFToken != "" {
		fmt.Println(LeetCodeSession, CSRFToken)
	} else {
		fmt.Println("Could not retrieve cookies")
	}
}

func (c *Chrome) getTokens() (string, string, error) {
	cookies, err := c.getCookies()
	if err != nil {
		return "", "", err
	}

	var LeetCodeSession string
	var CSRFToken string

	for _, cookie := range cookies {
		if cookie.Name == "LEETCODE_SESSION" {
			LeetCodeSession = cookie.Value
		} else if cookie.Name == "csrftoken" {
			CSRFToken = cookie.Value
		}
	}

	return LeetCodeSession, CSRFToken, nil
}

func (c *Chrome) getCookies() ([]Cookie, error) {
	if err := c.Send("Storage.getCookies", []string{".leetcode.com"}); err != nil {
		return nil, err
	}

	for {
		m := Message{}
		websocket.JSON.Receive(c.WS, &m)
		if m.ID == c.MessageID {
			v := struct {
				Cookies []Cookie `json:"Cookies"`
			}{}
			json.Unmarshal(m.Result, &v)

			return v.Cookies, nil
		}
	}
}

func (c *Chrome) Send(method string, params interface{}) error {
	c.MessageID++
	return websocket.JSON.Send(c.WS, s{"id": c.MessageID, "method": method, "params": params})
}

func (c *Chrome) SendMessageToTarget(method string, params interface{}) error {
	b, _ := json.Marshal(s{"id": c.MessageID + 1, "method": method, "params": params})
	return c.Send("Target.sendMessageToTarget", s{"message": string(b), "sessionId": c.SessionID})
}

func PrettyPrint(v interface{}) error {
	b, err := json.MarshalIndent(v, "", "  ")

	if err == nil {
		fmt.Println(string(b))
	}

	return err
}
