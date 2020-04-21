package main

import (
	"encoding/json"
	"github.com/GregLahaye/yoyo"
	"github.com/GregLahaye/yoyo/styles"
)

type Question struct {
	QuestionID     string        `json:"questionId"`
	Title          string        `json:"title"`
	Slug           string        `json:"titleSlug"`
	Content        string        `json:"content"`
	IsPaidOnly     bool          `json:"isPaidOnly"`
	Difficulty     string        `json:"difficulty"`
	Tags           []Tag         `json:"topicTags"`
	CodeSnippets   []CodeSnippet `json:"codeSnippets"`
	Status         string        `json:"status"`
	SampleTestCase string        `json:"sampleTestCase"`
}

type CodeSnippet struct {
	Lang string `json:"lang"`
	Slug string `json:"langSlug"`
	Code string `json:"code"`
}

const questionsDirectory = "questions"

func (u *User) GetQuestion(arg string) (Question, error) {
	var question Question

	problem, err := u.FindProblem(arg)
	if err != nil {
		return question, err
	}

	filename := QuestionFilename(problem.Stat.ID)
	if err := CacheRetrieve(filename, &question); err != nil {
		question, err = u.DownloadQuestion(problem.Stat.Slug)
		if err != nil {
			return question, err
		}

		if err = CacheStore(filename, question); err != nil {
			return question, err
		}
	}

	return question, nil
}

func (u *User) DownloadQuestion(slug string) (Question, error) {
	var question Question

	s := yoyo.Start(styles.Point)
	defer s.End()

	data := dict{"variables": dict{"titleSlug": slug}, "operationName": "questionData", "query": "query questionData($titleSlug: String!) {\n  question(titleSlug: $titleSlug) {\n    questionId\n    title\n    titleSlug\n    content\n    isPaidOnly\n    difficulty\n    isLiked\n    topicTags {\n      name\n      slug\n    }\n    codeSnippets {\n      lang\n      langSlug\n      code\n    }\n    stats\n    status\n    sampleTestCase\n    }\n}"}
	body, err := u.Request("POST", baseUrl+"/graphql", data)
	if err != nil {
		return question, err
	}

	v := struct {
		Data struct {
			Question Question `json:"question"`
		} `json:"data"`
	}{}
	if err = json.Unmarshal(body, &v); err != nil {
		return question, err
	}

	return v.Data.Question, nil
}

func (u *User) DownloadAll() error {
	if err := CacheDestroy(problemsFilename); err != nil {
		return err
	}

	if err := CacheDestroy(questionsDirectory); err != nil {
		return err
	}

	problems, err := u.GetProblems()
	if err != nil {
		return err
	}

	for _, problem := range problems {
		if _, err := u.GetQuestion(problem.Stat.Slug); err != nil {
			return err
		} else {
			DisplayProblem(problem)
		}
	}

	return nil
}
