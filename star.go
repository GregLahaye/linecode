package main

import (
	"encoding/json"
)

func (u *User) Star(arg string) error {
	problem, err := u.FindProblem(arg)
	if err != nil {
		return err
	}

	data := dict{"favorite_id_hash": u.Hash, "question_id": problem.Stat.ID}
	_, err = u.Request("POST", baseUrl+"/list/api/questions", data)
	if err != nil {
		return err
	}

	return CacheDestroy(problemsFilename)
}

func (u *User) Unstar(arg string) error {
	problem, err := u.FindProblem(arg)
	if err != nil {
		return err
	}

	url := baseUrl + "/list/api/questions/" + u.Hash + "/" + IntToString(problem.Stat.ID)
	_, err = u.Request("DELETE", url, nil)
	if err != nil {
		return err
	}

	return CacheDestroy(problemsFilename)
}

func (u *User) FindFavorites() error {
	body, err := u.Request("GET", baseUrl+"/list/api/questions", nil)
	if err != nil {
		return err
	}

	f := struct {
		Favorites struct {
			Private []struct {
				Hash string `json:"id_hash"`
				Name string `json:"name"`
			} `json:"private_favorites"`
		} `json:"favorites"`
	}{}
	err = json.Unmarshal(body, &f)
	if err != nil {
		return err
	}

	for _, i := range f.Favorites.Private {
		if i.Name == "Favorite" {
			u.Hash = i.Hash
		}
	}

	return nil
}
