package cmd

import (
	"flag"
	"fmt"
	"github.com/GregLahaye/browser"
	"github.com/GregLahaye/input"
	"github.com/GregLahaye/linecode/leetcode"
	"github.com/GregLahaye/linecode/store"
	"os"
	"os/exec"
	"strings"
)

var questionCmd = &Command{
	Name:    "question",
	Aliases: []string{"q", "problem", "p"},
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
	NArg: 1,
}

var testcase string
var d bool
var testCmd = &Command{
	Name:    "test",
	Aliases: []string{"t", "run"},
	Run: func(cmd *Command, args []string) error {
		filename := leetcode.FindFile(args[0])
		testcase = strings.ReplaceAll(testcase, "\\n", "\n")

		if !d && strings.TrimSpace(testcase) == "" {
			testcase, _ = input.MultilineInput("Testcase (optional): ")
		}

		submission, err := leetcode.TestCode(filename, testcase)
		if err != nil {
			return err
		}
		fmt.Println(submission)
		return nil
	},
	NArg:  1,
	Flags: flag.NewFlagSet("test", flag.ContinueOnError),
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
	NArg: 1,
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

var editCmd = &Command{
	Name:    "edit",
	Aliases: []string{"e"},
	Run: func(cmd *Command, args []string) error {
		filename := leetcode.FindFile(args[0])
		if store.DoesNotExist(filename) {
			q, err := leetcode.GetQuestion(args[0])
			if err != nil {
				return err
			}
			if q.PaidOnly {
				return fmt.Errorf("%s is a locked question", q.Slug)
			}
			if err = leetcode.SaveSnippet(q); err != nil {
				return err
			}
		}
		e := exec.Command("vim", filename)
		e.Stdin = os.Stdin
		e.Stdout = os.Stdout
		err := e.Run()
		return err
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
	questionCmd.AddCommands(testCmd, submitCmd, starCmd, unstarCmd, editCmd, openCmd)
	testCmd.Flags.StringVar(&testcase, "t", "", "testcase")
	testCmd.Flags.BoolVar(&d, "d", false, "testcase")
}
