package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/GregLahaye/yoyo"
	"github.com/GregLahaye/yoyo/styles"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

type Problems struct {
	Username       string    `json:"user_name"`
	Solved         int       `json:"num_solved"`
	Total          int       `json:"num_total"`
	AcceptedEasy   int       `json:"ac_easy"`
	AcceptedMedium int       `json:"ac_medium"`
	AcceptedHard   int       `json:"ac_hard"`
	Problems       []Problem `json:"stat_status_pairs"`
	FrequencyHigh  int       `json:"frequency_high"`
	FrequencyMid   int       `json:"frequency_mid"`
	Category       string    `json:"category_slug"`
}

type Problem struct {
	Stat struct {
		ID             int    `json:"question_id"`
		Live           bool   `json:"question__article__live"`
		ArticleSlug    string `json:"question__article__slug"`
		Title          string `json:"question__title"`
		TitleSlug      string `json:"question__title_slug"`
		Hidden         bool   `json:"question__hide"`
		TotalAccepted  int    `json:"total_acs"`
		TotalSubmitted int    `json:"total_submitted"`
		DisplayID      int    `json:"frontend_question_id"`
		IsNew          bool   `json:"is_new_question"`
	} `json:"stat"`
	Status     string `json:"status"`
	Difficulty struct {
		Level int `json:"level"`
	} `json:"difficulty"`
	PaidOnly  bool `json:"paid_only"`
	Starred   bool `json:"is_favor"`
	Frequency int  `json:"frequency"`
	Progress  int  `json:"progress"`
}

type Data struct {
	Data struct {
		Question RawQuestion `json:"question"`
	} `json:"data"`
}

type RawQuestion struct {
	QuestionID string `json:"questionId"`
	Title      string `json:"title"`
	TitleSlug  string `json:"titleSlug"`
	Content    string `json:"content"`
	IsPaidOnly bool   `json:"isPaidOnly"`
	Difficulty string `json:"difficulty"`
	TopicTags  []struct {
		Name string `json:"name"`
		Slug string `json:"slug"`
	} `json:"topicTags"`
	CodeSnippets   []CodeSnippet   `json:"codeSnippets"`
	Stats          json.RawMessage `json:"stats"`
	Status         string          `json:"status"`
	SampleTestCase string          `json:"sampleTestCase"`
}

type Question struct {
	QuestionID string `json:"questionId"`
	Title      string `json:"title"`
	TitleSlug  string `json:"titleSlug"`
	Content    string `json:"content"`
	IsPaidOnly bool   `json:"isPaidOnly"`
	Difficulty string `json:"difficulty"`
	TopicTags  []struct {
		Name string `json:"name"`
		Slug string `json:"slug"`
	} `json:"topicTags"`
	CodeSnippets   []CodeSnippet `json:"codeSnippets"`
	Stats          Stats         `json:"stats"`
	Status         string        `json:"status"`
	SampleTestCase string        `json:"sampleTestCase"`
}

type CodeSnippet struct {
	Lang     string `json:"lang"`
	LangSlug string `json:"langSlug"`
	Code     string `json:"code"`
}

type Stats struct {
	TotalAccepted      string `json:"totalAccepted"`
	TotalSubmission    string `json:"totalSubmission"`
	TotalAcceptedRaw   int    `json:"totalAcceptedRaw"`
	TotalSubmissionRaw int    `json:"totalSubmissionRaw"`
	AcceptanceRate     string `json:"acRate"`
}

type RunResult struct {
	InterpretID string `json:"interpret_id"`
	TestCase    string `json:"test_case"`
}

type SubmissionResult struct {
	SubmissionID int `json:"submission_id"`
}

type Submission struct {
	Success bool   `json:"run_success"`
	State   string `json:"state"`
	Status  string `json:"status_msg"`
	Judge   string

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

type Favorites struct {
	Favorites struct {
		PrivateFavorites []struct {
			Hash string `json:"id_hash"`
			Name string `json:"name"`
		} `json:"private_favorites"`
	} `json:"favorites"`
}

const problemsFilename = "problems.json"

func (u *User) Request(method, url string, body dict) (*http.Response, error) {
	client := &http.Client{}

	b, err := json.Marshal(body)
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
	req.Header.Set("Referer", "https://leetcode.com/")
	req.Header.Set("Content-Type", "application/json")

	return client.Do(req)
}

func (u *User) Retry(id string) (Submission, error) {
	submission, err := u.VerifyResult(id)
	if err != nil {
		return Submission{}, err
	}
	for submission.State != "SUCCESS" {
		time.Sleep(time.Second * 1)
		submission, err = u.VerifyResult(id)
		if err != nil {
			return Submission{}, err
		}
	}

	return submission, nil
}

func (u *User) GetSlug(id int) (string, error) {
	problems, err := u.GetProblems()
	if err != nil {
		return "", err
	}

	for _, p := range problems.Problems {
		if p.Stat.ID == id {
			return p.Stat.TitleSlug, nil
		}
	}

	return "", errors.New("slug not found")
}

func (u *User) GetProblems() (Problems, error) {
	var problems Problems
	if err := LoadStruct(problemsFilename, &problems); err != nil {
		problems, err = u.DownloadProblems()
		if err != nil {
			return Problems{}, err
		}

		if err = SaveStruct(problemsFilename, problems); err != nil {
			return Problems{}, err
		}
	}

	return problems, nil
}

func (u *User) DownloadProblems() (Problems, error) {
	s := yoyo.Start(styles.Simple)
	defer s.End()

	resp, err := u.Request("GET", "https://leetcode.com/api/problems/algorithms/", nil)
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

func (u *User) GetQuestion(id int) (Question, error) {
	slug, err := u.GetSlug(id)
	if err != nil {
		return Question{}, err
	}

	s := yoyo.Start(styles.Simple)
	defer s.End()

	data := dict{"variables": dict{"titleSlug": slug}, "operationName": "questionData", "query": "query questionData($titleSlug: String!) {\n  question(titleSlug: $titleSlug) {\n    questionId\n    title\n    titleSlug\n    content\n    isPaidOnly\n    difficulty\n    isLiked\n    topicTags {\n      name\n      slug\n    }\n    codeSnippets {\n      lang\n      langSlug\n      code\n    }\n    stats\n    status\n    sampleTestCase\n    }\n}"}
	resp, err := u.Request("POST", "https://leetcode.com/graphql", data)
	if err != nil {
		return Question{}, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return Question{}, err
	}

	d := Data{}
	if err = json.Unmarshal(body, &d); err != nil {
		return Question{}, err
	}

	q, err := parse(d)
	if err != nil {
		return Question{}, err
	}

	return q, nil
}

func (u *User) TestCode(id int, filename, testcase string) (Submission, error) {
	q, err := u.GetQuestion(id)
	slug := q.TitleSlug
	if testcase == "" {
		testcase = q.SampleTestCase
		fmt.Println(testcase)
	}

	s := yoyo.Start(styles.Simple)
	defer s.End()

	code, err := ReadFile(filename)
	if err != nil {
		return Submission{}, err
	}

	data := dict{"lang": u.Language.Slug, "question_id": id, "test_mode": false, "typed_code": code, "data_input": testcase}
	resp, err := u.Request("POST", "https://leetcode.com/problems/"+slug+"/interpret_solution/", data)
	if err != nil {
		return Submission{}, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return Submission{}, err
	}

	result := RunResult{}
	if err = json.Unmarshal(body, &result); err != nil {
		return Submission{}, err
	}

	return u.Retry(result.InterpretID)
}

func (u *User) SubmitCode(id int, filename string) (Submission, error) {
	slug, err := u.GetSlug(id)
	if err != nil {
		return Submission{}, err
	}

	s := yoyo.Start(styles.Simple)
	defer s.End()

	code, err := ReadFile(filename)
	if err != nil {
		return Submission{}, err
	}

	data := dict{"lang": u.Language.Slug, "question_id": id, "test_mode": false, "typed_code": code}
	resp, err := u.Request("POST", "https://leetcode.com/problems/"+slug+"/submit/", data)
	if err != nil {
		return Submission{}, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return Submission{}, err
	}

	result := SubmissionResult{}
	if err = json.Unmarshal(body, &result); err != nil {
		return Submission{}, err
	}

	return u.Retry(strconv.Itoa(result.SubmissionID))
}

func (u *User) VerifyResult(id string) (Submission, error) {
	resp, err := u.Request("GET", "https://leetcode.com/submissions/detail/"+id+"/check/", nil)
	if err != nil {
		return Submission{}, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return Submission{}, err
	}

	result := Submission{}
	if err = json.Unmarshal(body, &result); err != nil {
		return Submission{}, err
	}

	return result, nil
}

func parse(raw Data) (Question, error) {
	q := Question{}

	if v, err := strconv.Unquote(string(raw.Data.Question.Stats)); err != nil {
		return q, err
	} else {
		if err = json.Unmarshal([]byte(v), &q.Stats); err != nil {
			return q, err
		}
	}

	q.QuestionID = raw.Data.Question.QuestionID
	q.Title = raw.Data.Question.Title
	q.TitleSlug = raw.Data.Question.TitleSlug
	q.Content = raw.Data.Question.Content
	q.IsPaidOnly = raw.Data.Question.IsPaidOnly
	q.Difficulty = raw.Data.Question.Difficulty
	q.TopicTags = raw.Data.Question.TopicTags
	q.CodeSnippets = raw.Data.Question.CodeSnippets
	q.Status = raw.Data.Question.Status
	q.SampleTestCase = raw.Data.Question.SampleTestCase

	return q, nil
}

func (u *User) FindFavorites() error {
	resp, err := u.Request("GET", "https://leetcode.com/list/api/questions", nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	f := Favorites{}
	err = json.Unmarshal(body, &f)
	if err != nil {
		return err
	}

	for _, i := range f.Favorites.PrivateFavorites {
		if i.Name == "Favorite" {
			u.Hash = i.Hash
		}
	}

	return nil
}

func (u *User) Star(id int) error {
	data := dict{"favorite_id_hash": u.Hash, "question_id": id}
	resp, err := u.Request("POST", "https://leetcode.com/list/api/questions", data)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

func (u *User) UnStar(id int) error {
	l := "https://leetcode.com/list/api/questions/" + u.Hash + "/" + IntToString(id)
	resp, err := u.Request("DELETE", l, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}
