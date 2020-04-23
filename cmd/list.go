package cmd

import (
	"fmt"
	"github.com/GregLahaye/linecode/filter"
	"github.com/GregLahaye/linecode/leetcode"
	"sort"
)

var listHolder filter.Holder

var listCmd = &Command{
	Name:    "list",
	Aliases: []string{"l"},
	Flags:   filter.Flags("list", &listHolder),
	Run: func(cmd *Command, args []string) error {
		f := listHolder.Parse()

		problems, err := leetcode.GetProblems()
		if err != nil {
			return err
		}

		if f.Sort {
			s := func(i, j int) bool {
				a := problems[i]
				b := problems[j]
				x := float64(a.Stat.TotalAccepted) / float64(a.Stat.TotalSubmitted)
				y := float64(b.Stat.TotalAccepted) / float64(b.Stat.TotalSubmitted)
				return x < y
			}

			sort.Slice(problems, s)
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
	listCmd.Flags = filter.Flags("list", &listHolder)
}
