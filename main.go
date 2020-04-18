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
		Test()
		//fmt.Println("Invalid option")
	}
}

func Test() {
	fmt.Println(Background(Bright(Black)))
	for i := 0; i < 255; i++ {
		s := strconv.Itoa(i)
		fmt.Print(Foreground(TwoFiftySix(s)))
		fmt.Print(PadString(s, 4, true))
	}
}
