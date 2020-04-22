package leetcode

import "encoding/json"

func Star(id, hash string) error {
	data := dict{"favorite_id_hash": hash, "question_id": id}
	_, err := request("POST", "/list/api/questions", data)
	if err != nil {
		return err
	}

	return CacheDestroy(problemsFilename)
}

func Unstar(id, hash string) error {
	_, err := u.Request("DELETE", "/list/api/questions/" + hash + "/" + id, nil)
	if err != nil {
		return err
	}

	return CacheDestroy(problemsFilename)
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
}
