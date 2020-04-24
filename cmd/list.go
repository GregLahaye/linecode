package cmd

import (
	"fmt"
	"github.com/GregLahaye/linecode/filter"
	"github.com/GregLahaye/linecode/leetcode"
)

var listHolder filter.Holder

var listCmd = &Command{
	Name:    "list",
	Aliases: []string{"l"},
	Flags:   filter.PreferencesFlags("list", &listHolder),
	Run: func(cmd *Command, args []string) error {
		f := listHolder.Parse()

		all, err := leetcode.GetProblems()
		if err != nil {
			return err
		}

		tags, err := leetcode.GetTags()
		if err != nil {
			return err
		}

		problems := filter.Array(all, tags, f)

		for _, p := range problems {
			fmt.Println(p)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommands(listCmd)
}
