package leetcode

import (
	"encoding/json"
	"github.com/GregLahaye/linecode/linecode"
)

func FetchTags() ([]linecode.Tag, error) {
	body, err := u.Request("GET", "/problems/api/tags/", nil)
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
