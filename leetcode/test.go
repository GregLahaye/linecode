package leetcode

import (
	"encoding/json"
	"github.com/GregLahaye/linecode/linecode"
)

func TestCode(id int, slug, language, code, testcase string) (linecode.Submission, error) {
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
