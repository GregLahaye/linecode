package cmd

import (
	"flag"
	"fmt"
)

var questionCmd = &Command{
	Name: "question",
	Aliases: []string{"q"},
	Run: func(cmd *Command, args []string) error{
		fmt.Println("show this:", args)
		return nil
	},
	ArgN: 1,
	Flags: flag.NewFlagSet("show", flag.ContinueOnError),
}

var starCmd = &Command{
	Name: "star",
	Aliases: []string{"fav"},
	Run: func(cmd *Command, args []string) error{
		fmt.Println("star")
		return nil
	},
}

var unstarCmd = &Command{
	Name: "unstar",
	Aliases: []string{"unfav"},
	Run: func(cmd *Command, args []string) error{
		fmt.Println("unstar")
		return nil
	},
}

func init() {
	rootCmd.AddCommands(questionCmd)
	questionCmd.AddCommands(starCmd, unstarCmd)
}

