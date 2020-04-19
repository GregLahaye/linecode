package main

import (
	"fmt"
	"github.com/GregLahaye/yogurt"
	"github.com/GregLahaye/yogurt/colors"
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
		if err := u.ListProblems(); err != nil {
			log.Fatal(err)
		}
	case "show":
		id, err := strconv.Atoi(os.Args[2])
		if err != nil {
			log.Fatal(err)
		}

		if err = u.DisplayQuestion(id); err != nil {
			log.Fatal(err)
		}
	case "test":
		filename := os.Args[2]
		parts := strings.Split(filename, ".")
		id, err := strconv.Atoi(parts[0])
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("Please enter a testcase:")
		testcase, err := MultilineInput()
		if err != nil {
			log.Fatal(err)
		}

		submission, err := u.TestCode(id, filename, testcase)
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

		submission, err := u.SubmitCode(id, filename)
		if err != nil {
			log.Fatal(err)
		}
		submission.Judge = "large"
		DisplaySubmission(submission)
	default:
		fmt.Println("Invalid option")
	}
}
