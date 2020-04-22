package cmd

import (
	"flag"
	"fmt"
	"github.com/GregLahaye/linecode/leetcode"
)

var save bool
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
		if save {
			leetcode.SaveSnippet(question)
		}
		fmt.Println(question)
		return nil
	},
	ArgN:  1,
	Flags: flag.NewFlagSet("question", flag.ContinueOnError),
}

var starCmd = &Command{
	Name: "star",
	Run: func(cmd *Command, args []string) error {
		arg := args[0]
		leetcode.Star(arg)
		return nil
	},
}

var unstarCmd = &Command{
	Name: "unstar",
	Run: func(cmd *Command, args []string) error {
		arg := args[0]
		leetcode.Unstar(arg)
		return nil
	},
}

func init() {
	rootCmd.AddCommands(questionCmd)
	questionCmd.AddCommands(starCmd, unstarCmd)
	questionCmd.Flags.BoolVar(&save, "s", true, "save question snippet")
}
