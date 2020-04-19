package main

import (
	"bufio"
	"fmt"
	"os"
)

var languages = [][]string{
	{"cpp", "cpp"},
	{"java", "java"},
	{"python", "py"},
	{"python3", "py"},
	{"c", "c"},
	{"csharp", "cs"},
	{"javascript", "js"},
	{"ruby", "rb"},
	{"swift", "swift"},
	{"golang", "go"},
	{"scala", "scala"},
	{"kotlin", "kt"},
	{"rust", "rs"},
	{"php", "php"},
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

func LanguageToExtension(lang string) string {
	for _, l := range languages {
		if l[0] == lang {
			return l[1]
		}
	}

	return ""
}

func ExtensionToLanguage(ext string) string {
	for _, l := range languages {
		if l[1] == ext {
			return l[0]
		}
	}

	return ""
}
