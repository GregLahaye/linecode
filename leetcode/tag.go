package leetcode

import (
	"encoding/json"
	"github.com/GregLahaye/linecode/linecode"
	"github.com/GregLahaye/linecode/store"
)

const tagsFilename = "tags.json"

func GetTags() ([]linecode.Tag, error) {
	var tags []linecode.Tag
	if err := store.ReadFromCache(&tags, tagsFilename); err == nil {
		return tags, nil
	}

	tags, err := FetchTags()
	if err != nil {
		return tags, err
	}

	err = store.SaveToCache(tags, tagsFilename)

	return tags, err
}

func FetchTags() ([]linecode.Tag, error) {
	body, err := request("GET", "/problems/api/tags/", nil)
	if err != nil {
		return nil, err
	}

	v := struct {
		Tags []linecode.Tag `json:"topics"`
	}{}
	if err = json.Unmarshal(body, &v); err != nil {
		return nil, err
	}

	return v.Tags, nil
}
