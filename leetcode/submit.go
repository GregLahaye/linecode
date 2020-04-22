package leetcode

import (
	"encoding/json"
	"github.com/GregLahaye/linecode/config"
	"github.com/GregLahaye/linecode/convert"
	"github.com/GregLahaye/linecode/linecode"
	"github.com/GregLahaye/linecode/store"
)

func SubmitCode(filename string) (linecode.Submission, error) {
	var submission linecode.Submission

	id, slug, err := parseFilename(filename)
	if err != nil {
		return submission, err
	}

	c, _ := config.Config()
	language := c.Language

	code, err := store.ReadFile(filename)
	if err != nil {
		return submission, err
	}

	clearQuestion(id, slug)

	return submitCode(id, slug, language, code)
}

func submitCode(id, slug, language, code string) (linecode.Submission, error) {
	var submission linecode.Submission

	data := dict{"lang": language, "question_id": id, "typed_code": code}
	body, err := request("POST", "/problems/"+slug+"/submit/", data)
	if err != nil {
		return submission, err
	}

	v := struct {
		SubmissionID int `json:"submission_id"`
	}{}
	if err = json.Unmarshal(body, &v); err != nil {
		return submission, err
	}

	return retry(convert.IntToString(v.SubmissionID))
}
