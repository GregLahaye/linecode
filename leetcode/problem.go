package leetcode

import (
	"encoding/json"
	"fmt"
	"github.com/GregLahaye/input"
	"github.com/GregLahaye/linecode/config"
	"github.com/GregLahaye/linecode/linecode"
	"github.com/GregLahaye/linecode/store"
	"net/url"
	"strconv"
)

var problemsFilename = "problems.json"

func GetProblems() ([]linecode.Problem, error) {
	var problems []linecode.Problem
	if err := store.ReadFromCache(&problems, problemsFilename); err == nil {
		return problems, nil
	}

	problems, err := FetchProblems()
	if err != nil {
		return problems, err
	}

	err = store.SaveToCache(problems, problemsFilename)

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

func Search(arg string) (int, string, error) {
	if id, err := strconv.Atoi(arg); err == nil {
		if problem, err := findByID(id); err == nil {
			return problem.Stat.ID, problem.Stat.Slug, err
		}
	}

	if problem, err := findBySlug(arg); err == nil {
		return problem.Stat.ID, problem.Stat.Slug, nil
	}

	if problem, err := findByFilename(arg); err == nil {
		return problem.Stat.ID, problem.Stat.Slug, nil
	}

	if problem, err := findByQuery(arg); err == nil {
		return problem.Stat.ID, problem.Stat.Slug, nil
	}

	return 0, "", fmt.Errorf("problem not found")
}

func findByID(id int) (linecode.Problem, error) {
	problems, err := GetProblems()
	if err != nil {
		return linecode.Problem{}, err
	}

	for _, problem := range problems {
		if problem.Stat.ID == id {
			return problem, nil
		}
	}

	return linecode.Problem{}, fmt.Errorf("problem not found")
}

func findBySlug(slug string) (linecode.Problem, error) {
	problems, err := GetProblems()
	if err != nil {
		return linecode.Problem{}, err
	}

	for _, problem := range problems {
		if problem.Stat.Slug == slug {
			return problem, nil
		}
	}

	return linecode.Problem{}, fmt.Errorf("problem not found")
}

func findByFilename(filename string) (linecode.Problem, error) {
	var problem linecode.Problem

	i, slug, err := parseFilename(filename)
	if err != nil {
		return problem, err
	}

	if id, err := strconv.Atoi(i); err == nil {
		if problem, err = findByID(id); err == nil {
			return problem, nil
		}
	}

	if problem, err = findBySlug(slug); err == nil {
		return problem, nil
	}

	return problem, err
}

func findByQuery(query string) (linecode.Problem, error) {
	u := "/problems/api/filter-questions/" + url.PathEscape(query)
	body, err := request("GET", u, nil)
	if err != nil {
		return linecode.Problem{}, err
	}

	var ids []int
	if err := json.Unmarshal(body, &ids); err != nil {
		return linecode.Problem{}, err
	}

	m := len(ids)
	if m > 10 {
		m = 10
	}

	var problems []linecode.Problem
	for _, id := range ids[:m] {
		if problem, err := findByID(id); err == nil {
			problems = append(problems, problem)
		}
	}

	if len(ids) < 1 {
		return linecode.Problem{}, fmt.Errorf("problem not found")
	} else if len(problems) == 1 {
		return problems[0], nil
	}

	var s []string
	for _, p := range problems {
		s = append(s, p.String())
	}

	i := input.Select(s)

	return problems[i], nil
}

func FindFile(arg string) string {
	if store.DoesExist(arg) {
		return arg
	}

	id, slug, err := Search(arg)
	if err != nil {
		return arg
	}

	c, err := config.Config()
	if err != nil {
		return arg
	}

	l := linecode.FindLanguage(c.Language)
	filename := fmt.Sprintf("%d.%s.%s", id, slug, l.Extension)
	return filename
}
