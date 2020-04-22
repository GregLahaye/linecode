package leetcode

import (
	"encoding/json"
	"fmt"
	"github.com/GregLahaye/linecode/linecode"
	"github.com/GregLahaye/linecode/store"
	"path"
	"strconv"
	"strings"
	"time"
)

const wait = 1

func retry(id string) (linecode.Submission, error) {
	var submission linecode.Submission

	submission, err := verify(id)
	if err != nil {
		return submission, err
	}

	for submission.State != "SUCCESS" {
		time.Sleep(time.Second * wait)
		submission, err = verify(id)
		if err != nil {
			return submission, err
		}
	}

	return submission, nil
}

func verify(id string) (linecode.Submission, error) {
	var submission linecode.Submission

	body, err := request("GET", "/submissions/detail/"+id+"/check/", nil)
	if err != nil {
		return submission, err
	}

	if err = json.Unmarshal(body, &submission); err != nil {
		fmt.Println(string(body))
		return submission, err
	}

	return submission, nil
}

func parseFilename(filename string) (string, string, error) {
	parts := strings.Split(filename, ".")
	if _, err := strconv.Atoi(parts[0]); err != nil {
		return "", "", err
	}
	return parts[0], parts[1], nil
}

func clearQuestion(id, slug string) {
	filename := fmt.Sprintf("%s.%s.json", id, slug)
	p := path.Join(questionDirectory, filename)
	_ = store.RemoveFromCache(p)
	_ = store.RemoveFromCache(problemsFilename)
}
