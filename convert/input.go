package convert

import (
	"bufio"
	"fmt"
	"github.com/GregLahaye/yogurt"
	"golang.org/x/crypto/ssh/terminal"
	"os"
	"strings"
)

func GetChar() (r rune, err error) {
	oldState, err := terminal.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		return
	}
	defer terminal.Restore(int(os.Stdin.Fd()), oldState)

	reader := bufio.NewReader(os.Stdin)

	r, _, err = reader.ReadRune()

	return
}

func StringInput() (string, error) {
	reader := bufio.NewReader(os.Stdin)
	s, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(s), nil
}

func MultilineInput(prompt string) (string, error) {
	s := ""
	reader := bufio.NewReader(os.Stdin)
	fmt.Println(prompt)
	for {
		if i, err := reader.ReadString('\n'); err != nil {
			return s, err
		} else if i == "\n" || i == "\r\n" {
			return s, nil
		} else {
			s += i
		}
	}
}

func Confirm(prompt string) bool {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Printf(prompt)
		s, err := reader.ReadString('\n')
		if err != nil {
			return false
		}
		c := s[0]

		if c == 'y' || c == 'Y' {
			return true
		} else if c == 'n' || c == 'N' {
			return false
		}
	}
}

func Select(options []string) int {
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
