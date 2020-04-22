package cmd

import (
	"fmt"
	"github.com/GregLahaye/linecode/leetcode"
)

var submitCmd = &Command{
	Name:    "submit",
	Aliases: []string{"s"},
	Run: func(cmd *Command, args []string) error {
		submission, err := leetcode.SubmitCode(args[0])
		if err != nil {
			return err
		}
		submission.Judge = "large"
		fmt.Println(submission)
		return nil
	},
	ArgN: 1,
}

func init() {
	rootCmd.AddCommands(submitCmd)
}
