package leetcode

import (
	"encoding/json"
	"github.com/GregLahaye/linecode/linecode"
	"github.com/GregLahaye/linecode/store"
	"path"
)

var questionDirectory = "questions"

func GetQuestion() (linecode.Question, error) {
	var question linecode.Question

	p := path.Join(questionDirectory, "???")
	if err := store.ReadFromCache(&question, p); err == nil {
		return question, nil
	}

	if question, err := FetchQuestion("SLUG"); err != nil {
		return question, err
	}

	err := store.SaveToCache(question, p)

	return question, err
}

func FetchQuestion(slug string) (linecode.Question, error) {
	var question linecode.Question

	data := dict{"variables": dict{"titleSlug": slug}, "operationName": "questionData", "query": "query questionData($titleSlug: String!) {\n  question(titleSlug: $titleSlug) {\n    questionId\n    title\n    titleSlug\n    content\n    isPaidOnly\n    difficulty\n    isLiked\n    topicTags {\n      name\n      slug\n    }\n    codeSnippets {\n      lang\n      langSlug\n      code\n    }\n    stats\n    status\n    sampleTestCase\n    }\n}"}
	body, err := request("POST", "/graphql", data)
	if err != nil {
		return question, err
	}

	v := struct {
		Data struct {
			Question linecode.Question `json:"question"`
		} `json:"data"`
	}{}
	if err = json.Unmarshal(body, &v); err != nil {
		return question, err
	}

	return v.Data.Question, nil
}
