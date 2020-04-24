package linecode

import (
	"github.com/GregLahaye/convert"
	"github.com/GregLahaye/yogurt"
	"github.com/GregLahaye/yogurt/colors"
	"strings"
)

type Question struct {
	ID             string        `json:"questionId"`
	Title          string        `json:"title"`
	Slug           string        `json:"titleSlug"`
	Content        string        `json:"content"`
	PaidOnly       bool          `json:"isPaidOnly"`
	Difficulty     string        `json:"difficulty"`
	Tags           []Tag         `json:"topicTags"`
	CodeSnippets   []CodeSnippet `json:"codeSnippets"`
	Status         string        `json:"status"`
	SampleTestCase string        `json:"sampleTestCase"`
}

type CodeSnippet struct {
	Lang string `json:"lang"`
	Slug string `json:"langSlug"`
	Code string `json:"code"`
}

type Tag struct {
	Slug      string `json:"slug"`
	Name      string `json:"name"`
	Questions []int  `json:"questions"`
}

func (q Question) String() string {
	var s strings.Builder

	for _, c := range q.CodeSnippets {
		s.WriteString(" ")
		s.WriteString(yogurt.Foreground(colors.Black))
		s.WriteString(yogurt.Background(colors.DarkOrange))
		s.WriteString(" " + c.Slug + " ")
		s.WriteString(yogurt.ResetForeground)
		s.WriteString(yogurt.ResetBackground)
		s.WriteString(" ")
	}

	s.WriteString("\n\n #")
	s.WriteString(q.ID)
	s.WriteString(" - ")
	s.WriteString(q.Title)
	if q.Status == Accepted {
		s.WriteString(yogurt.Foreground(colors.Lime))
		s.WriteString(" [Accepted]")
		s.WriteString(yogurt.ResetForeground)
	}
	s.WriteString("\n\n")

	s.WriteString(convert.ParseHTML(q.Content))

	return s.String()
}
