package linecode

import (
	"fmt"
	"github.com/GregLahaye/convert"
	"github.com/GregLahaye/yogurt"
	"github.com/GregLahaye/yogurt/colors"
	"strings"
)

const none = "dne"

func DisplayStat(problems []Problem) {
	type v struct {
		Difficulty Difficulty
		All        int
		Accepted   int
	}
	d := []v{
		{Difficulty{Easy}, 0, 0},
		{Difficulty{Medium}, 0, 0},
		{Difficulty{Hard}, 0, 0},
	}

	for _, problem := range problems {
		d[problem.Difficulty.Level-1].All++
		if problem.Status == Accepted {
			d[problem.Difficulty.Level-1].Accepted++
		}
	}

	var s strings.Builder
	for _, i := range d {
		if i.All < 1 {
			continue
		}
		p := (float64(i.Accepted) / float64(i.All)) * 100
		s.WriteString(i.Difficulty.Color())
		s.WriteString(convert.PadString(i.Difficulty.String(), 6, false))
		s.WriteString(yogurt.ResetForeground)
		s.WriteString(" ")
		s.WriteString(convert.PadString(convert.IntToString(i.Accepted), 4, true))
		s.WriteString(" / ")
		s.WriteString(convert.PadString(convert.IntToString(i.All), 4, false))
		s.WriteString(" ")
		s.WriteString(convert.PadString(convert.FloatToString(p), 5, true))
		s.WriteString("%\n")
	}

	fmt.Println(s.String())
}

func DisplayGraph(problems []Problem) {
	highest := problems[0].Stat.ID
	cols := 50
	rows := (highest / cols) + 1

	var s strings.Builder
	for i := 0; i < rows; i++ {
		s.WriteString(convert.PadString(convert.IntToString(i*cols), 4, true))
		for j := 0; j < cols; j++ {
			switch getStatus(i*cols+j, problems) {
			case Accepted:
				s.WriteString(yogurt.Foreground(colors.Lime))
				s.WriteString(" ■")
				s.WriteString(yogurt.ResetForeground)
			case NotAccepted:
				s.WriteString(yogurt.Foreground(colors.Red1))
				s.WriteString(" ■")
				s.WriteString(yogurt.ResetForeground)
			case none:
				s.WriteString("  ")
			default:
				s.WriteString(yogurt.Foreground(colors.Grey19))
				s.WriteString(" ■")
				s.WriteString(yogurt.ResetForeground)
			}
		}
		s.WriteString("\n")
	}

	fmt.Println(s.String())
}

func getStatus(id int, problems []Problem) string {
	for _, problem := range problems {
		if problem.Stat.ID == id {
			return problem.Status
		}
	}

	return none
}
