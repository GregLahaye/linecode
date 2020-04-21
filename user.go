package main

import "fmt"

type User struct {
	Language       Language `json:"language"`
	Editor         string   `json:"editor"`
	TerminalEditor bool     `json:"terminal_editor"`
	Hash           string   `json:"favorites_hash"`
	Credentials    struct {
		Session   string `json:"session"`
		CSRFToken string `json:"csrf_token"`
	} `json:"credentials"`
}

const userFilename = "user.json"

func LoadUser() (User, error) {
	var u User

	if err := CacheRetrieve(userFilename, &u); err != nil {
		if err := u.Login(); err != nil {
			fmt.Println("Couldn't retrieve credentials from Chrome")
			if Confirm("Would you like to manually enter credentials? (Y/N) ") {
				fmt.Print("LEETCODE_SESSION: ")
				if u.Credentials.Session, err = StringInput(); err != nil {
					return u, err
				}

				fmt.Print("CSRF Token: ")
				if u.Credentials.CSRFToken, err = StringInput(); err != nil {
					return u, err
				}
			} else {
				return u, err
			}
		}

		u.Language = SelectLanguage()

		fmt.Print("Default editor: ")
		if u.Editor, err = StringInput(); err != nil {
			return u, err
		}

		u.TerminalEditor = Confirm("Is this a terminal editor (e.g. vim)? ")

		if err = u.FindFavorites(); err != nil {
			return u, err
		}

		if err = CacheStore(userFilename, u); err != nil {
			return u, err
		}
	}

	return u, nil
}
