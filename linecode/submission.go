package linecode

import (
	"encoding/json"
	"github.com/GregLahaye/convert"
	"github.com/GregLahaye/yogurt"
	"github.com/GregLahaye/yogurt/colors"
	"strings"
)

type Submission struct {
	QuestionID string `json:"question_id"`
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

func (submission Submission) String() string {
	var s strings.Builder

	ok := submission.Success
	answer := string(submission.Answer)
	testcase := submission.Input
	if testcase == "" {
		testcase = submission.LastTestcase
	}
	passed := submission.TotalCorrect
	total := submission.TotalTestcases

	var expected string
	var stdout string
	if submission.Judge == "large" {
		answer = string(submission.Output)
		expected = string(submission.ExpectedOutput)
		stdout = string(submission.StdOut)
	} else {
		stdout = string(submission.Output)
		expected = string(submission.ExpectedAnswer)
		if !submission.Correct {
			ok = false
		}
	}

	if passed != total {
		ok = false
	}
	if submission.Status != "Accepted" {
		ok = false
	}

	if submission.Status == "Runtime Error" {
		s.WriteString(yogurt.Foreground(colors.Orange3))
		s.WriteString(" Runtime Error ")
		s.WriteString(yogurt.ResetForeground)
		s.WriteString("\n ")
		s.WriteString(submission.RuntimeError)
	} else if submission.Status == "Compile Error" {
		s.WriteString(yogurt.Foreground(colors.Orange3))
		s.WriteString(" Compile Error ")
		s.WriteString(yogurt.ResetForeground)
		s.WriteString("\n ")
		s.WriteString(submission.CompileError)
	} else if submission.Status == "Time Limit Exceeded" {
		s.WriteString(" Time Limit Exceeded ")
	} else if ok {
		s.WriteString(yogurt.Foreground(colors.Lime))
		s.WriteString(" Accepted ")
		s.WriteString(yogurt.ResetForeground)

		s.WriteString("\n Runtime: " + submission.Runtime)
		if submission.RuntimePercentile > 0 {
			s.WriteString(", faster than ")
			s.WriteString(convert.FloatToString(submission.RuntimePercentile))
			s.WriteString("%")
		}

		s.WriteString("\n Memory: " + submission.Memory)
		if submission.MemoryPercentile > 0 {
			s.WriteString(", less than ")
			s.WriteString(convert.FloatToString(submission.MemoryPercentile))
			s.WriteString("%")
		}
	} else {
		s.WriteString(yogurt.Foreground(colors.Red1))
		s.WriteString(" Wrong Answer ")
		s.WriteString(yogurt.ResetForeground)

		if total > 0 {
			s.WriteString("\nPassed ")
			s.WriteString(convert.IntToString(passed))
			s.WriteString(" / ")
			s.WriteString(convert.IntToString(total))
			s.WriteString(" test cases")
		}

		if testcase != "" {
			s.WriteString("\nTestcase: \n" + testcase)
		}

		s.WriteString("\n stdout: ")
		s.WriteString(stdout)
		s.WriteString("\n Output: ")
		s.WriteString(answer)
		s.WriteString("\n Expected: ")
		s.WriteString(expected)
	}

	return s.String()
}
