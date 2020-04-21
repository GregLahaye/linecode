package main

import (
	"encoding/json"
	"github.com/GregLahaye/yoyo"
	"github.com/GregLahaye/yoyo/styles"
)

type Tag struct {
	Slug      string `json:"slug"`
	Name      string `json:"name"`
	Questions []int  `json:"questions"`
}

const tagsFilename = "tags.json"

func (u *User) GetTags() ([]Tag, error) {
	var tags []Tag

	if err := CacheRetrieve(tagsFilename, &tags); err != nil {
		tags, err = u.DownloadTags()
		if err != nil {
			return tags, err
		}

		if err = CacheStore(tagsFilename, tags); err != nil {
			return tags, err
		}
	}

	return tags, nil
}

func (u *User) DownloadTags() ([]Tag, error) {
	s := yoyo.Start(styles.Point)
	defer s.End()

	body, err := u.Request("GET", baseUrl+"/problems/api/tags/", nil)
	if err != nil {
		return nil, err
	}

	v := struct {
		Tags []Tag `json:"topics"`
	}{}
	if err = json.Unmarshal(body, &v); err != nil {
		return nil, err
	}

	return v.Tags, nil
}
