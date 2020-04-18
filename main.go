package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
)

func main() {
	defer fmt.Print(ForegroundReset, BackgroundReset)
	u := User{}
	if err := u.Login(); err != nil {
		log.Fatal(err)
	}

	fmt.Print(Foreground(Grey78))
	defer fmt.Println(ForegroundReset)

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

		if err = u.ShowQuestion(id); err != nil {
			log.Fatal(err)
		}
	case "run":
		id, err := strconv.Atoi(os.Args[2])
		if err != nil {
			log.Fatal(err)
		}
		language := os.Args[3]
		filename := os.Args[4]

		submission, err := u.RunCode(id, language, filename)
		if err != nil {
			log.Fatal(err)
		}
		PrettyPrint(submission)
	case "submit":
		id, err := strconv.Atoi(os.Args[2])
		if err != nil {
			log.Fatal(err)
		}
		language := os.Args[3]
		filename := os.Args[4]

		submission, err := u.SubmitCode(id, language, filename)
		if err != nil {
			log.Fatal(err)
		}
		PrettyPrint(submission)
	default:
		fmt.Println("Invalid option")
	}
}
