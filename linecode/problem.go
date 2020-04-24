package linecode

import (
	"github.com/GregLahaye/convert"
	"github.com/GregLahaye/yogurt"
	"github.com/GregLahaye/yogurt/colors"
	"golang.org/x/crypto/ssh/terminal"
	"os"
	"strings"
)

type Problem struct {
	Stat struct {
		ID             int    `json:"question_id"`
		Title          string `json:"question__title"`
		Slug           string `json:"question__title_slug"`
		TotalAccepted  int    `json:"total_acs"`
		TotalSubmitted int    `json:"total_submitted"`
	} `json:"stat"`
	Status     string     `json:"status"`
	Difficulty Difficulty `json:"difficulty"`
	PaidOnly   bool       `json:"paid_only"`
	Starred    bool       `json:"is_favor"`
}

type Difficulty struct {
	Level int `json:"level"`
}

func (d Difficulty) String() string {
	switch d.Level {
	case Easy:
		return "Easy"
	case Medium:
		return "Medium"
	case Hard:
		return "Hard"
	}
	return ""
}

func (d Difficulty) Color() string {
	switch d.Level {
	case Easy:
		return yogurt.Foreground(colors.Lime)
	case Medium:
		return yogurt.Foreground(colors.Orange1)
	case Hard:
		return yogurt.Foreground(colors.Red1)
	}
	return ""
}

func (p Problem) String() string {
	var s strings.Builder

	a := (float64(p.Stat.TotalAccepted) / float64(p.Stat.TotalSubmitted)) * 100

	if p.PaidOnly {
		s.WriteString(yogurt.Foreground(colors.Yellow2))
		s.WriteString("$")
		s.WriteString(yogurt.ResetForeground)
	} else {
		s.WriteString(" ")
	}

	if p.Starred {
		s.WriteString(yogurt.Foreground(colors.Pink1))
		s.WriteString("*")
		s.WriteString(yogurt.ResetForeground)
	} else {
		s.WriteString(" ")
	}

	if p.Status == Accepted {
		s.WriteString(yogurt.Foreground(colors.Lime))
		s.WriteString("#")
		s.WriteString(yogurt.ResetForeground)
	} else {
		s.WriteString(" ")
	}

	// get width of terminal

	w, _, _ := terminal.GetSize(int(os.Stdout.Fd()))
	s.WriteString(" [")
	s.WriteString(convert.PadString(convert.IntToString(p.Stat.ID), 4, true))
	s.WriteString("] ")
	// the length of the string without the slug is 30
	s.WriteString(convert.PadString(p.Stat.Slug, w-30, false))
	s.WriteString("  ")
	s.WriteString(p.Difficulty.Color())
	s.WriteString(convert.PadString(p.Difficulty.String(), 6, false))
	s.WriteString(yogurt.ResetForeground)
	s.WriteString("  (")
	s.WriteString(convert.FloatToString(a))
	s.WriteString(" %)")

	return s.String()
}
