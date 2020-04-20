package main

import "fmt"

type User struct {
	Language    Language `json:"language"`
	Editor      string   `json:"editor"`
	Hash        string   `json:"favorites_hash"`
	Credentials struct {
		Session   string `json:"session"`
		CSRFToken string `json:"csrf_token"`
	} `json:"credentials"`
}

const userFilename = "user.json"

func LoadUser() (User, error) {
	var u User
	if err := LoadStruct(userFilename, &u); err != nil {
		if err := u.Login(); err != nil {
			return u, err
		}

		u.Language = SelectLanguage()

		fmt.Print("Default editor: ")
		if u.Editor, err = StringInput(); err != nil {
			return u, err
		}

		if err = u.FindFavorites(); err != nil {
			return u, err
		}

		if err = SaveStruct(userFilename, u); err != nil {
			return u, err
		}
	}

	return u, nil
}
