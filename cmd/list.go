package cmd

import "fmt"

var listCmd = &Command{
	Name: "list",
	Run: func(cmd *Command, args []string) error {
		// list questions
		f := fh.parse()
		fmt.Println(f)
		return nil
	},
}

func init() {
	rootCmd.AddCommands(listCmd)
	listCmd.Flags = filterFlags("list", &fh)
	listCmd.Flags.Bool("help", false, "display help")
}
