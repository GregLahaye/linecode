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
	case "test":
		return u.DisplayTest(args[1])
	case "submit":
		return u.DisplaySubmit(args[1])
	}

	return errors.New("unknown subcommand")
}

func main() {
	if err := root(os.Args[1:]); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

/*
func none() {
	switch subcommand {
	case "open":
		slug := os.Args[2]
		if id, err := strconv.Atoi(slug); err == nil {
			if s, err := u.GetSlug(id); err == nil {
				slug = s
			}
		}

		found := false
		problems, err := u.GetProblems()
		if err != nil {
			log.Fatal(err)
		}

		for _, p := range problems.Problems {
			if p.Stat.Slug == slug {
				found = true
			}
		}

		if !found {
			id, err := u.GetID(slug)
			if err != nil {
				log.Fatal(err)
			}

			slug, err = u.GetSlug(id)
			if err != nil {
				log.Fatal(err)
			}
		}

		if !Open(baseUrl + "/problems/" + slug + "/") {
			fmt.Println("Failed to open browser")
		}
	case "code":
		slug := os.Args[2]

		problems, err := u.GetProblems()
		if err != nil {
			log.Fatal(err)
		}

		var problem Problem
		found := false
		if id, err := strconv.Atoi(slug); err == nil {
			for _, p := range problems.Problems {
				if p.Stat.ID == id {
					problem = p
					found = true
				}
			}
		} else {
			for _, p := range problems.Problems {
				if p.Stat.Slug == slug {
					problem = p
					found = true
				}
			}
		}

		if !found {
			id, err := u.GetID(slug)
			if err != nil {
				log.Fatal(err)
			}

			for _, p := range problems.Problems {
				if p.Stat.ID == id {
					problem = p
					found = true
				}
			}
		}

		if found {
			q, err := u.GetQuestion(problem.Stat.Slug)
			if err != nil {
				log.Fatal(err)
			}

			if q.IsPaidOnly {
				fmt.Printf("%s is a locked question\n", q.Slug)
			} else {
				filename := IntToString(problem.Stat.ID) + "." + problem.Stat.Slug + "." + u.Language.Extension
				if _, err = os.Stat(filename); os.IsNotExist(err) {
					var code string
					for _, l := range q.CodeSnippets {
						if l.Slug == u.Language.Slug {
							code = l.Code
						}
					}

					c := u.Language.Comment.Start + "\n" + ParseHTML(q.Content) + "\n" + u.Language.Comment.End + "\n\n\n" + code

					err = ioutil.WriteFile(filename, []byte(c), os.ModePerm)
					if err != nil {
						log.Fatal(err)
					}
				}

				if err := u.OpenEditor(filename); err != nil {
					log.Fatal(err)
				}
			}
		} else {
			fmt.Println("Could not find question")
		}
	case "star":
		id, err := strconv.Atoi(os.Args[2])
		if err != nil {
			log.Fatal(err)
		}

		if err = u.Star(id); err != nil {
			log.Fatal(err)
		}
	case "unstar":
		id, err := strconv.Atoi(os.Args[2])
		if err != nil {
			log.Fatal(err)
		}

		if err = u.UnStar(id); err != nil {
			log.Fatal(err)
		}
	case "tags":
		if err := u.ListTags(); err != nil {
			log.Fatal(err)
		}
	case "stats":
		var r []rune
		if len(os.Args) > 2 {
			s := os.Args[2]
			r = []rune(s)
		}

		var tags []string
		if len(os.Args) > 3 {
			tags = os.Args[3:]
		}

		if err := u.DisplayStatistics(r, tags); err != nil {
			log.Fatal(err)
		}
	case "graph":
		if err := u.DisplayGraph(); err != nil {
			log.Fatal(err)
		}
	case "cache":
		dir, err := CacheDir("")
		if err != nil {
			log.Fatal(dir)
		}

		fmt.Println(dir)
	case "download":
		if err := u.DownloadAll(); err != nil {
			log.Fatal(err)
		}
	case "delete":
		files := map[string]string{"all": "", "chrome": "chrome", "user": userFilename, "problems": problemsFilename, "questions": questionsDirectory, "tags": tagsFilename}

		var filename string
		var found bool
		for k, v := range files {
			if k == os.Args[2] {
				filename = v
				found = false
			}
		}

		if found {
			if Confirm("Are you sure? (Y/N)") {
				if err := CacheDestroy(filename); err != nil {
					log.Fatal(err)
				}
			}
		} else {
			fmt.Printf("%s is not a valid option\n", os.Args[2])
		}
	default:
		fmt.Println("Invalid option")
	}
}
*/
