package main

import (
	"fmt"
	"github.com/GregLahaye/yogurt"
)

type Language struct {
	Name      string
	Slug      string
	Extension string
}

var languages = []Language{
	{Name: "C++", Slug: "cpp", Extension: "cpp"},
	{Name: "Java", Slug: "java", Extension: "java"},
	{Name: "Python", Slug: "python", Extension: "py"},
	{Name: "Python3", Slug: "python3", Extension: "py"},
	{Name: "C", Slug: "c", Extension: "c"},
	{Name: "C#", Slug: "csharp", Extension: "cs"},
	{Name: "JavaScript", Slug: "javascript", Extension: "js"},
	{Name: "Ruby", Slug: "ruby", Extension: "rb"},
	{Name: "Swift", Slug: "swift", Extension: "swift"},
	{Name: "Go", Slug: "golang", Extension: "go"},
	{Name: "Scala", Slug: "scala", Extension: "scala"},
	{Name: "Kotlin", Slug: "kotlin", Extension: "kt"},
	{Name: "Rust", Slug: "rust", Extension: "rs"},
	{Name: "PHP", Slug: "php", Extension: "php"},
}

func SelectLanguage() Language {
	fmt.Println("Please select your default language")

	fmt.Print(yogurt.DisableCursor)
	defer fmt.Print(yogurt.EnableCursor)

	for _, l := range languages {
		fmt.Printf(" [ ] %s\n", l.Name)
	}

	yogurt.CursorUp(len(languages))
	yogurt.CursorForward(2)
	fmt.Print("x")
	yogurt.CursorBackward(1)

	i := 0
	done := false
	for !done {
		c, _ := GetChar()
		switch c {
		case 'j':
			if i < len(languages)-1 {
				fmt.Print(" ")
				yogurt.CursorBackward(1)
				yogurt.CursorDown(1)
				fmt.Print("x")
				yogurt.CursorBackward(1)
				i++
			}
		case 'k':
			if i > 0 {
				fmt.Print(" ")
				yogurt.CursorBackward(1)
				yogurt.CursorUp(1)
				fmt.Print("x")
				yogurt.CursorBackward(1)
				i--
			}
		case 13:
			done = true
		}
	}

	yogurt.CursorUp(i)
	for j := 0; j < len(languages); j++ {
		fmt.Printf(yogurt.ClearLine)
		yogurt.CursorDown(1)
	}

	yogurt.CursorUp(len(languages))
	yogurt.SetColumn(0)

	return languages[i]
}
