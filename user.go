package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type User struct {
	Language    Language `json:"language"`
	Credentials struct {
		Session   string `json:"session"`
		CSRFToken string `json:"csrf_token"`
	} `json:"credentials"`
}

func LoadUser() (User, error) {
	u, err := LoadUserFromFile()
	if err != nil {
		u = User{}
		if err := u.Login(); err != nil {
			return u, err
		}

		u.Language = SelectLanguage()

		if err = SaveUser(u); err != nil {
			return u, err
		}
	}

	return u, nil
}

func LoadUserFromFile() (User, error) {
	f, err := os.Open(userDataFilename)
	if err != nil {
		return User{}, err
	}

	b, err := ioutil.ReadAll(f)
	if err != nil {
		return User{}, err
	}

	var u User
	if err = json.Unmarshal(b, &u); err != nil {
		return User{}, err
	}

	return u, nil
}

func SaveUser(u User) error {
	b, err := json.MarshalIndent(u, "", "  ")
	if err != nil {
		return err
	}

	if err = ioutil.WriteFile(userDataFilename, b, os.ModePerm); err != nil {
		return err
	}

	return nil
}
