package main

import (
	"encoding/json"
	"fmt"
	"github.com/GregLahaye/yogurt"
	"github.com/GregLahaye/yogurt/colors"
	"io/ioutil"
	"os"
	"strconv"
)

func PrettyPrint(v interface{}) {
	b, err := json.MarshalIndent(v, "", "  ")

	if err == nil {
		fmt.Println(string(b))
	}
}

func (u *User) ListTags() error {
	s := ""

	tags, err := u.GetTags()
	if err != nil {
		return err
	}

	for _, tag := range tags {
		s += yogurt.Background(colors.Yellow1) + yogurt.Foreground(colors.Black) + " " + tag.Slug + " " + yogurt.ResetBackground + " "
	}
	s += yogurt.ResetForeground

	fmt.Println(s)

	return nil
}

func ProblemStatus(id int, problems []Problem) int {
	for _, p := range problems {
		if p.Stat.ID == id {
			if p.Status == "ac" {
				return 1
			} else if p.Status == "notac" {
				return 2
			} else {
				return 0
			}
		}
	}

	return -1
}

func HighestID(problems []Problem) int {
	id := 0
	for _, p := range problems {
		if p.Stat.ID > id {
			id = p.Stat.ID
		}
	}

	return id
}

func (u *User) DisplayGraph() error {
	problems, err := u.GetProblems()
	if err != nil {
		return err
	}
	highest := HighestID(problems)

	cols := 50
	rows := highest / cols
	s := ""
	for i := 0; i < rows; i++ {
		s += PadString(IntToString(i*cols), 4, true) + " "
		for j := 0; j < cols; j++ {
			id := i*cols + j
			if id < highest {
				switch ProblemStatus(id, problems) {
				case 0:
					s += yogurt.Foreground(colors.Grey19) + "■ " + yogurt.ResetForeground
				case 1:
					s += yogurt.Foreground(colors.Lime) + "■ " + yogurt.ResetForeground
				case 2:
					s += yogurt.Foreground(colors.Red1) + "● " + yogurt.ResetForeground
				default:
					s += " "
				}
			}
		}
		s += "\n\n"
	}

	fmt.Print(s)

	return nil
}

func (u *User) DisplayStatistics(f Filter) error {
	problems, err := u.GetProblems()
	if err != nil {
		return err
	}

	tags, err := u.GetTags()
	if err != nil {
		return err
	}

	filtered := FilterProblems(problems, tags, f)
	if err != nil {
		return err
	}

	type difficulty struct {
		Difficulty string
		All        int
		Accepted   int
	}

	a := []difficulty{
		{"Easy", 0, 0},
		{"Medium", 0, 0},
		{"Hard", 0, 0},
	}

	for _, p := range filtered {
		l := p.Difficulty.Level - 1
		a[l].All++
		if p.Status == "ac" {
			a[l].Accepted++
		}
	}

	s := ""
	for _, d := range a {
		switch d.Difficulty {
		case "Easy":
			s += yogurt.Foreground(colors.Lime)
		case "Medium":
			s += yogurt.Foreground(colors.DarkOrange)
		case "Hard":
			s += yogurt.Foreground(colors.Red1)
		}

		s += " ● " + PadString(d.Difficulty, 7, false)
		s += PadString(IntToString(d.Accepted)+"/"+IntToString(d.All), 9, true)
		p := float64(d.Accepted) / float64(d.All)
		s += "  " + PadString("("+FloatToString(p*100)+"%)", 9, true)
		s += "  " + ProgressBar(p, 30)
		s += yogurt.ResetForeground + "\n"
	}

	fmt.Print(s)

	return nil
}

func ProgressBar(f float64, size int) string {
	s := ""
	bars := int(f * float64(size))
	for i := 0; i < bars; i++ {
		s += "█"
	}
	for i := 0; i < size-bars; i++ {
		s += "░"
	}
	return s
}

func (u *User) ListProblems(f Filter) error {
	problems, err := u.GetProblems()
	if err != nil {
		return err
	}

	tags, err := u.GetTags()
	if err != nil {
		return err
	}

	filtered := FilterProblems(problems, tags, f)

	for _, p := range filtered {
		DisplayProblem(p)
	}

	return nil
}

func FilterProblems(problems []Problem, tags []Tag, f Filter) []Problem {
	var filtered []Problem
	for _, p := range problems {
		if FilterProblem(p, tags, f) {
			filtered = append(filtered, p)
		}
	}

	return filtered
}

func DisplayProblem(p Problem) {
	s := ""

	if p.Starred {
		s += yogurt.Foreground(colors.Yellow1) + "*" + yogurt.ResetForeground
	} else {
		s += " "
	}

	if p.PaidOnly {
		s += yogurt.Foreground(colors.Yellow1) + "$" + yogurt.ResetForeground
	} else {
		s += " "
	}

	if p.Status == "ac" {
		s += yogurt.Foreground(colors.Green) + "#" + yogurt.ResetForeground
	} else {
		s += " "
	}

	s += "[" + PadString(strconv.Itoa(p.Stat.ID), 4, true) + "] "

	s += PadString(p.Stat.Slug, 80, false) + " "

	switch p.Difficulty.Level {
	case 1:
		s += yogurt.Foreground(colors.Lime) + "Easy   "
	case 2:
		s += yogurt.Foreground(colors.DarkOrange) + "Medium "
	case 3:
		s += yogurt.Foreground(colors.Red1) + "Hard   "
	}
	s += yogurt.ResetForeground

	// TODO: p.Stat.AcceptanceRate
	f := (float64(p.Stat.TotalAccepted) / float64(p.Stat.TotalSubmitted)) * 100
	s += "(" + strconv.FormatFloat(f, 'f', 2, 64) + "%)"

	fmt.Println(s)
}

func (u *User) DisplayQuestion(arg string, save, open bool) error {
	q, err := u.GetQuestion(arg)
	if err != nil {
		return err
	}

	if q.IsPaidOnly {
		fmt.Print(" " + yogurt.Background(colors.Red3))
		fmt.Print("[" + PadString(q.QuestionID, 4, true) + "] " + q.Slug + " is a locked question")
		fmt.Print(yogurt.ResetBackground)
		return nil
	}

	s := ""
	for _, l := range q.CodeSnippets {
		s += yogurt.Background(colors.DarkOrange) + yogurt.Foreground(colors.Black) + " " + l.Slug + " " + yogurt.ResetBackground + " "
	}
	s += yogurt.ResetForeground

	s += "\n\n #" + q.QuestionID + " - " + q.Title
	s += "\n ● Tags: "
	for i, t := range q.Tags {
		s += t.Slug
		if i < len(q.Tags)-1 {
			s += ", "
		}
	}

	s += "\n ● Difficulty: "
	switch q.Difficulty {
	case "Easy":
		s += yogurt.Foreground(colors.Lime)
	case "Medium":
		s += yogurt.Foreground(colors.DarkOrange)
	case "Hard":
		s += yogurt.Foreground(colors.Red1)
	}
	s += q.Difficulty + yogurt.ResetForeground

	s += "\n ● Sample Test Case: " + strconv.Quote(q.SampleTestCase)

	if q.Status == "ac" {
		s += "\n ● " + yogurt.Background(colors.Lime) + yogurt.Foreground(colors.Black) + " Accepted " + yogurt.ResetBackground + yogurt.ResetForeground
	}

	content := ParseHTML(q.Content)
	s += "\n\n\nDescription: \n" + content

	fmt.Println(s)

	filename := q.QuestionID + "." + q.Slug + "." + u.Language.Extension

	if _, err = os.Stat(filename); (save || open) && os.IsNotExist(err) {
		var code string
		for _, l := range q.CodeSnippets {
			if l.Slug == u.Language.Slug {
				code = l.Code
			}
		}

		c := u.Language.Comment.Start + "\n" + content + "\n" + u.Language.Comment.End + "\n\n\n" + code

		err = ioutil.WriteFile(filename, []byte(c), os.ModePerm)
		if err != nil {
			return err
		}
	}

	if open {
		fmt.Print("Press enter to open code in editor...")
		StringInput()
		u.OpenEditor(filename)
	}

	return nil
}

func (u *User) DisplayTest(filename string) error {
	id, slug := SplitFilename(filename)

	tc, err := MultilineInput("Please enter a testcase: (optional)")
	if err != nil {
		return err
	}

	submission, err := u.TestCode(id, slug, filename, tc)
	if err != nil {
		return err
	}

	DisplaySubmission(submission)

	return nil
}

func (u *User) DisplaySubmit(filename string) error {
	id, slug := SplitFilename(filename)

	submission, err := u.SubmitCode(id, slug, filename)
	if err != nil {
		return err
	}

	DisplaySubmission(submission)

	return nil
}

func DisplaySubmission(submission Submission) {
	s := ""

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
		s += yogurt.Background(colors.DarkOrange) + yogurt.Foreground(colors.Black) + " Runtime Error "
		s += yogurt.ResetBackground + yogurt.ResetForeground
		s += "\n" + submission.RuntimeError
	} else if submission.Status == "Compile Error" {
		s += yogurt.Background(colors.DarkOrange) + yogurt.Foreground(colors.Black) + " Compile Error "
		s += yogurt.ResetBackground + yogurt.ResetForeground
		s += "\n" + submission.CompileError
	} else if submission.Status == "Time Limit Exceeded" {
		s += yogurt.Background(colors.Red1) + yogurt.Foreground(colors.Black) + " Time Limit Exceeded "
		s += yogurt.ResetBackground + yogurt.ResetForeground
	} else if ok {
		s += yogurt.Background(colors.Lime) + yogurt.Foreground(colors.Black) + " Accepted "
		s += yogurt.ResetBackground + yogurt.ResetForeground

		s += "\n ● Runtime: " + submission.Runtime
		if submission.RuntimePercentile > 0 {
			s += ", faster than " + FloatToString(submission.RuntimePercentile) + "%"
		}

		s += "\n ● Memory: " + submission.Memory
		if submission.MemoryPercentile > 0 {
			s += ", less than " + FloatToString(submission.MemoryPercentile) + "%"
		}

		CacheDestroy(QuestionFilename(submission.QuestionID))

		CacheDestroy(problemsFilename)
	} else {
		s += yogurt.Background(colors.Red1) + yogurt.Foreground(colors.Black) + " Wrong Answer "
		s += yogurt.ResetBackground + yogurt.ResetForeground

		if total > 0 {
			s += "\nPassed " + IntToString(passed) + " / " + IntToString(total) + " test cases"
		}

		if testcase != "" {
			s += "\nTestcase: \n" + testcase
		}

		s += "\n stdout: " + stdout
		s += "\n Output: " + answer
		s += "\n Expected: " + expected
	}

	fmt.Println(s)
}
