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
	defer fmt.Print(yogurt.ForegroundReset, yogurt.BackgroundReset)
	u := User{}
	if err := u.Login(); err != nil {
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
		lang := os.Args[3]

		if err = u.ShowQuestion(id, lang); err != nil {
			log.Fatal(err)
		}
	case "run":
		filename := os.Args[2]
		parts := strings.Split(filename, ".")
		id, err := strconv.Atoi(parts[0])
		if err != nil {
			log.Fatal(err)
		}
		ext := parts[len(parts)-1]
		lang := ExtensionToLanguage(ext)

		submission, err := u.RunCode(id, lang, filename)
		if err != nil {
			log.Fatal(err)
		}
		PrettyPrint(submission)
	case "submit":
		filename := os.Args[2]
		parts := strings.Split(filename, ".")
		id, err := strconv.Atoi(parts[0])
		if err != nil {
			log.Fatal(err)
		}
		ext := parts[len(parts)-1]
		lang := ExtensionToLanguage(ext)

		submission, err := u.SubmitCode(id, lang, filename)
		if err != nil {
			log.Fatal(err)
		}
		PrettyPrint(submission)
	default:
		fmt.Println("Invalid option")
	}
}
