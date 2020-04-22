package main

import (
	"encoding/json"
	"github.com/GregLahaye/yoyo"
	"github.com/GregLahaye/yoyo/styles"
)

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
