package cmd

import (
	"fmt"
	"github.com/GregLahaye/linecode/leetcode"
	"github.com/GregLahaye/linecode/linecode"
)

var statCmd = &Command{
	Name:    "stat",
	Aliases: []string{"stats"},
	Run: func(cmd *Command, args []string) error {
		fmt.Println("not yet implemented")
		return nil
	},
}

var graphCmd = &Command{
	Name:    "graph",
	Aliases: []string{"g"},
	Run: func(cmd *Command, args []string) error {
		problems, err := leetcode.GetProblems()
		if err != nil {
			return err
		}

		linecode.DisplayGraph(problems)

		return nil
	},
}

func init() {
	rootCmd.AddCommands(statCmd)
	statCmd.AddCommands(graphCmd)
}
