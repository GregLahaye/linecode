package leetcode

import (
	"encoding/json"
	"fmt"
	"github.com/GregLahaye/linecode/config"
	"github.com/GregLahaye/linecode/convert"
	"github.com/GregLahaye/linecode/linecode"
	"github.com/GregLahaye/linecode/store"
	"strings"
)

func TestCode(filename string) (linecode.Submission, error) {
	var submission linecode.Submission

	id, slug, err := parseFilename(filename)
	if err != nil {
		return submission, err
	}

	c, _ := config.Config()
	language := c.Language

	code, err := store.ReadFile(filename)
	if err != nil {
		return submission, err
	}

	testcase, err := convert.MultilineInput("Testcase (optional): ")
	if err != nil {
		return submission, err
	}

	if strings.TrimSpace(testcase) == "" {
		question, err := GetQuestion(id)
		if err != nil {
			return submission, err
		}
		testcase = question.SampleTestCase
		fmt.Println(testcase)
	}

	clearQuestion(id, slug)

	return testCode(id, slug, language, code, testcase)
}

func testCode(id, slug, language, code, testcase string) (linecode.Submission, error) {
	var submission linecode.Submission

	data := dict{"lang": language, "question_id": id, "typed_code": code, "data_input": testcase}
	body, err := request("POST", "/problems/"+slug+"/interpret_solution/", data)
	if err != nil {
		return submission, err
	}

	v := struct {
		InterpretID string `json:"interpret_id"`
	}{}
	if err = json.Unmarshal(body, &v); err != nil {
		return submission, err
	}

	return retry(v.InterpretID)
}
