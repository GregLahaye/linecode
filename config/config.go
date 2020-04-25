package config

import (
	"fmt"
	"github.com/GregLahaye/input"
	"github.com/GregLahaye/linecode/config/chrome"
	"github.com/GregLahaye/linecode/linecode"
	"github.com/GregLahaye/linecode/store"
)

type User struct {
	Language  string
	Hash      string
	Editor    string
	Terminal  bool
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
		fmt.Print("Couldn't retrieve credentials from Chrome\n\n")
		fmt.Print("Please visit chrome://settings/cookies/detail?site=leetcode.com\n")
		fmt.Print(" and enter credentials manually\n\n")
		fmt.Print("SESSION_ID: ")
		SessionID = input.String()
		fmt.Print("\nCSRFToken: ")
		CSRFToken = input.String()
		fmt.Print("\n")
	}

	fmt.Println("Default language: ")
	language := selectLanguage()

	fmt.Print("\nDefault editor: ")
	editor := input.String()
	var terminal bool
	switch editor {
	case "vim", "emacs", "nano", "vi":
		terminal = true
	case "code", "sublime", "atom", "notepad++", "brackets", "notepad":
		terminal = false
	default:
		terminal = input.Confirm("Is this a terminal editor? (e.g. vim)")
	}

	u.Language = language
	u.Editor = editor
	u.Terminal = terminal
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

	i := input.Select(s)

	return linecode.Languages[i].Slug
}
