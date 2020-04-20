package main

import (
	"fmt"
	"github.com/GregLahaye/yogurt"
	"github.com/GregLahaye/yogurt/colors"
	"io/ioutil"
	"os"
	"strconv"
)

func (u *User) ListTags() error {
	s := ""

	tags, err := u.GetTags()
	if err != nil {
		return err
	}

	for _, tag := range tags.Topics {
		s += yogurt.Background(colors.Yellow1) + yogurt.Foreground(colors.Black) + " " + tag.Slug + " " + yogurt.ResetBackground + " "
	}
	s += yogurt.ResetForeground

	fmt.Println(s)

	return nil
}

func (u *User) DisplayStatistics(filters []rune, slugs []string) error {
	problems, err := u.FilteredProblems(filters, slugs)
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

	for _, p := range problems {
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

func (u *User) ListProblems(filters []rune, slugs []string) error {
	problems, err := u.FilteredProblems(filters, slugs)
	if err != nil {
		return err
	}

	for _, p := range problems {
		DisplayProblem(p)
	}

	return nil
}

func (u *User) FilteredProblems(filters []rune, slugs []string) ([]Problem, error) {
	problems, err := u.GetProblems()
	if err != nil {
		return nil, err
	}

	tags, err := u.GetTags()
	if err != nil {
		return nil, err
	}

	s := "\n"
	for i, slug := range slugs {
		if !TagExists(slug, tags) {
			s += "Tag '" + slug + "' does not exist\n"
			slugs = append(slugs[:i], slugs[i+1:]...)
		}
	}

	var filtered []Problem
	for _, p := range problems.Problems {
		if Filter(p, filters) && (len(slugs) == 0 || HasAnyTag(p, slugs, tags)) {
			filtered = append(filtered, p)
		}
	}

	return filtered, err
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

	s += PadString(p.Stat.TitleSlug, 80, false) + " "

	switch p.Difficulty.Level {
	case 1:
		s += yogurt.Foreground(colors.Lime) + "Easy   "
	case 2:
		s += yogurt.Foreground(colors.DarkOrange) + "Medium "
	case 3:
		s += yogurt.Foreground(colors.Red1) + "Hard   "
	}
	s += yogurt.ResetForeground

	f := (float64(p.Stat.TotalAccepted) / float64(p.Stat.TotalSubmitted)) * 100
	s += "(" + strconv.FormatFloat(f, 'f', 2, 64) + "%)"

	fmt.Println(s)
}

func (u *User) DisplayQuestion(id int, save, open bool) error {
	q, err := u.GetQuestion(id)
	if err != nil {
		return err
	}

	s := ""
	for _, l := range q.CodeSnippets {
		s += yogurt.Background(colors.DarkOrange) + yogurt.Foreground(colors.Black) + " " + l.LangSlug + " " + yogurt.ResetBackground + " "
	}
	s += yogurt.ResetForeground

	s += "\n\n ● Tags: "
	for i, t := range q.TopicTags {
		s += t.Slug
		if i < len(q.TopicTags)-1 {
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

	filename := q.QuestionID + "." + q.TitleSlug + "." + u.Language.Extension

	if _, err = os.Stat(filename); (save || open) && os.IsNotExist(err) {
		var code string
		for _, l := range q.CodeSnippets {
			if l.LangSlug == u.Language.Slug {
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
		u.OpenEditor(filename)
	}

	return nil
}

func DisplaySubmission(m Submission) {
	s := ""

	ok := m.Success
	answer := string(m.Answer)
	testcase := m.Input
	if testcase == "" {
		testcase = m.LastTestcase
	}
	passed := m.TotalCorrect
	total := m.TotalTestcases

	var expected string
	var stdout string
	if m.Judge == "large" {
		answer = string(m.Output)
		expected = string(m.ExpectedOutput)
		stdout = string(m.StdOut)
	} else {
		stdout = string(m.Output)
		expected = string(m.ExpectedAnswer)
		if !m.Correct {
			ok = false
		}
	}

	if passed != total {
		ok = false
	}
	if m.Status != "Accepted" {
		ok = false
	}

	if m.Status == "Runtime Error" {
		s += yogurt.Background(colors.DarkOrange) + yogurt.Foreground(colors.Black) + " Runtime Error "
		s += yogurt.ResetBackground + yogurt.ResetForeground
		s += "\n" + m.RuntimeError
	} else if m.Status == "Compile Error" {
		s += yogurt.Background(colors.DarkOrange) + yogurt.Foreground(colors.Black) + " Compile Error "
		s += yogurt.ResetBackground + yogurt.ResetForeground
		s += "\n" + m.CompileError
	} else if ok {
		s += yogurt.Background(colors.Lime) + yogurt.Foreground(colors.Black) + " Accepted "
		s += yogurt.ResetBackground + yogurt.ResetForeground

		s += "\n ● Runtime: " + m.Runtime
		if m.RuntimePercentile > 0 {
			s += ", faster than " + FloatToString(m.RuntimePercentile) + "%"
		}

		s += "\n ● Memory: " + m.Memory
		if m.MemoryPercentile > 0 {
			s += ", less than " + FloatToString(m.MemoryPercentile) + "%"
		}

		Destroy(QuestionFilename(m.QuestionID))
		Destroy(problemsFilename)
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
