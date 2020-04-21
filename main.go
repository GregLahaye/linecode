package main

import (
	"fmt"
	"github.com/GregLahaye/yogurt"
	"github.com/GregLahaye/yogurt/colors"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
)

func main() {
	defer fmt.Print(yogurt.ResetForeground, yogurt.ResetBackground)

	u, err := LoadUser()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Print(yogurt.Foreground(colors.Grey78))

	arg := os.Args[1]
	switch arg {
	case "list":
		var r []rune
		if len(os.Args) > 2 {
			s := os.Args[2]
			r = []rune(s)
		}

		var tags []string
		if len(os.Args) > 3 {
			tags = os.Args[3:]
		}

		if err := u.ListProblems(r, tags); err != nil {
			log.Fatal(err)
		}
	case "show":
		slug := os.Args[2]
		if id, err := strconv.Atoi(slug); err == nil {
			if s, err := u.GetSlug(id); err == nil {
				slug = s
			}
		}
		open := false
		if len(os.Args) > 3 {
			if os.Args[3][0] == 'o' {
				open = true
			}
		}

		if err = u.DisplayQuestion(slug, true, open); err != nil {
			log.Fatal(err)
		}
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
			if p.Stat.TitleSlug == slug {
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

		if !Open("https://leetcode.com/problems/" + slug + "/") {
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
				if p.Stat.TitleSlug == slug {
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
			q, err := u.GetQuestion(problem.Stat.TitleSlug)
			if err != nil {
				log.Fatal(err)
			}

			if q.IsPaidOnly {
				fmt.Printf("%s is a locked question\n", q.TitleSlug)
			} else {
				filename := IntToString(problem.Stat.ID) + "." + problem.Stat.TitleSlug + "." + u.Language.Extension
				if _, err = os.Stat(filename); os.IsNotExist(err) {
					var code string
					for _, l := range q.CodeSnippets {
						if l.LangSlug == u.Language.Slug {
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
	case "test":
		filename := os.Args[2]
		parts := strings.Split(filename, ".")
		id, err := strconv.Atoi(parts[0])
		if err != nil {
			log.Fatal(err)
		}
		slug := parts[1]

		fmt.Println("Please enter a testcase: (optional)")
		testcase, err := MultilineInput()
		if err != nil {
			log.Fatal(err)
		}

		submission, err := u.TestCode(id, slug, filename, testcase)
		if err != nil {
			log.Fatal(err)
		}
		DisplaySubmission(submission)
	case "submit":
		filename := os.Args[2]
		parts := strings.Split(filename, ".")
		id, err := strconv.Atoi(parts[0])
		if err != nil {
			log.Fatal(err)
		}
		slug := parts[1]

		submission, err := u.SubmitCode(id, slug, filename)
		if err != nil {
			log.Fatal(err)
		}
		submission.Judge = "large"
		DisplaySubmission(submission)
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
	case "destroy":
		if Confirm("Are you sure you want to delete all cached files? (Y/N) ") {
			if err := Destroy(""); err != nil {
				log.Fatal(err)
			}
		}
	default:
		fmt.Println("Invalid option")
	}
}
