package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
)

func (u *User) Request(method, url string, data dict) ([]byte, error) {
	client := &http.Client{}

	b, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(method, url, bytes.NewReader(b))
	if err != nil {
		return nil, err
	}

	req.AddCookie(&http.Cookie{Name: "csrftoken", Value: u.Credentials.CSRFToken, Domain: ".leetcode.com"})
	req.AddCookie(&http.Cookie{Name: "LEETCODE_SESSION", Value: u.Credentials.Session, Domain: ".leetcode.com"})

	req.Header.Set("X-CSRFToken", u.Credentials.CSRFToken)
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	req.Header.Set("Referer", baseUrl+"/")
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func (u *User) FindProblem(arg string) (Problem, error) {
	var problem Problem

	problems, err := u.GetProblems()
	if err != nil {
		return problem, err
	}

	if id, err := strconv.Atoi(arg); err == nil {
		for _, p := range problems {
			if p.Stat.ID == id {
				return p, nil
			}
		}
	}

	for _, p := range problems {
		if p.Stat.Slug == arg {
			return p, nil
		}
	}

	slug := url.PathEscape(arg)
	body, err := u.Request("GET", baseUrl+"/problems/api/filter-questions/"+slug, nil)
	if err != nil {
		return problem, nil
	}

	var ids []int
	if err = json.Unmarshal(body, &ids); err != nil {
		return problem, err
	}

	length := len(ids)
	if length > 10 {
		length = 10
	}
	return u.SelectQuestion(ids[:length])
}

func (u *User) GetSlug(id int) (string, error) {
	problems, err := u.GetProblems()
	if err != nil {
		return "", err
	}

	for _, p := range problems {
		if p.Stat.ID == id {
			return p.Stat.Slug, nil
		}
	}

	return "", errors.New("slug not found")
}

func (u *User) GetID(slug string) (int, error) {
	problems, err := u.GetProblems()
	if err != nil {
		return 0, err
	}

	for _, p := range problems {
		if p.Stat.Slug == slug {
			return p.Stat.ID, nil
		}
	}

	fmt.Printf("Searching for questions matching '%s'\n", slug)
	slug = url.PathEscape(slug)
	body, err := u.Request("GET", baseUrl+"/problems/api/filter-questions/"+slug, nil)

	var ids []int
	if err = json.Unmarshal(body, &ids); err != nil {
		return 0, err
	}

	if len(ids) > 0 {
		max := 10
		if len(ids) < max {
			max = len(ids)
		}
		p, err := u.SelectQuestion(ids[:max])
		if err == nil {
			return p.Stat.ID, nil
		}
	}

	return 0, err
}
