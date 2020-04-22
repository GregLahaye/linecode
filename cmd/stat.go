package cmd

import "fmt"

var statCmd = &Command{
	Name: "stat",
	Aliases: []string{"stats"},
	Run: func(cmd *Command, args []string) error {
		fmt.Println("stat")
		return nil
	},
}

var graphCmd = &Command{
	Name: "graph",
	Aliases: []string{"g"},
	Run: func(cmd *Command, args []string) error {
		fmt.Println("graph")
		return nil
	},
}

func init() {
	rootCmd.AddCommands(statCmd)
	statCmd.AddCommands(graphCmd)
}

