package main

import (
	"fmt"
	"github.com/GregLahaye/yogurt"
	"github.com/GregLahaye/yogurt/colors"
	"golang.org/x/net/html"
	"io/ioutil"
	"os"
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
		s += yogurt.Foreground(colors.Yellow1) + "*" + yogurt.ForegroundReset
	} else {
		s += " "
	}

	if p.PaidOnly {
		s += yogurt.Foreground(colors.Yellow1) + "$" + yogurt.ForegroundReset
	} else {
		s += " "
	}

	if p.Status == "ac" {
		s += yogurt.Foreground(colors.Green) + "#" + yogurt.ForegroundReset
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
	s += yogurt.ForegroundReset

	f := (float64(p.Stat.TotalAccepted) / float64(p.Stat.TotalSubmitted)) * 100
	s += "(" + strconv.FormatFloat(f, 'f', 2, 64) + "%)"

	fmt.Println(s)
}

func (u *User) ShowQuestion(id int, lang string) error {
	q, err := u.GetQuestion(id)
	if err != nil {
		return err
	}

	s := ""
	for _, l := range q.CodeSnippets {
		s += yogurt.Background(colors.DarkOrange) + yogurt.Foreground(colors.Black) + " " + l.LangSlug + " " + yogurt.BackgroundReset + " "
	}
	s += yogurt.ForegroundReset

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
	s += q.Difficulty + yogurt.ForegroundReset

	s += "\n ● Sample Test Case: " + strconv.Quote(q.SampleTestCase)

	s += "\n\n\nDescription: \n" + ParseHTML(q.Content)

	fmt.Println(s)

	ext := LanguageToExtension(lang)
	filename := q.QuestionID + "." + q.TitleSlug + "." + ext

	if _, err = os.Stat(filename); err == nil {
		if !Confirm(filename + " already exists. Overwrite (Y/N) ") {
			return nil
		}
	}

	var code string
	for _, l := range q.CodeSnippets {
		if l.LangSlug == lang {
			code = l.Code
		}
	}

	err = ioutil.WriteFile(filename, []byte(code), os.ModePerm)
	if err != nil {
		return err
	}

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
