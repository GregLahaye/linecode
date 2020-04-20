package main

import (
	"fmt"
	"github.com/GregLahaye/yogurt"
)

type Language struct {
	Name      string
	Slug      string
	Extension string
	Comment   Comment
}

type Comment struct {
	Start string
	End   string
}

var languages = []Language{
	{
		Name:      "C++",
		Slug:      "cpp",
		Extension: "cpp",
		Comment: Comment{
			Start: `/*`,
			End:   `*/`,
		},
	},
	{
		Name:      "Java",
		Slug:      "java",
		Extension: "java",
		Comment: Comment{
			Start: `/*`,
			End:   `*/`,
		},
	},
	{
		Name:      "Python",
		Slug:      "python",
		Extension: "py",
		Comment: Comment{
			Start: `"""`,
			End:   `"""`,
		},
	},
	{
		Name:      "Python3",
		Slug:      "python3",
		Extension: "py",
		Comment: Comment{
			Start: `"""`,
			End:   `"""`,
		},
	},

	{
		Name:      "C",
		Slug:      "c",
		Extension: "c",
		Comment: Comment{
			Start: `/*`,
			End:   `*/`,
		},
	},
	{
		Name:      "C#",
		Slug:      "csharp",
		Extension: "cs",
		Comment: Comment{
			Start: `/*`,
			End:   `*/`,
		},
	},
	{
		Name:      "JavaScript",
		Slug:      "javascript",
		Extension: "js",
		Comment: Comment{
			Start: `/*`,
			End:   `*/`,
		},
	},
	{
		Name:      "Ruby",
		Slug:      "ruby",
		Extension: "rb",
		Comment: Comment{
			Start: `=begin`,
			End: `=end	`,
		},
	},
	{
		Name:      "Swift",
		Slug:      "swift",
		Extension: "swift",
		Comment: Comment{
			Start: `/*`,
			End:   `*/`,
		},
	},
	{
		Name:      "Go",
		Slug:      "golang",
		Extension: "go",
		Comment: Comment{
			Start: `/*`,
			End:   `*/`,
		},
	},
	{
		Name:      "Scala",
		Slug:      "scala",
		Extension: "scala",
		Comment: Comment{
			Start: `/*`,
			End:   `*/`,
		},
	},
	{
		Name:      "Kotlin",
		Slug:      "kotlin",
		Extension: "kt",
		Comment: Comment{
			Start: `/*`,
			End:   `*/`,
		},
	},
	{
		Name:      "Rust",
		Slug:      "rust",
		Extension: "rs",
		Comment: Comment{
			Start: `/*`,
			End:   `*/`,
		},
	},
	{
		Name:      "PHP",
		Slug:      "php",
		Extension: "php",
		Comment: Comment{
			Start: `/*`,
			End:   `*/`,
		},
	},
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
