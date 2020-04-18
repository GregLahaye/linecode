package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"
)

type User struct {
	LeetCodeSession string
	CSRFToken       string
}

type Problems struct {
	UserName        string    `json:"user_name"`
	NumSolved       int       `json:"num_solved"`
	NumTotal        int       `json:"num_total"`
	AcEasy          int       `json:"ac_easy"`
	AcMedium        int       `json:"ac_medium"`
	AcHard          int       `json:"ac_hard"`
	StatStatusPairs []Problem `json:"stat_status_pairs"`
	FrequencyHigh   int       `json:"frequency_high"`
	FrequencyMid    int       `json:"frequency_mid"`
	CategorySlug    string    `json:"category_slug"`
}

type Problem struct {
	Stat struct {
		QuestionID          int    `json:"question_id"`
		QuestionArticleLive bool   `json:"question__article__live"`
		QuestionArticleSlug string `json:"question__article__slug"`
		QuestionTitle       string `json:"question__title"`
		QuestionTitleSlug   string `json:"question__title_slug"`
		QuestionHide        bool   `json:"question__hide"`
		TotalAcs            int    `json:"total_acs"`
		TotalSubmitted      int    `json:"total_submitted"`
		FrontendQuestionID  int    `json:"frontend_question_id"`
		IsNewQuestion       bool   `json:"is_new_question"`
	} `json:"stat"`
	Difficulty struct {
		Level int `json:"level"`
	} `json:"difficulty"`
	PaidOnly  bool `json:"paid_only"`
	IsFavor   bool `json:"is_favor"`
	Frequency int  `json:"frequency"`
	Progress  int  `json:"progress"`
}

type Question struct {
	Content        string         `json:"content"`
	Stats          Stats          `json:"stats"`
	CodeDefinition CodeDefinition `json:"codeDefinition"`
	SampleTestCase string         `json:"sampleTestCase"`
	EnableRunCode  bool           `json:"enableRunCode"`
	MetaData       MetaData       `json:"metaData"`
}

type RawQuestion struct {
	Data struct {
		Question map[string]json.RawMessage `json:"question"`
	} `json:"data"`
}

type Stats struct {
	TotalAccepted      string `json:"totalAccepted"`
	TotalSubmission    string `json:"totalSubmission"`
	TotalAcceptedRaw   int    `json:"totalAcceptedRaw"`
	TotalSubmissionRaw int    `json:"totalSubmissionRaw"`
	AcRate             string `json:"acRate"`
}

type CodeDefinition []struct {
	Value       string `json:"value"`
	Text        string `json:"text"`
	DefaultCode string `json:"defaultCode"`
}
type MetaData struct {
	Name   string `json:"name"`
	Params []struct {
		Name string `json:"name"`
		Type string `json:"type"`
	} `json:"params"`
	Return struct {
		Type string `json:"type"`
		Size int    `json:"size"`
	} `json:"return"`
}

type RunResult struct {
	InterpretID string `json:"interpret_id"`
	TestCase    string `json:"test_case"`
}

type SubmissionResult struct {
	SubmissionID int `json:"submission_id"`
}

type Submission struct {
	StatusCode        int      `json:"status_code"`
	Lang              string   `json:"lang"`
	RunSuccess        bool     `json:"run_success"`
	RuntimeError      string   `json:"runtime_error"`
	FullRuntimeError  string   `json:"full_runtime_error"`
	StatusRuntime     string   `json:"status_runtime"`
	Memory            int      `json:"memory"`
	CodeAnswer        []string `json:"code_answer"`
	CodeOutput        []string `json:"code_output"`
	ElapsedTime       int      `json:"elapsed_time"`
	TaskFinishTime    int64    `json:"task_finish_time"`
	TotalCorrect      int      `json:"total_correct"`
	TotalTestcases    int      `json:"total_testcases"`
	RuntimePercentile float64  `json:"runtime_percentile"`
	StatusMemory      string   `json:"status_memory"`
	MemoryPercentile  float64  `json:"memory_percentile"`
	PrettyLang        string   `json:"pretty_lang"`
	SubmissionID      string   `json:"submission_id"`
	StatusMsg         string   `json:"status_msg"`
	State             string   `json:"state"`
}

func main() {
	u := User{}
	if err := Login(&u); err != nil {
		log.Fatal(err)
	} else {
		// we need to get a new CSRF token after request

		//problems, err := u.GetProblems()
		//if err != nil {
		//	log.Fatal(err)
		//}
		//PrettyPrint(problems)

		//q, err := u.GetQuestion("two-sum")
		//if err != nil {
		//	log.Fatal(err)
		//}

		code, err := ioutil.ReadFile("two-sum.py")
		if err != nil {
			log.Fatal(err)
		}
		//data := []byte(code)
		//err = ioutil.WriteFile("two-sum.py", data, os.ModePerm)
		//if err != nil {
		//	log.Fatal(err)
		//}

		v, _ := u.TestCode(1, "two-sum", "python3", string(code))
		r, _ := u.VerifyResult(string(v.InterpretID))
		for r.State != "SUCCESS" {
			time.Sleep(time.Second * 1)
			r, _ = u.VerifyResult(string(v.InterpretID))
		}
		PrettyPrint(r)

		m, _ := u.SubmitCode(1, "two-sum", "python3", string(code))
		r, _ = u.VerifyResult(strconv.Itoa(m.SubmissionID))
		for r.State != "SUCCESS" {
			time.Sleep(time.Second * 1)
			r, _ = u.VerifyResult(strconv.Itoa(m.SubmissionID))
		}
		PrettyPrint(r)
	}
}

func (u *User) SubmitCode(id int, slug, lang, code string) (SubmissionResult, error) {
	client := &http.Client{}

	data := dict{"data_input": "[2, 7, 11, 15]\n9", "lang": lang, "question_id": id, "test_mode": false, "typed_code": code}
	b, _ := json.Marshal(data)

	req, _ := http.NewRequest("POST", "https://leetcode.com/problems/"+slug+"/submit/", bytes.NewReader(b))

	req.AddCookie(&http.Cookie{Name: "csrftoken", Value: u.CSRFToken, Domain: ".leetcode.com"})
	req.AddCookie(&http.Cookie{Name: "LEETCODE_SESSION", Value: u.LeetCodeSession, Domain: ".leetcode.com"})

	req.Header.Set("X-CSRFToken", u.CSRFToken)
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	req.Header.Set("Referer", "https://leetcode.com/")
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return SubmissionResult{}, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	v := SubmissionResult{}
	json.Unmarshal(body, &v)

	return v, nil
}

func (u *User) TestCode(id int, slug, lang, code string) (RunResult, error) {
	client := &http.Client{}

	data := dict{"data_input": "[2, 7, 11, 15]\n9", "lang": lang, "question_id": id, "test_mode": false, "typed_code": code}
	b, _ := json.Marshal(data)

	req, _ := http.NewRequest("POST", "https://leetcode.com/problems/"+slug+"/interpret_solution/", bytes.NewReader(b))

	req.AddCookie(&http.Cookie{Name: "csrftoken", Value: u.CSRFToken, Domain: ".leetcode.com"})
	req.AddCookie(&http.Cookie{Name: "LEETCODE_SESSION", Value: u.LeetCodeSession, Domain: ".leetcode.com"})

	req.Header.Set("X-CSRFToken", u.CSRFToken)
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	req.Header.Set("Referer", "https://leetcode.com/")
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return RunResult{}, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	v := RunResult{}
	json.Unmarshal(body, &v)

	return v, nil
}

func (u *User) VerifyResult(id string) (Submission, error) {
	client := &http.Client{}

	req, _ := http.NewRequest("GET", "https://leetcode.com/submissions/detail/"+id+"/check/", nil)

	req.AddCookie(&http.Cookie{Name: "csrftoken", Value: u.CSRFToken, Domain: ".leetcode.com"})
	req.AddCookie(&http.Cookie{Name: "LEETCODE_SESSION", Value: u.LeetCodeSession, Domain: ".leetcode.com"})

	req.Header.Set("X-CSRFToken", u.CSRFToken)
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	req.Header.Set("Referer", "https://leetcode.com/")
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return Submission{}, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	result := Submission{}
	json.Unmarshal(body, &result)

	return result, nil
}

func (u *User) GetQuestion(slug string) (Question, error) {
	client := &http.Client{}

	query := dict{"variables": dict{"titleSlug": slug}, "operationName": "getQuestionDetail", "query": "query getQuestionDetail($titleSlug: String!) { question(titleSlug: $titleSlug) { content stats codeDefinition sampleTestCase enableRunCode metaData translatedContent } }"}
	b, _ := json.Marshal(query)

	req, _ := http.NewRequest("POST", "https://leetcode.com/graphql", bytes.NewReader(b))

	req.AddCookie(&http.Cookie{Name: "csrftoken", Value: u.CSRFToken, Domain: ".leetcode.com"})
	req.AddCookie(&http.Cookie{Name: "LEETCODE_SESSION", Value: u.LeetCodeSession, Domain: ".leetcode.com"})

	req.Header.Set("X-CSRFToken", u.CSRFToken)
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	req.Header.Set("Referer", "https://leetcode.com/problems/"+slug+"/description/")
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return Question{}, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return Question{}, err
	}

	raw := RawQuestion{}
	if err = json.Unmarshal(body, &raw); err != nil {
		return Question{}, err
	}

	q, err := parseQuestion(raw)
	if err != nil {
		return Question{}, err
	}

	return q, nil
}

func parseQuestion(raw RawQuestion) (Question, error) {
	q := Question{}

	if err := json.Unmarshal(raw.Data.Question["content"], &q.Content); err != nil {
		return q, err
	}

	if v, err := strconv.Unquote(string(raw.Data.Question["stats"])); err != nil {
		return q, err
	} else {
		if err = json.Unmarshal([]byte(v), &q.Stats); err != nil {
			return q, err
		}
	}

	if v, err := strconv.Unquote(string(raw.Data.Question["codeDefinition"])); err != nil {
		return q, err
	} else {
		if err = json.Unmarshal([]byte(v), &q.CodeDefinition); err != nil {
			return q, err
		}
	}

	if err := json.Unmarshal(raw.Data.Question["sampleTestCase"], &q.SampleTestCase); err != nil {
		return q, err
	}

	if err := json.Unmarshal(raw.Data.Question["enableRunCode"], &q.EnableRunCode); err != nil {
		return q, err
	}

	if v, err := strconv.Unquote(string(raw.Data.Question["metaData"])); err != nil {
		return q, err
	} else {
		if err = json.Unmarshal([]byte(v), &q.MetaData); err != nil {
			return q, err
		}
	}

	return q, nil
}

func (u *User) GetProblems() (Problems, error) {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", "https://leetcode.com/api/problems/algorithms/", nil)
	req.AddCookie(&http.Cookie{Name: "csrftoken", Value: u.CSRFToken, Domain: ".leetcode.com"})
	req.AddCookie(&http.Cookie{Name: "LEETCODE_SESSION", Value: u.LeetCodeSession, Domain: ".leetcode.com"})

	resp, err := client.Do(req)
	if err != nil {
		return Problems{}, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return Problems{}, err
	}

	problems := Problems{}
	if err = json.Unmarshal(body, &problems); err != nil {
		return Problems{}, err
	}

	return problems, nil
}

func PrettyPrint(v interface{}) error {
	b, err := json.MarshalIndent(v, "", "  ")

	if err == nil {
		fmt.Println(string(b))
	}

	return err
}
