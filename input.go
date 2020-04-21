package main

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/GregLahaye/yogurt"
	"github.com/GregLahaye/yogurt/colors"
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

func (u *User) SelectQuestion(ids []int) (Problem, error) {
	all, err := u.GetProblems()
	if err != nil {
		return Problem{}, err
	}

	fmt.Println(" Question not found, did you mean of these?")

	fmt.Print(yogurt.DisableCursor)
	defer fmt.Print(yogurt.EnableCursor)

	fmt.Println(" [ ]           No, I didn't")

	var problems []Problem
	for _, id := range ids {
		for _, p := range all {
			if p.Stat.ID == id {
				problems = append(problems, p)
			}
		}
	}

	for _, p := range problems {
		fmt.Print(" [ ] ")
		DisplayProblem(p)
	}

	length := len(problems) + 1
	yogurt.CursorUp(length)
	yogurt.CursorForward(2)
	fmt.Print("x")
	yogurt.CursorBackward(1)

	i := 0

	defer func() {
		yogurt.CursorUp(i + 1)
		for j := 0; j < length+1; j++ {
			fmt.Printf(yogurt.ClearLine)
			yogurt.CursorDown(1)
		}

		yogurt.CursorUp(length + 1)
		yogurt.SetColumn(0)
	}()

	done := false
	for !done {
		c, _ := GetChar()
		switch c {
		case 'j':
			if i < length-1 {
				yogurt.SetColumn(0)
				fmt.Print(" [ ] ")
				if i > 0 {
					DisplayProblem(problems[i-1])
					yogurt.CursorUp(1)
				}
				i++

				yogurt.CursorDown(1)
				yogurt.SetColumn(0)
				fmt.Print(" [x] ")
				if i <= len(problems) {
					fmt.Print(yogurt.Foreground(colors.Yellow1))
					DisplayProblem(problems[i-1])
					fmt.Print(yogurt.ResetForeground)
					yogurt.CursorUp(1)
				}
			}
		case 'k':
			if i > 0 {
				yogurt.SetColumn(0)
				fmt.Print(" [ ] ")
				DisplayProblem(problems[i-1])
				yogurt.CursorUp(1)
				i--

				yogurt.CursorUp(1)
				yogurt.SetColumn(0)
				fmt.Print(" [x] ")
				if i > 0 {
					fmt.Print(yogurt.Foreground(colors.Yellow1))
					DisplayProblem(problems[i-1])
					fmt.Print(yogurt.ResetForeground)
					yogurt.CursorUp(1)
				}
			}
		case 13:
			done = true
		case 3, 113:
			return Problem{}, errors.New("selection cancelled")
		}
	}

	index := i - 1
	if index >= 0 && index < len(problems) {
		return problems[index], nil
	} else {
		return Problem{}, errors.New("selection cancelled")
	}
}
