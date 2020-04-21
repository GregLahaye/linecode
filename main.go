package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
)

const project = "linecode"
const baseUrl = "https://leetcode.com"

func root(args []string) error {
	if len(args) < 1 {
		return errors.New("You must pass a sub-command")
	}

	u, err := LoadUser()
	if err != nil {
		return err
	}

	subcommand := args[0]
	switch subcommand {
	case "list":
		f := parseFilterFlags()
		return u.ListProblems(f)
	case "show":
		fs := flag.NewFlagSet("", flag.ContinueOnError)
		save := fs.Bool("s", true, "save code snippet")
		open := fs.Bool("o", false, "open code snippet in editor")
		_ = fs.Parse(args[2:])

		return u.DisplayQuestion(args[1], *save, *open)
	case "open":
		problem, err := u.FindProblem(args[1])
		if err != nil {
			return err
		}

		url := baseUrl + "/problems/" + problem.Stat.Slug + "/"
		return Open(url)
	case "code":
		problem, err := u.FindProblem(args[1])
		if err != nil {
			return err
		}

		filename := IntToString(problem.Stat.ID) + "." + problem.Stat.Slug + "." + u.Language.Extension
		return u.OpenEditor(filename)
	case "test":
		return u.DisplayTest(args[1])
	case "submit":
		return u.DisplaySubmit(args[1])
	case "stats":
		f := parseFilterFlags()
		return u.DisplayStatistics(f)
	case "graph":
		return u.DisplayGraph()
	case "star":
		return u.Star(args[1])
	case "unstar":
		return u.Unstar(args[1])
	case "tags":
		return u.ListTags()
	case "download":
		return u.DownloadAll()
	case "destroy":
		return Destroy(args[1])
	}

	return errors.New("unknown subcommand")
}

func main() {
	if err := root(os.Args[1:]); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
