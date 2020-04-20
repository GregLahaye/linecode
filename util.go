package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"golang.org/x/crypto/ssh/terminal"
	"golang.org/x/net/html"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

func PrettyPrint(v interface{}) {
	b, err := json.MarshalIndent(v, "", "  ")

	if err == nil {
		fmt.Println(string(b))
	}
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

func ReadFile(filename string) (string, error) {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func SaveStruct(filename string, v interface{}) error {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}

	if err = ioutil.WriteFile(filename, b, os.ModePerm); err != nil {
		return err
	}

	return nil
}

func LoadStruct(filename string, v interface{}) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}

	b, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}

	if err = json.Unmarshal(b, v); err != nil {
		return err
	}

	return nil
}

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

func MultilineInput() (string, error) {
	s := ""
	reader := bufio.NewReader(os.Stdin)
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

func IntToString(i int) string {
	return strconv.Itoa(i)
}

func FloatToString(f float64) string {
	return strconv.FormatFloat(f, 'f', 2, 64)
}

func Filter(p Problem, f []rune) bool {
	if len(f) == 0 {
		return true
	}

	hit := false      // hit a positive case
	miss := false     // missed a negative case
	positive := false // positive case exists
	negative := false // negative case exists
	for _, c := range f {
		switch c {
		case 'e':
			if p.Difficulty.Level == 1 {
				hit = true
			}
			positive = true
		case 'm':
			if p.Difficulty.Level == 2 {
				hit = true
			}
			positive = true
		case 'h':
			if p.Difficulty.Level == 3 {
				hit = true
			}
			positive = true
		case 'E':
			if p.Difficulty.Level == 1 {
				return false
			}
			miss = true
			negative = true
		case 'M':
			if p.Difficulty.Level == 2 {
				return false
			}
			miss = true
			negative = true
		case 'H':
			if p.Difficulty.Level == 3 {
				return false
			}
			miss = true
			negative = true
		case 'a':
			if p.Status == "ac" {
				hit = true
			}
			positive = true
		case 'A':
			if p.Status == "ac" {
				return false
			}
			miss = true
			negative = true
		case 'l':
			if p.PaidOnly {
				hit = true
			}
			positive = true
		case 'L':
			if p.PaidOnly {
				return false
			}
			miss = true
			negative = true
		case 's':
			if p.Starred {
				hit = true
			}
			positive = true
		case 'S':
			if p.Starred {
				return false
			}
			miss = true
			negative = true
		}
	}

	if positive {
		return hit
	}

	if positive && negative {
		return hit && miss
	}

	return true
}
