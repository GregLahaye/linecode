package leetcode

import (
	"encoding/json"
	"github.com/GregLahaye/linecode/linecode"
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
		return submission, err
	}

	return submission, nil
}
