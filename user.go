package main

type User struct {
	Language    Language `json:"language"`
	Credentials struct {
		Session   string `json:"session"`
		CSRFToken string `json:"csrf_token"`
	} `json:"credentials"`
}

func LoadUser() (User, error) {
	var u User
	if err := LoadStruct(userFilename, &u); err != nil {
		if err := u.Login(); err != nil {
			return u, err
		}

		u.Language = SelectLanguage()

		if err = SaveStruct(userFilename, u); err != nil {
			return u, err
		}
	}

	return u, nil
}
