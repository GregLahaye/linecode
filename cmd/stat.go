package cmd

import (
	"github.com/GregLahaye/linecode/filter"
	"github.com/GregLahaye/linecode/leetcode"
	"github.com/GregLahaye/linecode/linecode"
)

var statHolder filter.Holder

var statCmd = &Command{
	Name:    "stat",
	Aliases: []string{"stats"},
	Flags:   filter.Flags("stat", &statHolder),
	Run: func(cmd *Command, args []string) error {
		f := statHolder.Parse()

		all, err := leetcode.GetProblems()
		if err != nil {
			return err
		}

		tags, err := leetcode.GetTags()
		if err != nil {
			return err
		}

		var problems []linecode.Problem
		for _, problem := range all {
			if filter.Check(problem, tags, f) {
				problems = append(problems, problem)
			}
		}

		linecode.DisplayStat(problems)
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
	statCmd.Flags = filter.Flags("stat", &statHolder)
}
