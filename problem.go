package main

import (
	"encoding/json"
	"github.com/GregLahaye/yoyo"
	"github.com/GregLahaye/yoyo/styles"
)

const problemsFilename = "problems.json"

func (u *User) GetProblems() ([]Problem, error) {
	var problems []Problem

	if err := CacheRetrieve(problemsFilename, &problems); err != nil {
		problems, err = u.DownloadProblems()
		if err != nil {
			return problems, err
		}

		if err = CacheStore(problemsFilename, problems); err != nil {
			return problems, err
		}
	}

	return problems, nil
}

func (u *User) DownloadProblems() ([]Problem, error) {
	var problems []Problem

	s := yoyo.Start(styles.Point)
	defer s.End()

	body, err := u.Request("GET", baseUrl+"/api/problems/algorithms/", nil)
	if err != nil {
		return problems, err
	}

	v := struct {
		Problems []Problem `json:"stat_status_pairs"`
	}{}
	if err = json.Unmarshal(body, &v); err != nil {
		return problems, err
	}
	problems = v.Problems

	return problems, nil
}
