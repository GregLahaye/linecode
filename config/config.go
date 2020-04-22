package config

import (
	"fmt"
	"github.com/GregLahaye/linecode/config/chrome"
	"github.com/GregLahaye/linecode/convert"
	"github.com/GregLahaye/linecode/linecode"
	"github.com/GregLahaye/linecode/store"
)

type User struct {
	Language  string
	Hash      string
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

type AAA []fmt.Stringer

func setup() (User, error) {
	u := User{}

	SessionID, CSRFToken, err := chrome.RetrieveCredentials()
	if err != nil {
		return u, err
	}

	language := selectLanguage()

	u.Language = language
	u.SessionID = SessionID
	u.CSRFToken = CSRFToken

	err = store.SaveToConfig(u, userFilename)

	return u, err
}

func selectLanguage() string {
	var s []string
	for _, l := range linecode.Languages {
		s = append(s, l.String())
	}

	i := convert.Select(s)

	return linecode.Languages[i].Slug
}
