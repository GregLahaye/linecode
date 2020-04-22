package chrome

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

var tokenNotFoundError = errors.New("could not find tokens")

type dict map[string]interface{}

func RetrieveCredentials() (string, string, error) {
	var SessionID, CSRFToken string

	c, err := open()
	if err != nil {
		return SessionID, CSRFToken, err
	}
	defer c.close()

	if err = c.getSessionID(); err != nil {
		return SessionID, CSRFToken, err
	}

	SessionID, CSRFToken, err = c.findTokens()
	if err == tokenNotFoundError {
		if err = c.waitForLogin(); err != nil {
			return SessionID, CSRFToken, err
		}

		SessionID, CSRFToken, err = c.findTokens()
		if err != nil {
			return SessionID, CSRFToken, err
		}
	} else if err != nil {
		return SessionID, CSRFToken, err
	}

	return SessionID, CSRFToken, nil
}

func open() (*chrome, error) {
	dir, err := getUserDataDir()
	if err != nil {
		return nil, err
	}

	args := []string{"https://leetcode.com/accounts/login/", "--remote-debugging-port=0", "--user-data-dir=" + dir}
	cmd := exec.Command(locate(), args...)

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

func (c *chrome) close() {
	c.send("Browser.close", nil)
	c.ws.Close()
}

func getUserDataDir() (string, error) {
	cfg, _ := os.UserConfigDir()
	dir := path.Join(cfg, "linecode", "chrome")

	err := os.MkdirAll(dir, os.ModePerm)
	return dir, err
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
