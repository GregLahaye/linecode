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
		s += "*"
	} else {
		s += " "
	}

	if p.PaidOnly {
		s += "$"
	} else {
		s += " "
	}

	if p.Status == "ac" {
		s += "#"
	} else {
		s += " "
	}

	s += " [" + PadString(strconv.Itoa(p.Stat.ID), 4, true) + "] "

	s += PadString(p.Stat.Title, 80, false) + " "

	switch p.Difficulty.Level {
	case 1:
		s += "Easy   "
	case 2:
		s += "Medium "
	case 3:
		s += "Hard   "
	}

	f := (float64(p.Stat.TotalAccepted) / float64(p.Stat.TotalSubmitted)) * 100
	s += "(" + strconv.FormatFloat(f, 'f', 2, 64) + "%)"

	fmt.Println(s)
}

func (u *User) ShowQuestion(slug string) error {
	q, err := u.GetQuestion(slug)
	if err != nil {
		return err
	}

	s := ""
	for _, l := range q.CodeSnippets {
		s += l.LangSlug + "   "
	}

	s += "\n\nSample Test Case: " + strconv.Quote(q.SampleTestCase)

	s += "\n\n\n" + ParseHTML(q.Content)

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
