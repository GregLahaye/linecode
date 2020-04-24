package leetcode

import (
	"encoding/json"
	"fmt"
	"github.com/GregLahaye/convert"
	"github.com/GregLahaye/linecode/config"
	"github.com/GregLahaye/linecode/store"
)

func Star(arg string) error {
	hash, err := GetHash()
	if err != nil {
		return err
	}

	id, _, err := Search(arg)
	if err != nil {
		return err
	}

	data := dict{"favorite_id_hash": hash, "question_id": id}
	_, err = request("POST", "/list/api/questions", data)
	if err != nil {
		return err
	}

	return store.RemoveFromCache(problemsFilename)
}

func Unstar(arg string) error {
	hash, err := GetHash()
	if err != nil {
		return err
	}

	id, _, err := Search(arg)
	if err != nil {
		return err
	}

	_, err = request("DELETE", "/list/api/questions/"+hash+"/"+convert.IntToString(id), nil)
	if err != nil {
		return err
	}

	return store.RemoveFromCache(problemsFilename)
}

func GetHash() (string, error) {
	c, err := config.Config()
	if err == nil && c.Hash != "" {
		return c.Hash, nil
	}

	return FetchHash()
}

func FetchHash() (string, error) {
	body, err := request("GET", "/list/api/questions", nil)
	if err != nil {
		return "", err
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
		return "", err
	}

	for _, i := range f.Favorites.Private {
		if i.Name == "Favorite" {
			return i.Hash, nil
		}
	}

	return "", fmt.Errorf("could not find hash")
}
