package cmd

import (
	"fmt"
	"github.com/GregLahaye/linecode/filter"
	"github.com/GregLahaye/linecode/leetcode"
)

var listCmd = &Command{
	Name:    "list",
	Aliases: []string{"l"},
	Run: func(cmd *Command, args []string) error {
		f := fh.Parse()

		problems, err := leetcode.GetProblems()
		if err != nil {
			return err
		}

		tags, err := leetcode.GetTags()
		if err != nil {
			return err
		}

		for _, p := range problems {
			if filter.Check(p, tags, f) {
				fmt.Println(p)
			}
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommands(listCmd)
	listCmd.Flags = filter.Flags("list", &fh)
	listCmd.Flags.Bool("help", false, "display help")
}
