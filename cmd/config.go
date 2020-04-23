package cmd

import (
	"fmt"
	"github.com/GregLahaye/linecode/store"
)

var configCmd = &Command{
	Name: "config",
	Run: func(cmd *Command, args []string) error {
		fmt.Println(store.ConfigDir())
		return nil
	},
}

func init() {
	rootCmd.AddCommands(configCmd)
}
