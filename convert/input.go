package convert

import (
	"bufio"
	"fmt"
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