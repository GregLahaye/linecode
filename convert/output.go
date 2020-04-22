package convert

import (
	"fmt"
	"github.com/GregLahaye/yogurt"
	"strings"
)

func Selection(options []string) int {
	l := len(options)

	var s strings.Builder
	s.WriteString(yogurt.DisableCursor)
	for _, o := range options {
		s.WriteString(" [ ] ")
		s.WriteString(o)
		s.WriteString("\n")
	}

	s.WriteString(yogurt.SetColumn(3))
	s.WriteString(yogurt.CursorUp(l))
	s.WriteString("x")
	s.WriteString(yogurt.CursorBackward(1))

	fmt.Print(s.String())

	i := 0
	done := false
	for !done {
		c, _ := GetChar()
		switch c {
		case 3:
			done = true
		case 13:
			done = true
		case 'j':
			if i+1 < l {
				fmt.Print(" ")
				fmt.Print(yogurt.CursorBackward(1))
				fmt.Print(yogurt.CursorDown(1))
				fmt.Print("x")
				fmt.Print(yogurt.CursorBackward(1))
				i++
			}
		case 'k':
			if i > 0 {
				fmt.Print(" ")
				fmt.Print(yogurt.CursorBackward(1))
				fmt.Print(yogurt.CursorUp(1))
				fmt.Print("x")
				fmt.Print(yogurt.CursorBackward(1))
				i--
			}
		}
	}

	s.Reset()

	s.WriteString(yogurt.SetColumn(0))
	s.WriteString(yogurt.CursorUp(i))
	fmt.Print(yogurt.EnableCursor)
	for j := 0; j < l; j++ {
		s.WriteString(yogurt.ClearLine)
		s.WriteString(yogurt.CursorDown(1))
	}
	s.WriteString(yogurt.CursorUp(l))

	fmt.Print(s.String())

	return i
}
