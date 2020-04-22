package leetcode

import (
	"encoding/json"
	"github.com/GregLahaye/linecode/linecode"
	"github.com/GregLahaye/linecode/store"
)

var problemsFilename = "problems.json"

func GetProblems() ([]linecode.Problem, error) {
	var problems []linecode.Problem
	if err := store.ReadFromCache(&problems, problemsFilename); err == nil {
		return problems, nil
	}

	if problems, err := FetchProblems(); err != nil {
		return problems, err
	}

	err := store.SaveToCache(problems, problemsFilename)

	return problems, err
}

func FetchProblems() ([]linecode.Problem, error) {
	var problems []linecode.Problem

	body, err := request("GET", "/api/problems/algorithms/", nil)
	if err != nil {
		return problems, err
	}

	v := struct {
		Problems []linecode.Problem `json:"stat_status_pairs"`
	}{}
	if err = json.Unmarshal(body, &v); err != nil {
		return problems, err
	}
	problems = v.Problems

	return problems, nil
}

