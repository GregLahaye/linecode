package cmd

import (
	"fmt"
	"github.com/GregLahaye/browser"
	"github.com/GregLahaye/linecode/leetcode"
)

var questionCmd = &Command{
	Name:    "question",
	Aliases: []string{"q"},
	Run: func(cmd *Command, args []string) error {
		arg := args[0]
		question, err := leetcode.GetQuestion(arg)
		if err != nil {
			return err
		}
		if question.PaidOnly {
			return fmt.Errorf("%s is a locked question", question.Slug)
		}
		leetcode.SaveSnippet(question)
		fmt.Println(question)
		return nil
	},
	ArgN: 1,
}

var testCmd = &Command{
	Name:    "test",
	Aliases: []string{"t"},
	Run: func(cmd *Command, args []string) error {
		filename := leetcode.FindFile(args[0])
		submission, err := leetcode.TestCode(filename)
		if err != nil {
			return err
		}
		fmt.Println(submission)
		return nil
	},
	ArgN: 1,
}

var submitCmd = &Command{
	Name:    "submit",
	Aliases: []string{"s"},
	Run: func(cmd *Command, args []string) error {
		filename := leetcode.FindFile(args[0])
		submission, err := leetcode.SubmitCode(filename)
		if err != nil {
			return err
		}
		submission.Judge = "large"
		fmt.Println(submission)
		return nil
	},
	ArgN: 1,
}

var starCmd = &Command{
	Name: "star",
	Run: func(cmd *Command, args []string) error {
		arg := args[0]
		return leetcode.Star(arg)
	},
}

var unstarCmd = &Command{
	Name: "unstar",
	Run: func(cmd *Command, args []string) error {
		arg := args[0]
		return leetcode.Unstar(arg)
	},
}

var openCmd = &Command{
	Name:    "open",
	Aliases: []string{"o"},
	Run: func(cmd *Command, args []string) error {
		_, slug, err := leetcode.Search(args[0])
		if err != nil {
			return err
		}
		url := fmt.Sprintf("%s/problems/%s/", leetcode.BaseURL, slug)
		fmt.Printf("Opening %s\n", url)
		return browser.Open(url)
	},
}

func init() {
	rootCmd.AddCommands(questionCmd)
	questionCmd.AddCommands(testCmd, submitCmd, starCmd, unstarCmd, openCmd)
}
