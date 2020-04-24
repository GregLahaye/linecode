package cmd

import (
	"flag"
	"fmt"
	"os"
)

type Command struct {
	Name     string
	Aliases  []string
	Run      func(cmd *Command, args []string) error
	NArg     int
	Flags    *flag.FlagSet
	commands []*Command
}

func (c *Command) AddCommands(args ...*Command) {
	for _, x := range args {
		c.commands = append(c.commands, x)
	}
}

func (c *Command) ParseFlags(args []string) error {
	if c.Flags == nil {
		return nil
	}

	return c.Flags.Parse(args)
}

func (c *Command) CountFlags() int {
	if c.Flags == nil {
		return 0
	}

	return c.Flags.NFlag()
}

func (c *Command) Execute() error {
	args := os.Args[1:]

	return rootCmd.execute(args)
}

func (c *Command) execute(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("no command given")
	}

	subcommand := args[0]
	args = args[1:]
	for _, x := range c.commands {
		if isMatch(x, subcommand) {
			// validate arguments
			if len(args) < x.NArg {
				return fmt.Errorf("not enough arguments")
			}

			// parse flags
			if err := x.ParseFlags(args); err != nil {
				return err
			}

			// there are sub-commands
			if x.commands != nil && len(args) > x.NArg+x.CountFlags() {
				// put positional arguments at end
				positional := args[:x.NArg]
				args = args[x.NArg:]
				args = append(args, positional...)

				return x.execute(args)
			}

			// run command
			return x.Run(x, args)
		}
	}

	return fmt.Errorf("command does not exist")
}

func isMatch(c *Command, s string) bool {
	if c.Name == s {
		return true
	}

	for _, x := range c.Aliases {
		if x == s {
			return true
		}
	}

	return false
}

var rootCmd = &Command{
	Name: "linecode",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
