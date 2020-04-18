package main

import (
	"fmt"
	"golang.org/x/net/html"
	"strconv"
	"strings"
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
		s += Foreground(Yellow1) + "*" + ForegroundReset
	} else {
		s += " "
	}

	if p.PaidOnly {
		s += Foreground(Yellow1) + "$" + ForegroundReset
	} else {
		s += " "
	}

	if p.Status == "ac" {
		s += Foreground(Green) + "#" + ForegroundReset
	} else {
		s += " "
	}

	s += "[" + PadString(strconv.Itoa(p.Stat.ID), 4, true) + "] "

	s += PadString(p.Stat.TitleSlug, 80, false) + " "

	switch p.Difficulty.Level {
	case 1:
		s += Foreground(Lime) + "Easy   "
	case 2:
		s += Foreground(DarkOrange) + "Medium "
	case 3:
		s += Foreground(Red1) + "Hard   "
	}
	s += ForegroundReset

	f := (float64(p.Stat.TotalAccepted) / float64(p.Stat.TotalSubmitted)) * 100
	s += "(" + strconv.FormatFloat(f, 'f', 2, 64) + "%)"

	fmt.Println(s)
}

func (u *User) ShowQuestion(id int) error {
	q, err := u.GetQuestion(id)
	if err != nil {
		return err
	}

	s := ""
	for _, l := range q.CodeSnippets {
		s += Background(DarkOrange) + Foreground(Black) + " " + l.LangSlug + " " + BackgroundReset + " "
	}
	s += ForegroundReset

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
		s += Foreground(Lime)
	case "Medium":
		s += Foreground(DarkOrange)
	case "Hard":
		s += Foreground(Red1)
	}
	s += q.Difficulty + ForegroundReset

	s += "\n ● Sample Test Case: " + strconv.Quote(q.SampleTestCase)

	s += "\n\n\nDescription: \n" + ParseHTML(q.Content)

	fmt.Println(s)

	return nil
}

func PadString(str string, max int, left bool) string {
	length := len(str)
	if length > max {
		return str
	}

	difference := max - length
	padding := strings.Repeat(" ", difference)
	if left {
		str = padding + str
	} else {
		str += padding
	}

	return str
}

func ParseHTML(h string) string {
	z := html.NewTokenizer(strings.NewReader(h))

	s := ""
	for {
		tt := z.Next()
		t := z.Token()

		switch tt {
		case html.ErrorToken:
			return s
		case html.TextToken:
			s += t.Data
		}
	}
}
