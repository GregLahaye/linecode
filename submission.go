package main

import (
	"encoding/json"
	"fmt"
	"github.com/GregLahaye/yoyo"
	"github.com/GregLahaye/yoyo/styles"
	"strconv"
	"time"
)

type Submission struct {
	QuestionID int    `json:"question_id"`
	Success    bool   `json:"run_success"`
	State      string `json:"state"`
	Status     string `json:"status_msg"`
	Judge      string

	Runtime           string  `json:"status_runtime"`
	RuntimePercentile float64 `json:"runtime_percentile"`
	Memory            string  `json:"status_memory"`
	MemoryPercentile  float64 `json:"memory_percentile"`

	Input        string `json:"input"`
	LastTestcase string `json:"last_testcase"`

	Correct        bool            `json:"correct_answer"`
	Answer         json.RawMessage `json:"code_answer"`
	Output         json.RawMessage `json:"code_output"`
	StdOut         json.RawMessage `json:"std_output"`
	ExpectedOutput json.RawMessage `json:"expected_output"`
	ExpectedAnswer json.RawMessage `json:"expected_code_answer"`

	TotalCorrect   int `json:"total_correct"`
	TotalTestcases int `json:"total_testcases"`

	RuntimeError string `json:"full_runtime_error"`
	CompileError string `json:"full_compile_error"`
}

const retryDelay = 1

func (u *User) TestCode(id int, slug, filename, testcase string) (Submission, error) {
	if testcase == "" {
		q, err := u.GetQuestion(slug)
		if err != nil {
			return Submission{}, nil
		}
		testcase = q.SampleTestCase
		fmt.Println(testcase)
	}

	s := yoyo.Start(styles.Point)
	defer s.End()

	code, err := ReadFile(filename)
	if err != nil {
		return Submission{}, err
	}

	data := dict{"lang": u.Language.Slug, "question_id": id, "test_mode": false, "typed_code": code, "data_input": testcase}
	body, err := u.Request("POST", baseUrl+"/problems/"+slug+"/interpret_solution/", data)
	if err != nil {
		return Submission{}, err
	}

	v := struct {
		InterpretID string `json:"interpret_id"`
	}{}
	if err = json.Unmarshal(body, &v); err != nil {
		return Submission{}, err
	}

	return u.Retry(v.InterpretID)
}

func (u *User) SubmitCode(id int, slug, filename string) (Submission, error) {
	var submission Submission

	s := yoyo.Start(styles.Point)
	defer s.End()

	code, err := ReadFile(filename)
	if err != nil {
		return submission, err
	}

	data := dict{"lang": u.Language.Slug, "question_id": id, "test_mode": false, "typed_code": code}
	body, err := u.Request("POST", baseUrl+"/problems/"+slug+"/submit/", data)
	if err != nil {
		return submission, err
	}

	v := struct {
		SubmissionID int `json:"submission_id"`
	}{}
	if err = json.Unmarshal(body, &v); err != nil {
		return submission, err
	}

	return u.Retry(strconv.Itoa(v.SubmissionID))
}

func (u *User) VerifyResult(id string) (Submission, error) {
	var submission Submission

	body, err := u.Request("GET", baseUrl+"/submissions/detail/"+id+"/check/", nil)
	if err != nil {
		return submission, err
	}

	if err = json.Unmarshal(body, &submission); err != nil {
		return submission, err
	}

	return submission, nil
}

func (u *User) Retry(id string) (Submission, error) {
	submission, err := u.VerifyResult(id)
	if err != nil {
		return submission, err
	}

	for submission.State != "SUCCESS" {
		time.Sleep(time.Second * retryDelay)
		submission, err = u.VerifyResult(id)
		if err != nil {
			return Submission{}, err
		}
	}

	return submission, nil
}
