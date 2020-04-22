package chrome

import (
	"encoding/json"
	"github.com/GregLahaye/linecode/leetcode"
	"golang.org/x/net/websocket"
)

func (c *chrome) send(method string, params interface{}) error {
	c.messageID++
	return websocket.JSON.Send(c.ws, dict{"id": c.messageID, "method": method, "params": params})
}

func (c *chrome) sendToTarget(method string, params interface{}) error {
	b, _ := json.Marshal(dict{"id": c.messageID + 1, "method": method, "params": params})
	return c.send("Target.sendMessageToTarget", dict{"message": string(b), "sessionId": c.sessionID})
}

func (c *chrome) getSessionID() error {
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
			if r.Request.Method == "POST" && r.Request.URL == leetcode.BaseURL+"/accounts/login/" {
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
