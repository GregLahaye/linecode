package main

import (
	"encoding/json"
	"github.com/GregLahaye/yoyo"
	"github.com/GregLahaye/yoyo/styles"
)

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
