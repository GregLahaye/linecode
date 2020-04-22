package cmd

import (
	"fmt"
	"github.com/GregLahaye/linecode/store"
)

var cacheCmd = &Command{
	Name: "cache",
	Run: func(cmd *Command, args []string) error {
		fmt.Println(store.CacheDir())
		return nil
	},
}

var removeCmd = &Command{
	Name:    "remove",
	Aliases: []string{"r", "delete", "d"},
	Run: func(cmd *Command, args []string) error {
		fmt.Println("remove", args)
		return nil
	},
	ArgN: 1,
}

func init() {
	rootCmd.AddCommands(cacheCmd)
	cacheCmd.AddCommands(removeCmd)
}
