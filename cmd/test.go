package cmd

import (
	"fmt"
	"github.com/GregLahaye/linecode/leetcode"
)

var testcase string
var testCmd = &Command{
	Name:    "test",
	Aliases: []string{"t"},
	Run: func(cmd *Command, args []string) error {
		submission, err := leetcode.TestCode(args[0], testcase)
		if err != nil {
			return err
		}
		fmt.Println(submission)
		return nil
	},
	ArgN: 1,
}

func init() {
	rootCmd.AddCommands(testCmd)
}
