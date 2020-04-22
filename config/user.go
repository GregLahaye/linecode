package config

import (
	"github.com/GregLahaye/linecode/chrome"
	"github.com/GregLahaye/linecode/store"
)

type User struct {
	Language string
	SessionID string
	CSRFToken string
}

const userFilename = "config.json"

func Config() (User, error) {
	u := User{}
	if err := store.ReadFromConfig(&u, userFilename); err != nil {
		return setup()
	}
	return u, nil
}

func setup() (User, error) {
	u := User{}

	SessionID, CSRFToken, err := chrome.RetrieveCredentials()
	if err != nil {
		return u, err
	}

	u.Language = "python3"
	u.SessionID = SessionID
	u.CSRFToken = CSRFToken

	err = store.SaveToConfig(u, userFilename)

	return u, err
}
