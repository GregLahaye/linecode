package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
)

func main() {
	u := User{}
	if err := u.Login(); err != nil {
		log.Fatal(err)
	}

	fmt.Print(Foreground(Grey78))
	defer fmt.Println(ForegroundReset())

	arg := os.Args[1]
	switch arg {
	case "list":
		if err := u.ListProblems(); err != nil {
			log.Fatal(err)
		}
	case "show":
		if err := u.ShowQuestion(os.Args[2]); err != nil {
			log.Fatal(err)
		}
	case "run":
		if code, err := ReadFile(os.Args[4]); err != nil {
			log.Fatal(err)
		} else {
			result, err := u.TestCode(1, os.Args[2], os.Args[3], string(code))
			if err != nil {
				log.Fatal(err)
			}
			submission, err := u.Retry(result.InterpretID)
			if err != nil {
				log.Fatal(err)
			}
			PrettyPrint(submission)
		}
	case "submit":
		if code, err := ReadFile(os.Args[4]); err != nil {
			log.Fatal(err)
		} else {
			result, err := u.SubmitCode(1, os.Args[2], os.Args[3], string(code))
			if err != nil {
				log.Fatal(err)
			}
			submission, err := u.Retry(strconv.Itoa(result.SubmissionID))
			if err != nil {
				log.Fatal(err)
			}
			PrettyPrint(submission)
		}
	default:
		fmt.Println("Invalid option")
	}
}
