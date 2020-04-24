package leetcode

import (
	"encoding/json"
	"fmt"
	"github.com/GregLahaye/linecode/config"
	"github.com/GregLahaye/linecode/convert"
	"github.com/GregLahaye/linecode/linecode"
	"github.com/GregLahaye/linecode/store"
	"path"
	"strings"
)

var questionDirectory = "questions"

func GetQuestion(arg string) (linecode.Question, error) {
	var question linecode.Question

	id, slug, err := Search(arg)
	if err != nil {
		return linecode.Question{}, err
	}

	filename := convert.IntToString(id) + "." + slug + ".json"
	p := path.Join(questionDirectory, filename)
	if err := store.ReadFromCache(&question, p); err == nil {
		return question, nil
	}

	question, err = FetchQuestion(slug)
	if err != nil {
		return question, err
	}

	err = store.SaveToCache(question, p)

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

func SaveSnippet(q linecode.Question) error {
	c, _ := config.Config()

	l := linecode.FindLanguage(c.Language)

	for _, s := range q.CodeSnippets {
		if s.Slug == c.Language {
			filename := fmt.Sprintf("%s.%s.%s", q.ID, q.Slug, l.Extension)
			snippet := createSnippet(convert.ParseHTML(q.Content), s.Code, l.Comment)
			if store.DoesNotExist(filename) {
				return store.WriteFile(snippet, filename)
			}
		}
	}

	return fmt.Errorf("could not find language")
}

func createSnippet(content, code string, comment linecode.Comment) string {
	var s strings.Builder

	s.WriteString(comment.Start)
	s.WriteString("\n")
	s.WriteString(content)
	s.WriteString("\n")
	s.WriteString(comment.End)
	s.WriteString("\n\n")
	s.WriteString(code)

	return s.String()
}
