package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
)

func main() {
	u := User{}
	if err := u.Login(); err != nil {
		log.Fatal(err)
	}

	arg := os.Args[1]
	switch arg {
	case "list":
		if err := u.ListProblems(); err != nil {
			log.Fatal(err)
		}
	case "show":
		if question, err := u.GetQuestion(os.Args[2]); err != nil {
			log.Fatal(err)
		} else {
			PrettyPrint(question)
		}
	case "run":
		if code, err := ReadFile(os.Args[4]); err != nil {
			log.Fatal(err)
		} else {
			result, err := u.TestCode(1, os.Args[2], os.Args[3], string(code))
			if err != nil {
				log.Fatal(err)
			}
			submission, err := u.VerifyResult(result.InterpretID)
			if err != nil {
				log.Fatal(err)
			}
			for submission.State != "SUCCESS" {
				time.Sleep(time.Second * 1)
				submission, err = u.VerifyResult(result.InterpretID)
				if err != nil {
					log.Fatal(err)
				}
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
			submission, err := u.VerifyResult(strconv.Itoa(result.SubmissionID))
			if err != nil {
				log.Fatal(err)
			}
			for submission.State != "SUCCESS" {
				time.Sleep(time.Second * 1)
				submission, err = u.VerifyResult(strconv.Itoa(result.SubmissionID))
				if err != nil {
					log.Fatal(err)
				}
			}
			PrettyPrint(submission)
		}
	default:
		fmt.Println("Invalid option")
	}
}
