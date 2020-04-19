package main

import (
	"fmt"
	"github.com/GregLahaye/yogurt"
	"github.com/GregLahaye/yogurt/colors"
	"io/ioutil"
	"os"
	"strconv"
)

func (u *User) ListProblems() error {
	problems, err := u.GetProblems()
	if err != nil {
		return err
	}

	for _, problem := range problems.Problems {
		DisplayProblem(problem)
	}

	return nil
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

func (u *User) DisplayQuestion(id int) error {
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

	s += "\n\n\nDescription: \n" + ParseHTML(q.Content)

	filename := q.QuestionID + "." + q.TitleSlug + "." + u.Language.Extension

	if _, err = os.Stat(filename); os.IsNotExist(err) {
		var code string
		for _, l := range q.CodeSnippets {
			if l.LangSlug == u.Language.Slug {
				code = l.Code
			}
		}

		err = ioutil.WriteFile(filename, []byte(code), os.ModePerm)
		if err != nil {
			return err
		}
	}

	return nil
}

func DisplaySubmission(m Submission) {
	s := ""

	// because different judge types have different json keys, we need to do some checks to determine the actual status
	// of the submission and the actual/expected answer/outputs

	switch m.Status {
	case "Accepted":
		s += yogurt.Background(colors.Lime) + yogurt.Foreground(colors.Black) + " Accepted "
		s += yogurt.ResetBackground + yogurt.ResetForeground

		s += "\n ● Runtime: " + strconv.Itoa(m.Runtime) + " ms"
		if m.RuntimePercentile > 0 {
			s += ", faster than " + FloatToString(m.RuntimePercentile)
		}

		s += "\n ● Memory: " + m.Memory + " ms"
		if m.MemoryPercentile > 0 {
			s += ", less than " + FloatToString(m.MemoryPercentile)
		}
	case "Wrong Answer":
		s += yogurt.Background(colors.Red1) + yogurt.Foreground(colors.Black) + " Wrong Answer "
		s += yogurt.ResetBackground + yogurt.ResetForeground

		s += "\nActual"
		s += "\n Answer: " + string(m.Answer)
		s += "\n Output: " + string(m.Output)

		s += "\nExpected"
		s += "\n Answer: " + string(m.ExpectedAnswer)
		s += "\n Output: " + string(m.ExpectedOutput)
	case "Runtime Error":
		s += yogurt.Background(colors.DarkOrange) + yogurt.Foreground(colors.Black) + " Runtime Error "
		s += yogurt.ResetBackground + yogurt.ResetForeground
		s += "\n" + m.RuntimeError
	case "Compile Error":
		s += yogurt.Background(colors.DarkOrange) + yogurt.Foreground(colors.Black) + " Compile Error "
		s += yogurt.ResetBackground + yogurt.ResetForeground
		s += "\n" + m.CompileError
	default:
		s += "Unknown submission status"
	}

	fmt.Println(s)
}
